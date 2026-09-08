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

// Package materializedview reconciles MaterializedView managed resources,
// including the asynchronous backfill create.
package materializedview

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
	mv "github.com/functional-team/provider-azure-adx/internal/adx/materializedview"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve   = "cannot observe materialized view"
	errCreate    = "cannot create materialized view"
	errUpdate    = "cannot update materialized view"
	errDelete    = "cannot delete materialized view"
	errOperation = "cannot check backfill operation"

	kind = "MaterializedView"
)

// Setup adds a controller that reconciles MaterializedView managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	k := base.Kind[*v1alpha1.MaterializedView]{GVK: v1alpha1.MaterializedViewGroupVersionKind, Object: &v1alpha1.MaterializedView{}, List: &v1alpha1.MaterializedViewList{}}
	return base.SetupGated(o, k, func() error {
		return base.Register(mgr, o, k, &connector{base: base.NewConnector(mgr, d)})
	})
}

type connector struct{ base *base.Connector }

func (c *connector) Connect(ctx context.Context, cr *v1alpha1.MaterializedView) (managed.TypedExternalClient[*v1alpha1.MaterializedView], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external{kc: kc, cache: c.base.Deps.Cache, kube: c.base.Kube}, nil
}

type external struct {
	kc    kusto.Client
	cache *snapshot.Cache
	kube  client.Client
}

func (e *external) observe(ctx context.Context, db, name string) (mv.Observed, bool, error) {
	if e.cache.Disabled() {
		res, err := e.kc.Mgmt(ctx, db, mv.ShowOne(name))
		if err != nil {
			if kerrors.IsNotFound(err) {
				return mv.Observed{}, false, nil
			}
			return mv.Observed{}, false, err
		}
		o, ok := mv.ParseOne(res, name)
		return o, ok, nil
	}
	all, err := snapshot.Load(ctx, e.cache, snapshot.Key{Endpoint: e.kc.Endpoint(), Database: db, Section: mv.Section}, func(ctx context.Context) (map[string]mv.Observed, error) {
		res, err := e.kc.Mgmt(ctx, db, mv.ShowAll())
		if err != nil {
			return nil, err
		}
		return mv.ParseRows(res), nil
	})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return mv.Observed{}, false, nil
		}
		return mv.Observed{}, false, err
	}
	o, ok := all[name]
	return o, ok, nil
}

func (e *external) invalidate(db string) {
	e.cache.Invalidate(e.kc.Endpoint(), db, mv.Section)
}

// operation looks up an async operation, first in the recent operations and
// then in the historic log form (entries older than about 6 hours disappear
// from the direct lookup).
func (e *external) operation(ctx context.Context, db, id string) (mv.Operation, bool, error) {
	show, err := mv.ShowOperation(id)
	if err != nil {
		return mv.Operation{}, false, err
	}
	res, err := e.kc.Mgmt(ctx, db, show)
	if err != nil && !kerrors.IsNotFound(err) {
		return mv.Operation{}, false, err
	}
	if err == nil {
		if op, ok := mv.ParseOperation(res); ok {
			return op, true, nil
		}
	}
	hist, err := mv.ShowOperationHistoric(id)
	if err != nil {
		return mv.Operation{}, false, err
	}
	res, err = e.kc.Mgmt(ctx, db, hist)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return mv.Operation{}, false, nil
		}
		return mv.Operation{}, false, err
	}
	op, ok := mv.ParseOperation(res)
	return op, ok, nil
}

// clearOperation forgets the tracked operation. The annotation is emptied
// rather than removed so that the merge patch used by PersistAnnotations
// carries the change.
func (e *external) clearOperation(ctx context.Context, cr *v1alpha1.MaterializedView) error {
	meta.AddAnnotations(cr, map[string]string{base.AnnotationOperationID: ""})
	if e.kube == nil {
		return nil
	}
	return base.PersistAnnotations(ctx, e.kube, cr)
}

