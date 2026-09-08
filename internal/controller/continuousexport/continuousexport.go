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

// Package continuousexport reconciles ContinuousExport managed resources.
package continuousexport

import (
	"context"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	ce "github.com/functional-team/provider-azure-adx/internal/adx/continuousexport"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve = "cannot observe continuous export"
	errCreate  = "cannot create continuous export"
	errUpdate  = "cannot update continuous export"
	errDelete  = "cannot delete continuous export"

	kind = "ContinuousExport"
)

// Setup adds a controller that reconciles ContinuousExport managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	k := base.Kind[*v1alpha1.ContinuousExport]{GVK: v1alpha1.ContinuousExportGroupVersionKind, Object: &v1alpha1.ContinuousExport{}, List: &v1alpha1.ContinuousExportList{}}
	return base.SetupGated(o, k, func() error {
		return base.Register(mgr, o, k, &connector{base: base.NewConnector(mgr, d)})
	})
}

type connector struct{ base *base.Connector }

func (c *connector) Connect(ctx context.Context, cr *v1alpha1.ContinuousExport) (managed.TypedExternalClient[*v1alpha1.ContinuousExport], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external{kc: kc, cache: c.base.Deps.Cache, kube: c.base.Kube}, nil
}

// NewExternal builds the external client directly (integration tests).
func NewExternal(kc kusto.Client, cache *snapshot.Cache, kube client.Client) managed.TypedExternalClient[*v1alpha1.ContinuousExport] {
	return &external{kc: kc, cache: cache, kube: kube}
}

type external struct {
	kc    kusto.Client
	cache *snapshot.Cache
	kube  client.Client
}

func (e *external) observe(ctx context.Context, db, name string) (ce.Observed, bool, error) {
	if e.cache.Disabled() {
		res, err := e.kc.Mgmt(ctx, db, ce.Show(name))
		if err != nil {
			if kerrors.IsNotFound(err) {
				return ce.Observed{}, false, nil
			}
			return ce.Observed{}, false, err
		}
		o, ok := ce.ParseOne(res, name)
		return o, ok, nil
	}
	all, err := snapshot.Load(ctx, e.cache, snapshot.Key{Endpoint: e.kc.Endpoint(), Database: db, Section: ce.Section}, func(ctx context.Context) (map[string]ce.Observed, error) {
		res, err := e.kc.Mgmt(ctx, db, ce.ShowAll())
		if err != nil {
			return nil, err
		}
		return ce.ParseRows(res), nil
	})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return ce.Observed{}, false, nil
		}
		return ce.Observed{}, false, err
	}
	o, ok := all[name]
	return o, ok, nil
}

func (e *external) invalidate(db string) {
	e.cache.Invalidate(e.kc.Endpoint(), db, ce.Section)
}

func (e *external) Observe(ctx context.Context, cr *v1alpha1.ContinuousExport) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, kind, "observe")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	d, err := ce.FromParams(name, cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalObservation{}, kerrors.NewBlocked(kerrors.ReasonInvalidSpec, "%v", err)
	}
	obs, ok, err := e.observe(ctx, db, name)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	if !ok {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	plan := ce.Diff(d, obs)
	reasons := append([]string(nil), plan.Reasons...)
	if !base.TextUpToDate(cr, ce.DesiredTexts(d), ce.ObservedTexts(obs)) {
		reasons = append(reasons, "query differs")
	}
	cr.Status.AtProvider = ce.Observation(obs)
	cr.Status.SetConditions(xpv2.Available())
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: len(reasons) == 0, Diff: strings.Join(reasons, "; ")}, nil
}

// write sends create-or-alter and/or the enable/disable toggle, then reads the
// export back to record the query hashes.
func (e *external) write(ctx context.Context, cr *v1alpha1.ContinuousExport, d ce.Desired, alter bool, toggle *bool) error {
	db := cr.Spec.ForProvider.Database
	if alter {
		if _, err := e.kc.Mgmt(ctx, db, ce.BuildCreateOrAlter(d)); err != nil {
			e.invalidate(db)
			return err
		}
	}
	if toggle != nil {
		if _, err := e.kc.Mgmt(ctx, db, ce.BuildToggle(d.Name, *toggle)); err != nil {
			e.invalidate(db)
			return err
		}
	}
	e.invalidate(db)
	if !alter {
		return nil
	}
	res, err := e.kc.Mgmt(ctx, db, ce.Show(d.Name))
	if err != nil {
		return errors.Wrap(err, "cannot read continuous export back after write")
	}
	if obs, ok := ce.ParseOne(res, d.Name); ok {
		base.SetHashes(cr, ce.DesiredTexts(d), ce.ObservedTexts(obs))
	}
	return nil
}

func (e *external) Create(ctx context.Context, cr *v1alpha1.ContinuousExport) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, kind, "create")
	d, err := ce.FromParams(meta.GetExternalName(cr), cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	// Kusto creates exports enabled; disable right after when the spec says so.
	var toggle *bool
	if !d.Enabled {
		off := false
		toggle = &off
	}
	if err := e.write(ctx, cr, d, true, toggle); err != nil && !kerrors.IsAlreadyExists(err) {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, cr *v1alpha1.ContinuousExport) (managed.ExternalUpdate, error) {
	ctx = kusto.WithOp(ctx, kind, "update")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	d, err := ce.FromParams(name, cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}
	obs, ok, err := e.observe(ctx, db, name)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}
	if !ok {
		return managed.ExternalUpdate{}, nil // gone; the next Observe recreates it
	}
	plan := ce.Diff(d, obs)
	alter := plan.Alter || !base.TextUpToDate(cr, ce.DesiredTexts(d), ce.ObservedTexts(obs))
	if err := e.write(ctx, cr, d, alter, plan.Toggle); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}
	if alter && e.kube != nil {
		if err := base.PersistAnnotations(ctx, e.kube, cr); err != nil {
			return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
		}
	}
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, cr *v1alpha1.ContinuousExport) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, kind, "delete")
	cr.Status.SetConditions(xpv2.Deleting())
	db := cr.Spec.ForProvider.Database
	_, err := e.kc.Mgmt(ctx, db, ce.BuildDelete(meta.GetExternalName(cr)))
	e.invalidate(db)
	if err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

func (e *external) Disconnect(_ context.Context) error { return nil }
