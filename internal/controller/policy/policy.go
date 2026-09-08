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

// Package policy reconciles all *Policy managed resources with one generic
// controller. Each kind contributes its Kusto definition (command shapes,
// comparison hints) and a function that converts the spec into the Kusto
// policy JSON; everything else is shared.
package policy

import (
	"context"
	"encoding/json"
	"sync/atomic"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/policy/v1alpha1"
	adxpolicy "github.com/functional-team/provider-azure-adx/internal/adx/policy"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve  = "cannot observe policy"
	errCreate   = "cannot create policy"
	errUpdate   = "cannot update policy"
	errDelete   = "cannot delete policy"
	errReadBack = "cannot read policy back after write"
)

// Policy is implemented by every policy kind (see apis/policy/v1alpha1) so
// one generic controller can reconcile all of them.
type Policy interface {
	resource.ModernManaged
	GetTarget() *v1alpha1.PolicyTarget
	GetPolicyObservation() v1alpha1.PolicyObservation
	SetPolicyObservation(v1alpha1.PolicyObservation)
}

// Def ties a policy kind to its Kusto definition and spec conversion.
type Def[T Policy] struct {
	Kind   base.Kind[T]
	Policy adxpolicy.Def
	// Desired converts the spec into the Kusto policy value (marshalled to
	// JSON for the default alter form and for comparison).
	Desired func(T) (any, error)
	// Column, if set, narrows the entity to one column (encoding policy).
	Column func(T) string
}

// state is shared by all reconciles of one kind.
type state struct {
	// batchUnsupported is set when ".show table * policy <name>" was rejected
	// by the cluster; the kind then falls back to single .show commands.
	batchUnsupported atomic.Bool
}

// SetupKind adds the controller for one policy kind.
func SetupKind[T Policy](mgr ctrl.Manager, o controller.Options, d base.Deps, def Def[T]) error {
	st := &state{}
	return base.SetupGated(o, def.Kind, func() error {
		return base.Register(mgr, o, def.Kind, &connector[T]{base: base.NewConnector(mgr, d), def: def, st: st})
	})
}

type connector[T Policy] struct {
	base *base.Connector
	def  Def[T]
	st   *state
}

func (c *connector[T]) Connect(ctx context.Context, cr T) (managed.TypedExternalClient[T], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external[T]{kc: kc, cache: c.base.Deps.Cache, kube: c.base.Kube, def: c.def, st: c.st}, nil
}

// NewExternal builds the external client for one policy kind directly
// (integration tests).
func NewExternal[T Policy](kc kusto.Client, cache *snapshot.Cache, kube client.Client, def Def[T]) managed.TypedExternalClient[T] {
	return &external[T]{kc: kc, cache: cache, kube: kube, def: def, st: &state{}}
}

type external[T Policy] struct {
	kc    kusto.Client
	cache *snapshot.Cache
	kube  client.Client
	def   Def[T]
	st    *state
}

func (e *external[T]) kind() string { return e.def.Kind.GVK.Kind }

func (e *external[T]) entity(cr T) (adxpolicy.Entity, error) {
	t := cr.GetTarget()
	ent := adxpolicy.Entity{Kind: t.Entity.Kind, Database: t.Database}
	if t.Entity.Kind != common.EntityKindDatabase {
		if t.Entity.Name == nil || *t.Entity.Name == "" {
			return ent, errors.New("entity.name is empty; set it or wait for entity.nameRef to resolve")
		}
		ent.Name = *t.Entity.Name
	}
	if e.def.Column != nil {
		ent.Column = e.def.Column(cr)
	}
	return ent, nil
}

func (e *external[T]) section() string { return "policy:" + e.def.Policy.Name }

func (e *external[T]) invalidate(db string) {
	e.cache.Invalidate(e.kc.Endpoint(), db, e.section())
}

// observe returns the entity's own policy JSON (not the inherited one).
func (e *external[T]) observe(ctx context.Context, ent adxpolicy.Entity) (json.RawMessage, bool, error) {
	if e.batchable(ent) {
		raw, ok, done, err := e.observeBatch(ctx, ent)
		if done {
			return raw, ok, err
		}
	}
	show, err := e.def.Policy.ShowCmd(ent)
	if err != nil {
		return nil, false, err
	}
	res, err := e.kc.Mgmt(ctx, ent.Database, show)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	raw, ok := adxpolicy.ParseShow(res)
	return raw, ok, nil
}

