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

// Package function reconciles Function managed resources.
package function

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
	"github.com/functional-team/provider-azure-adx/internal/adx/function"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve = "cannot observe function"
	errCreate  = "cannot create function"
	errUpdate  = "cannot update function"
	errDelete  = "cannot delete function"

	kind = "Function"
)

// Setup adds a controller that reconciles Function managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	k := base.Kind[*v1alpha1.Function]{GVK: v1alpha1.FunctionGroupVersionKind, Object: &v1alpha1.Function{}, List: &v1alpha1.FunctionList{}}
	return base.SetupGated(o, k, func() error {
		return base.Register(mgr, o, k, &connector{base: base.NewConnector(mgr, d)})
	})
}

type connector struct{ base *base.Connector }

func (c *connector) Connect(ctx context.Context, cr *v1alpha1.Function) (managed.TypedExternalClient[*v1alpha1.Function], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external{kc: kc, cache: c.base.Deps.Cache, kube: c.base.Kube}, nil
}

// NewExternal builds the external client directly; used by integration tests
// that drive the lifecycle against the emulator. With a nil kube client the
// hash annotations are kept in memory only.
func NewExternal(kc kusto.Client, cache *snapshot.Cache, kube client.Client) managed.TypedExternalClient[*v1alpha1.Function] {
	return &external{kc: kc, cache: cache, kube: kube}
}

type external struct {
	kc    kusto.Client
	cache *snapshot.Cache
	kube  client.Client
}

func (e *external) observe(ctx context.Context, db, name string) (function.Observed, bool, error) {
	if e.cache.Disabled() {
		res, err := e.kc.Mgmt(ctx, db, function.ShowFunction(name))
		if err != nil {
			if kerrors.IsNotFound(err) {
				return function.Observed{}, false, nil
			}
			return function.Observed{}, false, err
		}
		o, ok := function.ParseOne(res, name)
		return o, ok, nil
	}
	all, err := snapshot.Load(ctx, e.cache, snapshot.Key{Endpoint: e.kc.Endpoint(), Database: db, Section: function.Section}, func(ctx context.Context) (map[string]function.Observed, error) {
		res, err := e.kc.Mgmt(ctx, db, function.ShowFunctions())
		if err != nil {
			return nil, err
		}
		return function.ParseRows(res), nil
	})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return function.Observed{}, false, nil
		}
		return function.Observed{}, false, err
	}
	o, ok := all[name]
	return o, ok, nil
}

func (e *external) invalidate(db string) {
	e.cache.Invalidate(e.kc.Endpoint(), db, function.Section)
}

func (e *external) Observe(ctx context.Context, cr *v1alpha1.Function) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, kind, "observe")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	obs, ok, err := e.observe(ctx, db, name)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	if !ok {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	d := function.FromParams(name, cr.Spec.ForProvider)
	cr.Status.AtProvider = function.Observation(obs)
	cr.Status.SetConditions(xpv2.Available())
	upToDate := base.TextUpToDate(cr, function.DesiredTexts(d), function.ObservedTexts(obs)) && function.MetadataUpToDate(d, obs)
	diff := ""
	if !upToDate {
		diff = "body, parameters or metadata differ"
	}
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: upToDate, Diff: diff}, nil
}

// write runs create-or-alter, reads the function back and records the hashes
// (in memory) on cr.
func (e *external) write(ctx context.Context, cr *v1alpha1.Function) error {
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	d := function.FromParams(name, cr.Spec.ForProvider)
	res, err := e.kc.Mgmt(ctx, db, function.BuildCreateOrAlter(d))
	e.invalidate(db)
	if err != nil {
		return err
	}
	obs, ok := function.ParseOne(res, name)
	if !ok {
		// The command output did not echo the function; read it back.
		res, err = e.kc.Mgmt(ctx, db, function.ShowFunction(name))
		if err != nil {
			return err
		}
		obs, ok = function.ParseOne(res, name)
		if !ok {
			return errors.New("function not found after write")
		}
	}
	base.SetHashes(cr, function.DesiredTexts(d), function.ObservedTexts(obs))
	return nil
}

func (e *external) Create(ctx context.Context, cr *v1alpha1.Function) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, kind, "create")
	if err := e.write(ctx, cr); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	// Annotations set during Create are persisted by the reconciler.
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, cr *v1alpha1.Function) (managed.ExternalUpdate, error) {
	ctx = kusto.WithOp(ctx, kind, "update")
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

func (e *external) Delete(ctx context.Context, cr *v1alpha1.Function) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, kind, "delete")
	cr.Status.SetConditions(xpv2.Deleting())
	db := cr.Spec.ForProvider.Database
	_, err := e.kc.Mgmt(ctx, db, function.BuildDelete(meta.GetExternalName(cr)))
	e.invalidate(db)
	if err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

func (e *external) Disconnect(_ context.Context) error { return nil }
