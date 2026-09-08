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

// Package ingestionmapping reconciles IngestionMapping managed resources.
package ingestionmapping

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	"github.com/functional-team/provider-azure-adx/internal/adx/ingestionmapping"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve = "cannot observe ingestion mapping"
	errCreate  = "cannot create ingestion mapping"
	errUpdate  = "cannot update ingestion mapping"
	errDelete  = "cannot delete ingestion mapping"
	errNoTable = "spec.forProvider.table is empty; set it or wait for tableRef to resolve"

	kind = "IngestionMapping"
)

// Setup adds a controller that reconciles IngestionMapping managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	k := base.Kind[*v1alpha1.IngestionMapping]{GVK: v1alpha1.IngestionMappingGroupVersionKind, Object: &v1alpha1.IngestionMapping{}, List: &v1alpha1.IngestionMappingList{}}
	return base.SetupGated(o, k, func() error {
		return base.Register(mgr, o, k, &connector{base: base.NewConnector(mgr, d)})
	})
}

type connector struct{ base *base.Connector }

func (c *connector) Connect(ctx context.Context, cr *v1alpha1.IngestionMapping) (managed.TypedExternalClient[*v1alpha1.IngestionMapping], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external{kc: kc, cache: c.base.Deps.Cache}, nil
}

type external struct {
	kc    kusto.Client
	cache *snapshot.Cache
}

func (e *external) observe(ctx context.Context, d ingestionmapping.Desired, db string) (ingestionmapping.Observed, bool, error) {
	if e.cache.Disabled() {
		res, err := e.kc.Mgmt(ctx, db, ingestionmapping.ShowOne(d.Table, d.Kind, d.Name))
		if err != nil {
			if kerrors.IsNotFound(err) {
				return ingestionmapping.Observed{}, false, nil
			}
			return ingestionmapping.Observed{}, false, err
		}
		return ingestionmapping.ParseOne(res, d.Kind, d.Name)
	}
	all, err := snapshot.Load(ctx, e.cache, snapshot.Key{Endpoint: e.kc.Endpoint(), Database: db, Section: ingestionmapping.Section(d.Table)},
		func(ctx context.Context) (map[string]ingestionmapping.Observed, error) {
			res, err := e.kc.Mgmt(ctx, db, ingestionmapping.ShowAll(d.Table))
			if err != nil {
				return nil, err
			}
			return ingestionmapping.ParseRows(res)
		})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return ingestionmapping.Observed{}, false, nil
		}
		return ingestionmapping.Observed{}, false, err
	}
	o, ok := all[ingestionmapping.Key(d.Kind, d.Name)]
	return o, ok, nil
}

func (e *external) invalidate(db, table string) {
	e.cache.Invalidate(e.kc.Endpoint(), db, ingestionmapping.Section(table))
}

func desired(cr *v1alpha1.IngestionMapping) (ingestionmapping.Desired, error) {
	d := ingestionmapping.FromParams(meta.GetExternalName(cr), cr.Spec.ForProvider)
	if d.Table == "" {
		return d, errors.New(errNoTable)
	}
	return d, nil
}

func (e *external) Observe(ctx context.Context, cr *v1alpha1.IngestionMapping) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, kind, "observe")
	d, err := desired(cr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	obs, ok, err := e.observe(ctx, d, cr.Spec.ForProvider.Database)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	if !ok {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	cr.Status.AtProvider = ingestionmapping.Observation(obs)
	cr.Status.SetConditions(xpv2.Available())
	equal, diff := ingestionmapping.Equal(d, obs)
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: equal, Diff: diff}, nil
}

func (e *external) write(ctx context.Context, cr *v1alpha1.IngestionMapping) error {
	d, err := desired(cr)
	if err != nil {
		return err
	}
	c, err := ingestionmapping.BuildCreateOrAlter(d)
	if err != nil {
		return err
	}
	_, err = e.kc.Mgmt(ctx, cr.Spec.ForProvider.Database, c)
	e.invalidate(cr.Spec.ForProvider.Database, d.Table)
	return err
}

func (e *external) Create(ctx context.Context, cr *v1alpha1.IngestionMapping) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, kind, "create")
	if err := e.write(ctx, cr); err != nil && !kerrors.IsAlreadyExists(err) {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, cr *v1alpha1.IngestionMapping) (managed.ExternalUpdate, error) {
	ctx = kusto.WithOp(ctx, kind, "update")
	if err := e.write(ctx, cr); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, cr *v1alpha1.IngestionMapping) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, kind, "delete")
	cr.Status.SetConditions(xpv2.Deleting())
	d, err := desired(cr)
	if err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	_, err = e.kc.Mgmt(ctx, cr.Spec.ForProvider.Database, ingestionmapping.BuildDelete(d.Table, d.Kind, d.Name))
	e.invalidate(cr.Spec.ForProvider.Database, d.Table)
	if err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

func (e *external) Disconnect(_ context.Context) error { return nil }
