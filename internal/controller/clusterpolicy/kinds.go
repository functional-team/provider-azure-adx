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
	"encoding/json"
	"fmt"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"

	"github.com/functional-team/provider-azure-adx/apis/cluster/v1alpha1"
	adxclusterpolicy "github.com/functional-team/provider-azure-adx/internal/adx/clusterpolicy"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

// SetupAll adds the controllers of every cluster policy kind.
func SetupAll(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	for _, setup := range []func() error{
		func() error { return SetupKind(mgr, o, d, Callout()) },
		func() error { return SetupKind(mgr, o, d, Capacity()) },
		func() error { return SetupKind(mgr, o, d, Sandbox()) },
		func() error { return SetupKind(mgr, o, d, QueryWeakConsistency()) },
		func() error { return SetupKind(mgr, o, d, ManagedIdentity()) },
		func() error { return SetupKind(mgr, o, d, MultiDatabaseAdmins()) },
		func() error { return SetupKind(mgr, o, d, RequestClassification()) },
	} {
		if err := setup(); err != nil {
			return err
		}
	}
	return nil
}

func boolOr(b *bool, def bool) bool {
	if b == nil {
		return def
	}
	return *b
}

// rawJSON parses a spec JSON blob into a generic value for comparison and
// marshalling. Empty or invalid JSON is a spec error.
func rawJSON(raw []byte, field string) (any, error) {
	if len(strings.TrimSpace(string(raw))) == 0 || strings.TrimSpace(string(raw)) == "null" {
		return nil, fmt.Errorf("%s is empty", field)
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("%s is not valid JSON: %w", field, err)
	}
	if m, ok := v.(map[string]any); ok && len(m) == 0 {
		return nil, fmt.Errorf("%s is an empty object", field)
	}
	return v, nil
}

type calloutJSON struct {
	CalloutType     string `json:"CalloutType"`
	CalloutURIRegex string `json:"CalloutUriRegex"`
	CanCall         bool   `json:"CanCall"`
}

// Callout defines the CalloutPolicy kind (JSON array).
func Callout() Def[*v1alpha1.CalloutPolicy] {
	return Def[*v1alpha1.CalloutPolicy]{
		Kind:   base.Kind[*v1alpha1.CalloutPolicy]{GVK: v1alpha1.CalloutPolicyGroupVersionKind, Object: &v1alpha1.CalloutPolicy{}, List: &v1alpha1.CalloutPolicyList{}},
		Policy: adxclusterpolicy.Def{Name: "callout"},
		Desired: func(cr *v1alpha1.CalloutPolicy) (any, error) {
			out := make([]calloutJSON, 0, len(cr.Spec.ForProvider.Callouts))
			for _, c := range cr.Spec.ForProvider.Callouts {
				out = append(out, calloutJSON{CalloutType: c.CalloutType, CalloutURIRegex: c.CalloutURIRegex, CanCall: c.CanCall})
			}
			return out, nil
		},
	}
}

// Capacity defines the CapacityPolicy kind. The spec carries the raw policy
// JSON; only its keys are compared. Kusto cannot delete the capacity policy
// and only accepts it via ".alter-merge", hence Merge.
func Capacity() Def[*v1alpha1.CapacityPolicy] {
	return Def[*v1alpha1.CapacityPolicy]{
		Kind:   base.Kind[*v1alpha1.CapacityPolicy]{GVK: v1alpha1.CapacityPolicyGroupVersionKind, Object: &v1alpha1.CapacityPolicy{}, List: &v1alpha1.CapacityPolicyList{}},
		Policy: adxclusterpolicy.Def{Name: "capacity", NoDelete: true, Merge: true},
		Desired: func(cr *v1alpha1.CapacityPolicy) (any, error) {
			return rawJSON(cr.Spec.ForProvider.Policy.Raw, "policy")
		},
	}
}

type sandboxJSON struct {
	SandboxKind           string  `json:"SandboxKind"`
	VirtualMachineSize    *string `json:"VirtualMachineSize,omitempty"`
	InitializeOnStartup   *bool   `json:"InitializeOnStartup,omitempty"`
	MaxNumberOfSandboxes  *int64  `json:"MaxNumberOfSandboxes,omitempty"`
	MaxCPUPerSandbox      *int64  `json:"MaxCpuPerSandbox,omitempty"`
	MaxMemoryMbPerSandbox *int64  `json:"MaxMemoryMbPerSandbox,omitempty"`
}

