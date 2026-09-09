/*
Copyright 2026 The provider-azure-adx Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clusterpolicy

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/cluster/v1alpha1"
	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/fake"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

var showCols = []string{"PolicyName", "EntityName", "Policy", "ChildEntities", "EntityType"}

func ptr[T any](v T) *T { return &v }

// cluster simulates the cluster-level policy store of one Kusto cluster.
type cluster struct {
	policies map[string]string // name -> JSON ("" = null)
	// reformat rewrites the JSON stored by an alter, like Kusto echoing defaults.
	reformat func(name, js string) string
}

func newCluster() *cluster { return &cluster{policies: map[string]string{}} }

func (c *cluster) handle(db string, command cmd.Command) (*kusto.Result, error) {
	if db != "" {
		return nil, errors.New("cluster commands must run with an empty database, got " + db)
	}
	text := command.String()
	switch {
	case strings.HasPrefix(text, ".show cluster policy "):
		name := strings.TrimPrefix(text, ".show cluster policy ")
		var pol any
		if js := c.policies[name]; js != "" {
			pol = js
		}
		return kusto.NewResult(kusto.NewTable("Table_0", showCols, []any{name, "", pol, nil, "Cluster"})), nil
	// Both verbs land here: the capacity policy only accepts ".alter-merge"
	// (Def.Merge). This fake stores what it is given either way -- it does not
	// model the service's merge semantics, which the tests below don't rely on.
	case strings.HasPrefix(text, ".alter cluster policy "), strings.HasPrefix(text, ".alter-merge cluster policy "):
		rest := strings.TrimPrefix(strings.TrimPrefix(text, ".alter-merge cluster policy "), ".alter cluster policy ")
		name, body, _ := strings.Cut(rest, " ")
		js := body
		if i := strings.Index(body, "\n<| "); i >= 0 {
			// request classification: props @'...' then the query
			props := strings.TrimSuffix(strings.TrimPrefix(body[:i], "@'"), "'")
			query := body[i+len("\n<| "):]
			js = strings.TrimSuffix(props, "}") + `,"ClassificationFunction":` + jsonString(query) + "}"
		} else {
			js = strings.TrimSuffix(strings.TrimPrefix(js, "@'"), "'")
		}
		if c.reformat != nil {
			js = c.reformat(name, js)
		}
		c.policies[name] = js
		return kusto.NewResult(), nil
	case strings.HasPrefix(text, ".delete cluster policy "):
		delete(c.policies, strings.TrimPrefix(text, ".delete cluster policy "))
		return kusto.NewResult(), nil
	}
	return nil, errors.New("unexpected " + text)
}

func jsonString(s string) string {
	b := strings.Builder{}
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

func ext[T ClusterPolicy](kc kusto.Client, def Def[T]) *external[T] {
	return &external[T]{kc: kc, def: def}
}

// Deleting a cluster policy cannot be verified the usual way: ".show cluster
// policy X" keeps answering with Kusto's built-in default, so the resource
// always looks present. Before the fix the finalizer was never removed and the
// provider reissued ".delete cluster policy callout" once a minute forever
// (e2e run 34332137892).
func TestDeleteFinalizes(t *testing.T) {
	deleting := func(synced xpv2.Condition) *v1alpha1.CalloutPolicy {
		cr := &v1alpha1.CalloutPolicy{ObjectMeta: metav1.ObjectMeta{
			Name: "callout", Namespace: "ns",
			DeletionTimestamp: &metav1.Time{Time: time.Now()},
			Finalizers:        []string{"finalizer.managedresource.crossplane.io"},
		}}
		cr.Spec.ForProvider.Callouts = []v1alpha1.CalloutRule{{CalloutType: "sql", CalloutURIRegex: ".*", CanCall: true}}
		cr.SetConditions(xpv2.Deleting(), synced)
		return cr
	}
	for _, tc := range []struct {
		name       string
		cr         *v1alpha1.CalloutPolicy
		wantExists bool
	}{
		// Nothing has run yet: the policy must look present so that the
		// reconciler actually calls Delete and resets it to the default.
		{"before the delete ran", func() *v1alpha1.CalloutPolicy {
			cr := deleting(xpv2.ReconcileSuccess())
			cr.Status.ConditionedStatus = xpv2.ConditionedStatus{}
			return cr
		}(), true},
		// The delete ran and succeeded, so report it gone and let the
		// finalizer go, leaving the cluster on the default policy.
		{"delete succeeded", deleting(xpv2.ReconcileSuccess()), false},
		// The delete failed. Giving up here would drop the finalizer while
		// claiming a reset that never happened.
		{"delete failed", deleting(xpv2.ReconcileError(errors.New("boom"))), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cl := newCluster()
			cl.policies["callout"] = `[{"CalloutType":"sql","CalloutUriRegex":".*","CanCall":true}]`
			e := ext(fake.New("http://e").OnFn("", cl.handle), Callout())
			got, err := e.Observe(context.Background(), tc.cr)
			if err != nil {
				t.Fatal(err)
			}
			if got.ResourceExists != tc.wantExists {
				t.Errorf("ResourceExists = %v, want %v", got.ResourceExists, tc.wantExists)
			}
		})
	}
}

// Kusto has no delete command for every cluster policy: it answers with a
// syntax error at the policy name. Retrying that forever would keep the
// resource from finalizing, so Delete must treat it as "nothing to delete".
func TestDeleteToleratesMissingCommand(t *testing.T) {
	synErr := errors.New("SYN0002: Request is invalid and cannot be processed: Syntax error: SYN0002: A recognition error occurred. [line:position=1:23]")
	kc := fake.New("http://e").OnFn("", func(_ string, command cmd.Command) (*kusto.Result, error) {
		if strings.HasPrefix(command.String(), ".delete cluster policy ") {
			return nil, synErr
		}
		return kusto.NewResult(), nil
	})
	e := ext(kc, Sandbox())
	cr := &v1alpha1.SandboxPolicy{ObjectMeta: metav1.ObjectMeta{Name: "sandbox", Namespace: "ns"}}
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Errorf("a policy without a delete command must not fail the delete: %v", err)
	}
	// Any other permanent failure must still surface.
	kc = fake.New("http://e").OnFn("", func(_ string, command cmd.Command) (*kusto.Result, error) {
		if strings.HasPrefix(command.String(), ".delete cluster policy ") {
			return nil, errors.New("Forbidden: principal is not allowed")
		}
		return kusto.NewResult(), nil
	})
	if _, err := ext(kc, Sandbox()).Delete(context.Background(), cr); err == nil {
		t.Error("an unrelated delete failure must be reported")
	}
}

func TestCalloutLifecycle(t *testing.T) {
	cl := newCluster()
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := ext(kc, Callout())
	cr := &v1alpha1.CalloutPolicy{ObjectMeta: metav1.ObjectMeta{Name: "callout", Namespace: "ns"}}
	cr.Spec.ForProvider.Callouts = []v1alpha1.CalloutRule{{CalloutType: "sql", CalloutURIRegex: `.*\.database\.windows\.net`, CanCall: true}}

	got, err := e.Observe(context.Background(), cr)
	if err != nil || got.ResourceExists {
		t.Fatalf("initial observe: %+v %v", got, err)
	}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(cl.policies["callout"], `"CalloutType":"sql"`) {
		t.Errorf("stored: %s", cl.policies["callout"])
	}
	if cr.GetAnnotations()[AnnotationOwner] != "ns/callout" {
		t.Errorf("owner annotation: %v", cr.GetAnnotations())
	}
	if a, _ := base.Hashes(cr); a != "" {
		t.Error("structural policies must not record text hashes")
	}
	got, err = e.Observe(context.Background(), cr)
	if err != nil || !got.ResourceExists || !got.ResourceUpToDate {
		t.Fatalf("after create: %+v %v", got, err)
	}
	if cr.Status.AtProvider.Policy == "" || cr.GetCondition(xpv2.TypeReady).Status != "True" {
		t.Error("status must be filled and Available")
	}
	// Cluster drift: someone added a rule.
	cl.policies["callout"] = `[{"CalloutType":"sql","CalloutUriRegex":".*\\.database\\.windows\\.net","CanCall":true},{"CalloutType":"webapi","CalloutUriRegex":".*","CanCall":true}]`
	got, _ = e.Observe(context.Background(), cr)
	if got.ResourceUpToDate || !strings.Contains(got.Diff, "elements") {
		t.Fatalf("drift must be detected: %+v", got)
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	got, _ = e.Observe(context.Background(), cr)
	if !got.ResourceUpToDate {
		t.Fatal("update must repair drift")
	}
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if _, ok := cl.policies["callout"]; ok {
		t.Error("delete must remove the policy")
	}
	if cr.GetCondition(xpv2.TypeReady).Reason != xpv2.ReasonDeleting {
		t.Error("Delete must set Deleting")
	}
	if kc.Count(".delete cluster policy callout") != 1 {
		t.Errorf("commands: %v", kc.Commands())
	}
}

func TestCapacityAlwaysExistsNoDelete(t *testing.T) {
	cl := newCluster()
	cl.policies["capacity"] = `{"IngestionCapacity":{"ClusterMaximumConcurrentOperations":512,"CoreUtilizationCoefficient":0.75},"ExportCapacity":{"ClusterMaximumConcurrentOperations":100,"CoreUtilizationCoefficient":0.25}}`
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := ext(kc, Capacity())
	cr := &v1alpha1.CapacityPolicy{ObjectMeta: metav1.ObjectMeta{Name: "cap", Namespace: "ns"}}
	cr.Spec.ForProvider.Policy = apiextensionsv1.JSON{Raw: []byte(`{"IngestionCapacity":{"CoreUtilizationCoefficient":0.75}}`)}

	got, err := e.Observe(context.Background(), cr)
	if err != nil || !got.ResourceExists || !got.ResourceUpToDate {
		t.Fatalf("subset of defaults must be up to date: %+v %v", got, err)
	}
	cr.Spec.ForProvider.Policy = apiextensionsv1.JSON{Raw: []byte(`{"IngestionCapacity":{"CoreUtilizationCoefficient":0.5},"ExportCapacity":{"ClusterMaximumConcurrentOperations":100}}`)}
	got, _ = e.Observe(context.Background(), cr)
	if got.ResourceUpToDate || !strings.Contains(got.Diff, "CoreUtilizationCoefficient") {
		t.Fatalf("change must be detected: %+v", got)
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(cl.policies["capacity"], `"CoreUtilizationCoefficient":0.5`) {
		t.Errorf("stored: %s", cl.policies["capacity"])
	}
	before := len(kc.Commands())
	if _, err := e.Delete(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if len(kc.Commands()) != before || cl.policies["capacity"] == "" {
		t.Error("capacity policy cannot be deleted; Delete must be a no-op")
	}
	// Invalid spec JSON is a guardrail error.
	cr.Spec.ForProvider.Policy = apiextensionsv1.JSON{Raw: []byte(`{not json`)}
	if _, err := e.Observe(context.Background(), cr); !kerrors.IsBlocked(err) {
		t.Errorf("invalid JSON must be Blocked(InvalidSpec): %v", err)
	}
	cr.Spec.ForProvider.Policy = apiextensionsv1.JSON{Raw: []byte(`{}`)}
	if _, err := e.Observe(context.Background(), cr); !kerrors.IsBlocked(err) {
		t.Errorf("empty object must be Blocked(InvalidSpec): %v", err)
	}
}

func TestRequestClassificationKQLTolerance(t *testing.T) {
	cl := newCluster()
	cl.reformat = func(name, js string) string {
		if name != "request_classification" {
			return js
		}
		// Kusto echoes the function with collapsed whitespace.
		return strings.ReplaceAll(js, `\n`, " ")
	}
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := ext(kc, RequestClassification())
	cr := &v1alpha1.RequestClassificationPolicy{ObjectMeta: metav1.ObjectMeta{Name: "rc", Namespace: "ns"}}
	cr.Spec.ForProvider.Query = "iff(request_properties.current_application == \"Kusto.Explorer\",\n\"Ad-hoc queries\",\n\"default\")"

	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	alter := kc.Commands()[0]
	if !strings.HasPrefix(alter, ".alter cluster policy request_classification @'{\"IsEnabled\":true}'\n<| iff(") {
		t.Errorf("alter form: %q", alter)
	}
	if a, o := base.Hashes(cr); a == "" || o == "" {
		t.Fatal("KQL policies must record hashes")
	}
	for i := 0; i < 5; i++ {
		got, err := e.Observe(context.Background(), cr)
		if err != nil || !got.ResourceExists || !got.ResourceUpToDate {
			t.Fatalf("observe %d must be up to date despite reformatting: %+v %v", i, got, err)
		}
	}
	// Spec change -> detected.
	cr.Spec.ForProvider.Query = "\"default\""
	got, _ := e.Observe(context.Background(), cr)
	if got.ResourceUpToDate {
		t.Fatal("spec change must be detected")
	}
	if _, err := e.Update(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	got, _ = e.Observe(context.Background(), cr)
	if !got.ResourceUpToDate {
		t.Fatal("after update up to date")
	}
	// Disabled flag.
	cr.Spec.ForProvider.Enabled = ptr(false)
	got, _ = e.Observe(context.Background(), cr)
	if got.ResourceUpToDate || !strings.Contains(got.Diff, "IsEnabled") {
		t.Fatalf("IsEnabled change must be detected: %+v", got)
	}
	// Cluster drift in the function text.
	cr.Spec.ForProvider.Enabled = nil
	cl.policies["request_classification"] = `{"IsEnabled":true,"ClassificationFunction":"\"other\""}`
	got, _ = e.Observe(context.Background(), cr)
	if got.ResourceUpToDate {
		t.Fatal("cluster drift must be detected")
	}
	// Empty query is a spec error.
	cr.Spec.ForProvider.Query = "  "
	if _, err := e.Observe(context.Background(), cr); !kerrors.IsBlocked(err) {
		t.Errorf("empty query must be Blocked: %v", err)
	}
}

func TestMultiDatabaseAdminsHashOnly(t *testing.T) {
	cl := newCluster()
	cl.reformat = func(name, js string) string {
		// Shape unverified: pretend Kusto returns something else entirely.
		return `{"Whatever":["aadapp=1;t"]}`
	}
	kc := fake.New("http://e").OnFn("", cl.handle)
	e := ext(kc, MultiDatabaseAdmins())
	cr := &v1alpha1.MultiDatabaseAdminsPolicy{ObjectMeta: metav1.ObjectMeta{Name: "mda", Namespace: "ns"}}
	cr.Spec.ForProvider.Principals = []common.Principal{"aadapp=1;t"}
	if _, err := e.Create(context.Background(), cr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(kc.Commands()[0], `@'{"Principals":["aadapp=1;t"]}'`) {
		t.Errorf("alter: %s", kc.Commands()[0])
	}
	got, err := e.Observe(context.Background(), cr)
	if err != nil || !got.ResourceUpToDate {
		t.Fatalf("hash tolerance must accept the echoed shape: %+v %v", got, err)
	}
	cr.Spec.ForProvider.Principals = append(cr.Spec.ForProvider.Principals, "aaduser=x@y")
	got, _ = e.Observe(context.Background(), cr)
	if got.ResourceUpToDate {
		t.Fatal("spec change must be detected")
	}
	cl.policies["multidatabaseadmins"] = `{"Whatever":["changed"]}`
	cr.Spec.ForProvider.Principals = cr.Spec.ForProvider.Principals[:1]
	got, _ = e.Observe(context.Background(), cr)
	if got.ResourceUpToDate {
		t.Fatal("cluster drift must be detected")
	}
	cr.Spec.ForProvider.Principals = []common.Principal{" "}
	if _, err := e.Observe(context.Background(), cr); !kerrors.IsBlocked(err) {
		t.Errorf("empty principal must be Blocked: %v", err)
	}
}

func TestErrorsAndOtherKinds(t *testing.T) {
	boom := fake.New("http://e").On("", nil, errors.New("boom"))
	cr := &v1alpha1.SandboxPolicy{ObjectMeta: metav1.ObjectMeta{Name: "sb", Namespace: "ns"}}
	cr.Spec.ForProvider.Sandboxes = []v1alpha1.SandboxRule{{SandboxKind: "PythonExecution", MaxNumberOfSandboxes: ptr(int64(16))}}
	if _, err := ext(boom, Sandbox()).Observe(context.Background(), cr); err == nil || !strings.Contains(err.Error(), errObserve) {
		t.Errorf("observe error: %v", err)
	}
	if _, err := ext(boom, Sandbox()).Create(context.Background(), cr); err == nil {
		t.Error("create error expected")
	}
	if _, err := ext(boom, Sandbox()).Delete(context.Background(), cr); err == nil {
		t.Error("delete error expected")
	}
	nf := fake.New("http://e").On("", nil, fake.HTTPError(400, `{"error":{"code":"BadRequest_EntityNotFound","@message":"not found","@permanent":true}}`))
	if got, err := ext(nf, Sandbox()).Observe(context.Background(), cr); err != nil || got.ResourceExists {
		t.Errorf("not found: %+v %v", got, err)
	}
	if _, err := ext(nf, Sandbox()).Delete(context.Background(), cr); err != nil {
		t.Errorf("not found on delete is fine: %v", err)
	}

	d, err := Sandbox().Desired(cr)
	if err != nil || d.([]sandboxJSON)[0].SandboxKind != "PythonExecution" || d.([]sandboxJSON)[0].VirtualMachineSize != nil {
		t.Errorf("sandbox desired: %+v %v", d, err)
	}
	qw := &v1alpha1.QueryWeakConsistencyPolicy{}
	if _, err := QueryWeakConsistency().Desired(qw); err == nil {
		t.Error("empty query weak consistency must fail")
	}
	qw.Spec.ForProvider.PercentageOfNodes = ptr(int64(50))
	if d, err := QueryWeakConsistency().Desired(qw); err != nil || *d.(queryWeakConsistencyJSON).PercentageOfNodes != 50 {
		t.Errorf("qwc desired: %+v %v", d, err)
	}
	if _, ok := QueryWeakConsistency().Policy.DeleteCmd(); ok {
		t.Error("query weak consistency cannot be deleted")
	}
	mi := &v1alpha1.ClusterManagedIdentityPolicy{}
	mi.Spec.ForProvider.Identities = []v1alpha1.ClusterManagedIdentityEntry{{ObjectID: "system", AllowedUsages: []string{"NativeIngestion", "ExternalTable"}}}
	d, err = ManagedIdentity().Desired(mi)
	if err != nil || d.([]managedIdentityJSON)[0].AllowedUsages != "NativeIngestion, ExternalTable" {
		t.Errorf("managed identity desired: %+v %v", d, err)
	}
	c, err := ManagedIdentity().Policy.AlterCmd(d)
	if err != nil || c.String() != `.alter cluster policy managed_identity @'[{"ObjectId":"system","AllowedUsages":"NativeIngestion, ExternalTable"}]'` {
		t.Errorf("managed identity alter: %s %v", c, err)
	}
	if len(kindNames()) != 7 {
		t.Errorf("expected 7 cluster policy kinds, got %d", len(kindNames()))
	}
}

func kindNames() []string {
	return []string{Callout().Policy.Name, Capacity().Policy.Name, Sandbox().Policy.Name, QueryWeakConsistency().Policy.Name, ManagedIdentity().Policy.Name, MultiDatabaseAdmins().Policy.Name, RequestClassification().Policy.Name}
}