func (e *external) Observe(ctx context.Context, cr *v1alpha1.MaterializedView) (managed.ExternalObservation, error) { //nolint:gocyclo // async operation states plus the regular observe.
	ctx = kusto.WithOp(ctx, kind, "observe")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	d, err := mv.FromParams(name, cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalObservation{}, kerrors.NewBlocked(kerrors.ReasonInvalidSpec, "%v", err)
	}

	if opID := cr.GetAnnotations()[base.AnnotationOperationID]; opID != "" {
		op, found, err := e.operation(ctx, db, opID)
		if err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, errOperation)
		}
		switch {
		case found && op.Running():
			cr.Status.AtProvider.OperationID = opID
			cr.Status.AtProvider.OperationState = op.State
			cr.Status.SetConditions(xpv2.Creating())
			return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true}, nil
		case found && op.State != mv.StateCompleted:
			if err := e.clearOperation(ctx, cr); err != nil {
				return managed.ExternalObservation{}, errors.Wrap(err, errOperation)
			}
			return managed.ExternalObservation{}, errors.Errorf("backfill operation %s ended with state %s: %s", opID, op.State, op.Status)
		default:
			// Completed, or no longer known to the cluster: continue with a
			// regular observe of the view itself.
			if err := e.clearOperation(ctx, cr); err != nil {
				return managed.ExternalObservation{}, errors.Wrap(err, errOperation)
			}
		}
	}

	obs, ok, err := e.observe(ctx, db, name)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	if !ok {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	queryEqual := base.TextUpToDate(cr, mv.DesiredTexts(d), mv.ObservedTexts(obs))
	plan := mv.Diff(d, obs, queryEqual)
	cr.Status.AtProvider = mv.Observation(obs)
	cr.Status.SetConditions(xpv2.Available())
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: plan.Empty(), Diff: plan.String()}, nil
}

// readBack reads the view after a write and records the query hashes.
func (e *external) readBack(ctx context.Context, cr *v1alpha1.MaterializedView, d mv.Desired) {
	res, err := e.kc.Mgmt(ctx, cr.Spec.ForProvider.Database, mv.ShowOne(d.Name))
	if err != nil {
		return
	}
	if obs, ok := mv.ParseOne(res, d.Name); ok {
		base.SetHashes(cr, mv.DesiredTexts(d), mv.ObservedTexts(obs))
	}
}

func (e *external) Create(ctx context.Context, cr *v1alpha1.MaterializedView) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, kind, "create")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	d, err := mv.FromParams(name, cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	res, err := e.kc.Mgmt(ctx, db, mv.BuildCreate(d))
	e.invalidate(db)
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	if d.Backfill && err == nil {
		if id := mv.ParseOperationID(res); mv.ValidOperationID(id) {
			// Persisted by the reconciler together with the other create annotations.
			meta.AddAnnotations(cr, map[string]string{base.AnnotationOperationID: id})
		}
	}
	e.readBack(ctx, cr, d)
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, cr *v1alpha1.MaterializedView) (managed.ExternalUpdate, error) { //nolint:gocyclo // plan execution with the alter-rejected guardrail.
	ctx = kusto.WithOp(ctx, kind, "update")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	d, err := mv.FromParams(name, cr.Spec.ForProvider)
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
	plan := mv.Diff(d, obs, base.TextUpToDate(cr, mv.DesiredTexts(d), mv.ObservedTexts(obs)))
	altered := false
	for _, s := range plan.Steps {
		if _, err := e.kc.Mgmt(ctx, db, s.Command); err != nil {
			e.invalidate(db)
			if s.Kind == mv.StepAlter && kerrors.Classify(err) == kerrors.Permanent {
				// Kusto refuses e.g. changed group-by expressions. No automatic
				// drop and recreate (decision F8); the message stays in Synced.
				return managed.ExternalUpdate{}, errors.Wrapf(err, "%s: %s", kerrors.ReasonMaterializedViewAlterRejected, errUpdate)
			}
			return managed.ExternalUpdate{}, errors.Wrapf(err, "%s (%s)", errUpdate, s.Kind)
		}
		if s.Kind == mv.StepAlter {
			altered = true
		}
	}
	e.invalidate(db)
	if altered {
		e.readBack(ctx, cr, d)
		if e.kube != nil {
			if err := base.PersistAnnotations(ctx, e.kube, cr); err != nil {
				return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
			}
		}
	}
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, cr *v1alpha1.MaterializedView) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, kind, "delete")
	cr.Status.SetConditions(xpv2.Deleting())
	db := cr.Spec.ForProvider.Database
	_, err := e.kc.Mgmt(ctx, db, mv.BuildDelete(meta.GetExternalName(cr)))
	e.invalidate(db)
	if err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

func (e *external) Disconnect(_ context.Context) error { return nil }
