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

// Package clusterpolicy reconciles all cluster-level *Policy managed
// resources with one generic controller, mirroring internal/controller/policy
// for database entities. Cluster policies have exactly one instance per
// cluster, so there is no entity: commands run at cluster level (database "").
package clusterpolicy

import (
	"context"
	"encoding/json"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/cluster/v1alpha1"
	adxclusterpolicy "github.com/functional-team/provider-azure-adx/internal/adx/clusterpolicy"
	adxpolicy "github.com/functional-team/provider-azure-adx/internal/adx/policy"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve  = "cannot observe cluster policy"
	errCreate   = "cannot create cluster policy"
	errUpdate   = "cannot update cluster policy"
	errDelete   = "cannot delete cluster policy"
	errReadBack = "cannot read cluster policy back after write"

	// AnnotationOwner records which managed resource last wrote the cluster
	// policy (spike S8). Two managed resources on the same cluster policy are
	// a user error; the annotation only helps to tell them apart when
	// debugging, the provider does not arbitrate between them.
	AnnotationOwner = "adx.functional.team/cluster-policy-owner"
)

// ClusterPolicy is implemented by every cluster policy kind (see
// apis/cluster/v1alpha1) so one generic controller can reconcile all of them.
type ClusterPolicy interface {
	resource.ModernManaged
	GetClusterPolicyObservation() v1alpha1.ClusterPolicyObservation
	SetClusterPolicyObservation(v1alpha1.ClusterPolicyObservation)
}

// Def ties a cluster policy kind to its Kusto definition and spec conversion.
type Def[T ClusterPolicy] struct {
	Kind   base.Kind[T]
	Policy adxclusterpolicy.Def
	// Desired converts the spec into the Kusto policy value.
	Desired func(T) (any, error)
}

// SetupKind adds the controller for one cluster policy kind.
func SetupKind[T ClusterPolicy](mgr ctrl.Manager, o controller.Options, d base.Deps, def Def[T]) error {
	return base.SetupGated(o, def.Kind, func() error {
		return base.Register(mgr, o, def.Kind, &connector[T]{base: base.NewConnector(mgr, d), def: def})
	})
}

type connector[T ClusterPolicy] struct {
	base *base.Connector
	def  Def[T]
}

func (c *connector[T]) Connect(ctx context.Context, cr T) (managed.TypedExternalClient[T], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external[T]{kc: kc, kube: c.base.Kube, def: c.def}, nil
}

type external[T ClusterPolicy] struct {
	kc   kusto.Client
	kube client.Client
	def  Def[T]
}

func (e *external[T]) kind() string { return e.def.Kind.GVK.Kind }

// observe returns the cluster's policy JSON. Cluster-level commands run with
// an empty database; the few cluster policies per cluster do not justify a
// cache section.
func (e *external[T]) observe(ctx context.Context) (json.RawMessage, bool, error) {
	res, err := e.kc.Mgmt(ctx, "", e.def.Policy.ShowCmd())
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	raw, ok := adxclusterpolicy.ParseShow(res)
	return raw, ok, nil
}

// resetDone reports whether the delete has already run and succeeded for this
// resource. The runtime records that as Ready=Deleting plus
// Synced=ReconcileSuccess before it requeues to verify the deletion. Checking
// Synced too is essential: Deleting is also marked when the delete itself
// failed, and treating that as done would drop the finalizer while claiming a
// reset that never happened.
func resetDone(cr ClusterPolicy) bool {
	return cr.GetCondition(xpv2.TypeReady).Reason == xpv2.ReasonDeleting &&
		cr.GetCondition(xpv2.TypeSynced).Reason == xpv2.ReasonReconcileSuccess
}