func (e *external[T]) batchable(ent adxpolicy.Entity) bool {
	return ent.Kind == common.EntityKindTable && ent.Column == "" && !e.def.Policy.NoBatch && !e.cache.Disabled() && !e.st.batchUnsupported.Load()
}

// observeBatch answers from the ".show table * policy <name>" section. done is
// false when the caller must fall back to a single .show.
func (e *external[T]) observeBatch(ctx context.Context, ent adxpolicy.Entity) (raw json.RawMessage, ok, done bool, err error) {
	all, err := snapshot.Load(ctx, e.cache, snapshot.Key{Endpoint: e.kc.Endpoint(), Database: ent.Database, Section: e.section()},
		func(ctx context.Context) (map[string]json.RawMessage, error) {
			res, err := e.kc.Mgmt(ctx, ent.Database, e.def.Policy.ShowAllTablesCmd())
			if err != nil {
				return nil, err
			}
			return adxpolicy.ParseShowAll(res), nil
		})
	switch {
	case err == nil:
		raw, ok = all[ent.Name]
		return raw, ok && !adxpolicy.IsNull(raw), true, nil
	case kerrors.IsNotFound(err):
		return nil, false, true, nil
	case kerrors.Classify(err) == kerrors.Permanent:
		// The wildcard form is not supported for this policy (spike S3);
		// remember that and use single .show commands from now on.
		e.st.batchUnsupported.Store(true)
		return nil, false, false, nil
	}
	return nil, false, true, err
}

func (e *external[T]) Observe(ctx context.Context, cr T) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, e.kind(), "observe")
	ent, err := e.entity(cr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	desired, err := e.def.Desired(cr)
	if err != nil {
		return managed.ExternalObservation{}, kerrors.NewBlocked(kerrors.ReasonInvalidSpec, "%v", err)
	}
	raw, ok, err := e.observe(ctx, ent)
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
	cr.SetPolicyObservation(v1alpha1.PolicyObservation{Entity: ent.Display(), Policy: adxpolicy.Compact(raw)})
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

// write alters the policy, reads it back and records the text hashes on cr.
func (e *external[T]) write(ctx context.Context, cr T) error {
	ent, err := e.entity(cr)
	if err != nil {
		return err
	}
	desired, err := e.def.Desired(cr)
	if err != nil {
		return kerrors.NewBlocked(kerrors.ReasonInvalidSpec, "%v", err)
	}
	alter, err := e.def.Policy.AlterCmd(ent, desired)
	if err != nil {
		return err
	}
	_, err = e.kc.Mgmt(ctx, ent.Database, alter)
	e.invalidate(ent.Database)
	if err != nil {
		return err
	}
	show, err := e.def.Policy.ShowCmd(ent)
	if err != nil {
		return err
	}
	res, err := e.kc.Mgmt(ctx, ent.Database, show)
	if err != nil {
		return errors.Wrap(err, errReadBack)
	}
	raw, _ := adxpolicy.ParseShow(res)
	cmp, err := e.def.Policy.CompareWith(desired, raw)
	if err != nil {
		return errors.Wrap(err, errReadBack)
	}
	if len(cmp.DesiredTexts) > 0 {
		base.SetHashes(cr, cmp.DesiredTexts, cmp.ObservedTexts)
	}
	return nil
}

func (e *external[T]) Create(ctx context.Context, cr T) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, e.kind(), "create")
	if err := e.write(ctx, cr); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	return managed.ExternalCreation{}, nil
}

func (e *external[T]) Update(ctx context.Context, cr T) (managed.ExternalUpdate, error) {
	ctx = kusto.WithOp(ctx, e.kind(), "update")
	if err := e.write(ctx, cr); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}
	if a, _ := base.Hashes(cr); a != "" && e.kube != nil {
		if err := base.PersistAnnotations(ctx, e.kube, cr); err != nil {
			return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
		}
	}
	return managed.ExternalUpdate{}, nil
}

func (e *external[T]) Delete(ctx context.Context, cr T) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, e.kind(), "delete")
	cr.SetConditions(xpv2.Deleting())
	ent, err := e.entity(cr)
	if err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	del, err := e.def.Policy.DeleteCmd(ent)
	if err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	_, err = e.kc.Mgmt(ctx, ent.Database, del)
	e.invalidate(ent.Database)
	if err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

func (e *external[T]) Disconnect(_ context.Context) error { return nil }
