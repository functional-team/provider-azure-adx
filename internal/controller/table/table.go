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

// Package table reconciles Table managed resources.
package table

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/adx/schema"
	"github.com/functional-team/provider-azure-adx/internal/adx/table"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve = "cannot observe table"
	errCreate  = "cannot create table"
	errUpdate  = "cannot update table"
	errDelete  = "cannot delete table"

	kind = "Table"
)

// Setup adds a controller that reconciles Table managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	k := base.Kind[*v1alpha1.Table]{GVK: v1alpha1.TableGroupVersionKind, Object: &v1alpha1.Table{}, List: &v1alpha1.TableList{}}
	return base.SetupGated(o, k, func() error {
		return base.Register(mgr, o, k, &connector{base: base.NewConnector(mgr, d)})
	})
}

type connector struct{ base *base.Connector }

func (c *connector) Connect(ctx context.Context, cr *v1alpha1.Table) (managed.TypedExternalClient[*v1alpha1.Table], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external{kc: kc, cache: c.base.Deps.Cache}, nil
}

// NewExternal builds the external client directly; used by integration tests
// that drive the lifecycle against the emulator. kube is unused for tables
// (they persist no annotations after Update) and kept for a uniform signature.
func NewExternal(kc kusto.Client, cache *snapshot.Cache, _ client.Client) managed.TypedExternalClient[*v1alpha1.Table] {
	return &external{kc: kc, cache: cache}
}

type external struct {
	kc    kusto.Client
	cache *snapshot.Cache
}

func (e *external) observe(ctx context.Context, db, name string) (table.Observed, bool, error) {
	if e.cache.Disabled() {
		res, err := e.kc.Mgmt(ctx, db, table.ShowCslSchema(name))
		if err != nil {
			if kerrors.IsNotFound(err) {
				return table.Observed{}, false, nil
			}
			return table.Observed{}, false, err
		}
		o, ok := table.ParseCslSchema(res)
		return o, ok, nil
	}
	d, err := schema.Cached(ctx, e.cache, e.kc, db)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return table.Observed{}, false, nil
		}
		return table.Observed{}, false, err
	}
	t, ok := d.Tables[name]
	if !ok {
		return table.Observed{}, false, nil
	}
	return table.FromSchema(t), true, nil
}

func (e *external) Observe(ctx context.Context, cr *v1alpha1.Table) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, kind, "observe")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	obs, ok, err := e.observe(ctx, db, name)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	if !ok {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	plan, err := table.Diff(table.FromParams(name, cr.Spec.ForProvider), obs)
	if err != nil {
		// Guardrail: surface the reason in Synced=False and do not call Update.
		cr.Status.AtProvider = table.Observation(obs, nil)
		return managed.ExternalObservation{}, err
	}
	cr.Status.AtProvider = table.Observation(obs, plan.DriftColumns)
	cr.Status.SetConditions(xpv2.Available())
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: plan.Empty(), Diff: plan.String()}, nil
}

func (e *external) Create(ctx context.Context, cr *v1alpha1.Table) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, kind, "create")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	for _, c := range table.BuildCreate(table.FromParams(name, cr.Spec.ForProvider)) {
		if _, err := e.kc.Mgmt(ctx, db, c); err != nil && !kerrors.IsAlreadyExists(err) {
			return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
		}
	}
	schema.Invalidate(e.cache, e.kc, db)
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, cr *v1alpha1.Table) (managed.ExternalUpdate, error) {
	ctx = kusto.WithOp(ctx, kind, "update")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	obs, ok, err := e.observe(ctx, db, name)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}
	if !ok {
		return managed.ExternalUpdate{}, nil // gone; next Observe recreates
	}
	plan, err := table.Diff(table.FromParams(name, cr.Spec.ForProvider), obs)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}
	for _, s := range plan.Steps {
		if _, err := e.kc.Mgmt(ctx, db, s.Command); err != nil {
			schema.Invalidate(e.cache, e.kc, db)
			return managed.ExternalUpdate{}, errors.Wrapf(err, "%s (%s)", errUpdate, s.Kind)
		}
	}
	schema.Invalidate(e.cache, e.kc, db)
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, cr *v1alpha1.Table) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, kind, "delete")
	cr.Status.SetConditions(xpv2.Deleting())
	db := cr.Spec.ForProvider.Database
	if _, err := e.kc.Mgmt(ctx, db, table.BuildDelete(meta.GetExternalName(cr))); err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	schema.Invalidate(e.cache, e.kc, db)
	return managed.ExternalDelete{}, nil
}

func (e *external) Disconnect(_ context.Context) error { return nil }