func (e *external[T]) Observe(ctx context.Context, cr T) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, e.kind(), "observe")
	// A cluster policy is a singleton that cannot be observed as absent: once
	// ours is deleted, ".show cluster policy X" answers with Kusto's built-in
	// default rather than null, so the resource looks like it still exists.
	// The finalizer was therefore never removed and the provider re-ran
	// ".delete cluster policy callout" once a minute, forever (e2e run
	// 34332137892). Report it gone once the delete has run, which leaves the
	// cluster on the default -- the point of deleting the resource.
	if meta.WasDeleted(cr) && resetDone(cr) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	desired, err := e.def.Desired(cr)
	if err != nil {
		return managed.ExternalObservation{}, kerrors.NewBlocked(kerrors.ReasonInvalidSpec, "%v", err)
	}
	raw, ok, err := e.observe(ctx)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	if !ok {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	res, err := e.def.Policy.CompareWith(desired, raw)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	cr.SetClusterPolicyObservation(v1alpha1.ClusterPolicyObservation{Policy: adxpolicy.Compact(raw)})
	cr.SetConditions(xpv2.Available())
	upToDate := res.Equal && base.TextUpToDate(cr, res.DesiredTexts, res.ObservedTexts)
	diff := ""
	if !upToDate {
		diff = res.Diff
		if diff == "" {
			diff = "KQL text differs"
		}
	}
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: upToDate, Diff: diff}, nil
}

// write alters the policy, reads it back and records the text hashes and the
// owner annotation on cr (in memory).
func (e *external[T]) write(ctx context.Context, cr T) error {
	desired, err := e.def.Desired(cr)
	if err != nil {
		return kerrors.NewBlocked(kerrors.ReasonInvalidSpec, "%v", err)
	}
	alter, err := e.def.Policy.AlterCmd(desired)
	if err != nil {
		return err
	}
	if _, err := e.kc.Mgmt(ctx, "", alter); err != nil {
		return err
	}
	res, err := e.kc.Mgmt(ctx, "", e.def.Policy.ShowCmd())
	if err != nil {
		return errors.Wrap(err, errReadBack)
	}
	raw, _ := adxclusterpolicy.ParseShow(res)
	cmp, err := e.def.Policy.CompareWith(desired, raw)
	if err != nil {
		return errors.Wrap(err, errReadBack)
	}
	if len(cmp.DesiredTexts) > 0 {
		base.SetHashes(cr, cmp.DesiredTexts, cmp.ObservedTexts)
	}
	meta.AddAnnotations(cr, map[string]string{AnnotationOwner: cr.GetNamespace() + "/" + cr.GetName()})
	return nil
}

func (e *external[T]) Create(ctx context.Context, cr T) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, e.kind(), "create")
	if err := e.write(ctx, cr); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	// Annotations set during Create are persisted by the reconciler.
	return managed.ExternalCreation{}, nil
}

func (e *external[T]) Update(ctx context.Context, cr T) (managed.ExternalUpdate, error) {
	ctx = kusto.WithOp(ctx, e.kind(), "update")
	if err := e.write(ctx, cr); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}
	if e.kube != nil {
		if err := base.PersistAnnotations(ctx, e.kube, cr); err != nil {
			return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
		}
	}
	return managed.ExternalUpdate{}, nil
}

func (e *external[T]) Delete(ctx context.Context, cr T) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, e.kind(), "delete")
	cr.SetConditions(xpv2.Deleting())
	del, ok := e.def.Policy.DeleteCmd()
	if !ok {
		// Kusto has no delete for this policy; the last applied values stay.
		return managed.ExternalDelete{}, nil
	}
	if _, err := e.kc.Mgmt(ctx, "", del); err != nil && !kerrors.IsNotFound(err) && !noDeleteCommand(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

// noDeleteCommand reports whether Kusto has no delete command for this policy
// at all. It answers with a syntax error pointing at the policy name, e.g.
// ".delete cluster policy multidatabaseadmins" fails with "SYN0002: A
// recognition error occurred. [line:position=1:23]" -- position 23 is exactly
// where the name starts, so the grammar does not accept it for delete.
//
// Retrying such a command can never succeed and would keep the resource from
// finalizing, so it counts as "nothing to delete": the policy keeps its values,
// like the NoDelete ones. Policies known to behave this way should get
// NoDelete instead, so no doomed command is sent in the first place; this is
// the net for the ones nobody has observed deleting yet
// (request_classification, managed_identity).
func noDeleteCommand(err error) bool {
	return strings.Contains(err.Error(), "SYN0002")
}

func (e *external[T]) Disconnect(_ context.Context) error { return nil }
