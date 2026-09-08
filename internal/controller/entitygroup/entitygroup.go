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

// Package entitygroup reconciles EntityGroup managed resources (Tier 3).
package entitygroup

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
	eg "github.com/functional-team/provider-azure-adx/internal/adx/entitygroup"
	"github.com/functional-team/provider-azure-adx/internal/adx/schema"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve = "cannot observe entity group"
	errCreate  = "cannot create entity group"
	errUpdate  = "cannot update entity group"
	errDelete  = "cannot delete entity group"

	kind = "EntityGroup"
)

// Setup adds a controller that reconciles EntityGroup managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	k := base.Kind[*v1alpha1.EntityGroup]{GVK: v1alpha1.EntityGroupGroupVersionKind, Object: &v1alpha1.EntityGroup{}, List: &v1alpha1.EntityGroupList{}}
	return base.SetupGated(o, k, func() error {
		return base.Register(mgr, o, k, &connector{base: base.NewConnector(mgr, d)})
	})
}

type connector struct{ base *base.Connector }

func (c *connector) Connect(ctx context.Context, cr *v1alpha1.EntityGroup) (managed.TypedExternalClient[*v1alpha1.EntityGroup], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external{kc: kc, cache: c.base.Deps.Cache, kube: c.base.Kube}, nil
}

// NewExternal builds the external client directly (integration tests).
func NewExternal(kc kusto.Client, cache *snapshot.Cache, kube client.Client) managed.TypedExternalClient[*v1alpha1.EntityGroup] {
	return &external{kc: kc, cache: cache, kube: kube}
}

type external struct {
	kc    kusto.Client
	cache *snapshot.Cache
	kube  client.Client
}

func (e *external) observe(ctx context.Context, db, name string) (eg.Observed, bool, error) {
	if e.cache.Disabled() {
		res, err := e.kc.Mgmt(ctx, db, eg.ShowOne(name))
		if err != nil {
			if kerrors.IsNotFound(err) {
				return eg.Observed{}, false, nil
			}
			return eg.Observed{}, false, err
		}
		o, ok := eg.ParseShow(res, name)
		return o, ok, nil
	}
	d, err := schema.Cached(ctx, e.cache, e.kc, db)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return eg.Observed{}, false, nil
		}
		return eg.Observed{}, false, err
	}
	g, ok := d.EntityGroups[name]
	if !ok {
		return eg.Observed{}, false, nil
	}
	return eg.FromSchema(g), true, nil
}

func (e *external) Observe(ctx context.Context, cr *v1alpha1.EntityGroup) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, kind, "observe")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	obs, ok, err := e.observe(ctx, db, name)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	if !ok {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	desired := cr.Spec.ForProvider.Entities
	cr.Status.AtProvider = eg.Observation(obs)
	cr.Status.SetConditions(xpv2.Available())
	upToDate := eg.Equal(desired, obs.Entities) || base.TextUpToDate(cr, eg.Texts(desired), eg.Texts(obs.Entities))
	diff := ""
	if !upToDate {
		diff = "entities differ"
	}
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: upToDate, Diff: diff}, nil
}

// write creates or alters the group and records the hashes from a read back.
func (e *external) write(ctx context.Context, cr *v1alpha1.EntityGroup) error {
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	desired := cr.Spec.ForProvider.Entities
	_, err := e.kc.Mgmt(ctx, db, eg.BuildCreateOrAlter(name, desired))
	schema.Invalidate(e.cache, e.kc, db)
	if err != nil {
		return err
	}
	res, err := e.kc.Mgmt(ctx, db, eg.ShowOne(name))
	if err != nil {
		// The read back only feeds the hash tolerance; a cluster that does not
		// support the single show form still converges through Equal().
		return nil //nolint:nilerr // intentional: hashes are optional
	}
	if obs, ok := eg.ParseShow(res, name); ok {
		base.SetHashes(cr, eg.Texts(desired), eg.Texts(obs.Entities))
	}
	return nil
}

func (e *external) Create(ctx context.Context, cr *v1alpha1.EntityGroup) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, kind, "create")
	if err := e.write(ctx, cr); err != nil && !kerrors.IsAlreadyExists(err) {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, cr *v1alpha1.EntityGroup) (managed.ExternalUpdate, error) {
	ctx = kusto.WithOp(ctx, kind, "update")
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

func (e *external) Delete(ctx context.Context, cr *v1alpha1.EntityGroup) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, kind, "delete")
	cr.Status.SetConditions(xpv2.Deleting())
	db := cr.Spec.ForProvider.Database
	_, err := e.kc.Mgmt(ctx, db, eg.BuildDelete(meta.GetExternalName(cr)))
	schema.Invalidate(e.cache, e.kc, db)
	if err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

func (e *external) Disconnect(_ context.Context) error { return nil }