// Sandbox defines the SandboxPolicy kind (JSON array).
func Sandbox() Def[*v1alpha1.SandboxPolicy] {
	return Def[*v1alpha1.SandboxPolicy]{
		Kind:   base.Kind[*v1alpha1.SandboxPolicy]{GVK: v1alpha1.SandboxPolicyGroupVersionKind, Object: &v1alpha1.SandboxPolicy{}, List: &v1alpha1.SandboxPolicyList{}},
		Policy: adxclusterpolicy.Def{Name: "sandbox"},
		Desired: func(cr *v1alpha1.SandboxPolicy) (any, error) {
			out := make([]sandboxJSON, 0, len(cr.Spec.ForProvider.Sandboxes))
			for _, s := range cr.Spec.ForProvider.Sandboxes {
				out = append(out, sandboxJSON{SandboxKind: s.SandboxKind, VirtualMachineSize: s.VirtualMachineSize, InitializeOnStartup: s.InitializeOnStartup, MaxNumberOfSandboxes: s.MaxNumberOfSandboxes, MaxCPUPerSandbox: s.MaxCPUPerSandbox, MaxMemoryMbPerSandbox: s.MaxMemoryMbPerSandbox})
			}
			return out, nil
		},
	}
}

type queryWeakConsistencyJSON struct {
	PercentageOfNodes                  *int64 `json:"PercentageOfNodes,omitempty"`
	MinimumNumberOfNodes               *int64 `json:"MinimumNumberOfNodes,omitempty"`
	MaximumNumberOfNodes               *int64 `json:"MaximumNumberOfNodes,omitempty"`
	SuperSlackerNumberOfNodesThreshold *int64 `json:"SuperSlackerNumberOfNodesThreshold,omitempty"`
	EnableMetadataPrefetch             *bool  `json:"EnableMetadataPrefetch,omitempty"`
	MaximumLagAllowedInMinutes         *int64 `json:"MaximumLagAllowedInMinutes,omitempty"`
	RefreshPeriodInSeconds             *int64 `json:"RefreshPeriodInSeconds,omitempty"`
}

// QueryWeakConsistency defines the QueryWeakConsistencyPolicy kind. Kusto
// cannot delete this policy.
func QueryWeakConsistency() Def[*v1alpha1.QueryWeakConsistencyPolicy] {
	return Def[*v1alpha1.QueryWeakConsistencyPolicy]{
		Kind:   base.Kind[*v1alpha1.QueryWeakConsistencyPolicy]{GVK: v1alpha1.QueryWeakConsistencyPolicyGroupVersionKind, Object: &v1alpha1.QueryWeakConsistencyPolicy{}, List: &v1alpha1.QueryWeakConsistencyPolicyList{}},
		Policy: adxclusterpolicy.Def{Name: "query_weak_consistency", NoDelete: true},
		Desired: func(cr *v1alpha1.QueryWeakConsistencyPolicy) (any, error) {
			p := cr.Spec.ForProvider
			d := queryWeakConsistencyJSON{PercentageOfNodes: p.PercentageOfNodes, MinimumNumberOfNodes: p.MinimumNumberOfNodes, MaximumNumberOfNodes: p.MaximumNumberOfNodes, SuperSlackerNumberOfNodesThreshold: p.SuperSlackerNumberOfNodesThreshold, EnableMetadataPrefetch: p.EnableMetadataPrefetch, MaximumLagAllowedInMinutes: p.MaximumLagAllowedInMinutes, RefreshPeriodInSeconds: p.RefreshPeriodInSeconds}
			if d == (queryWeakConsistencyJSON{}) {
				return nil, fmt.Errorf("at least one field must be set")
			}
			return d, nil
		},
	}
}

type managedIdentityJSON struct {
	ObjectID      string `json:"ObjectId"`
	AllowedUsages string `json:"AllowedUsages"`
}

// ManagedIdentity defines the ClusterManagedIdentityPolicy kind (same JSON as
// the database level policy, ".alter cluster policy managed_identity").
func ManagedIdentity() Def[*v1alpha1.ClusterManagedIdentityPolicy] {
	return Def[*v1alpha1.ClusterManagedIdentityPolicy]{
		Kind:   base.Kind[*v1alpha1.ClusterManagedIdentityPolicy]{GVK: v1alpha1.ClusterManagedIdentityPolicyGroupVersionKind, Object: &v1alpha1.ClusterManagedIdentityPolicy{}, List: &v1alpha1.ClusterManagedIdentityPolicyList{}},
		Policy: adxclusterpolicy.Def{Name: "managed_identity", SetFields: []string{"AllowedUsages"}},
		Desired: func(cr *v1alpha1.ClusterManagedIdentityPolicy) (any, error) {
			out := make([]managedIdentityJSON, 0, len(cr.Spec.ForProvider.Identities))
			for _, id := range cr.Spec.ForProvider.Identities {
				out = append(out, managedIdentityJSON{ObjectID: id.ObjectID, AllowedUsages: strings.Join(id.AllowedUsages, ", ")})
			}
			return out, nil
		},
	}
}

// multiDatabaseAdminsJSON is the assumed shape of the multidatabaseadmins
// policy. UNVERIFIED against a cluster (spike): the kind compares through the
// hash annotations only, so a different shape shows up as a Kusto error on
// the first write rather than as an endless update loop.
type multiDatabaseAdminsJSON struct {
	Principals []string `json:"Principals"`
}

// MultiDatabaseAdmins defines the MultiDatabaseAdminsPolicy kind.
func MultiDatabaseAdmins() Def[*v1alpha1.MultiDatabaseAdminsPolicy] {
	return Def[*v1alpha1.MultiDatabaseAdminsPolicy]{
		Kind: base.Kind[*v1alpha1.MultiDatabaseAdminsPolicy]{GVK: v1alpha1.MultiDatabaseAdminsPolicyGroupVersionKind, Object: &v1alpha1.MultiDatabaseAdminsPolicy{}, List: &v1alpha1.MultiDatabaseAdminsPolicyList{}},
		// Kusto has no ".delete cluster policy multidatabaseadmins": it answers
		// with a syntax error at the policy name (observed 2026-09-09), so
		// deleting the resource leaves the policy as it is.
		Policy: adxclusterpolicy.Def{Name: "multidatabaseadmins", HashOnly: true, NoDelete: true},
		Desired: func(cr *v1alpha1.MultiDatabaseAdminsPolicy) (any, error) {
			out := multiDatabaseAdminsJSON{Principals: make([]string, 0, len(cr.Spec.ForProvider.Principals))}
			for _, p := range cr.Spec.ForProvider.Principals {
				if strings.TrimSpace(string(p)) == "" {
					return nil, fmt.Errorf("principals must not contain empty entries")
				}
				out.Principals = append(out.Principals, string(p))
			}
			return out, nil
		},
	}
}

type requestClassificationJSON struct {
	IsEnabled              bool   `json:"IsEnabled"`
	ClassificationFunction string `json:"ClassificationFunction"`
}

// RequestClassification defines the RequestClassificationPolicy kind:
// ".alter cluster policy request_classification @'{"IsEnabled":true}' <| query".
// The .show JSON carries IsEnabled and ClassificationFunction (the query).
func RequestClassification() Def[*v1alpha1.RequestClassificationPolicy] {
	return Def[*v1alpha1.RequestClassificationPolicy]{
		Kind: base.Kind[*v1alpha1.RequestClassificationPolicy]{GVK: v1alpha1.RequestClassificationPolicyGroupVersionKind, Object: &v1alpha1.RequestClassificationPolicy{}, List: &v1alpha1.RequestClassificationPolicyList{}},
		Policy: adxclusterpolicy.Def{Name: "request_classification", KQLFields: []string{"ClassificationFunction"}, Alter: func(desired any) (cmd.Command, error) {
			d, ok := desired.(requestClassificationJSON)
			if !ok {
				return cmd.Command{}, fmt.Errorf("unexpected desired type %T", desired)
			}
			props, err := cmd.JSON(struct {
				IsEnabled bool `json:"IsEnabled"`
			}{IsEnabled: d.IsEnabled})
			if err != nil {
				return cmd.Command{}, err
			}
			return cmd.New(".alter cluster policy request_classification ", props, cmd.Pipe(d.ClassificationFunction)), nil
		}},
		Desired: func(cr *v1alpha1.RequestClassificationPolicy) (any, error) {
			p := cr.Spec.ForProvider
			if strings.TrimSpace(p.Query) == "" {
				return nil, fmt.Errorf("query is empty")
			}
			return requestClassificationJSON{IsEnabled: boolOr(p.Enabled, true), ClassificationFunction: p.Query}, nil
		},
	}
}
