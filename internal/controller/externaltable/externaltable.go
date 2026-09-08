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

// Package externaltable reconciles ExternalTable managed resources. Connection
// strings with secrets are write-only: they are resolved from Secrets, sent
// obfuscated and tracked through a hash annotation.
package externaltable

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/adx/v1alpha1"
	et "github.com/functional-team/provider-azure-adx/internal/adx/externaltable"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve = "cannot observe external table"
	errCreate  = "cannot create external table"
	errUpdate  = "cannot update external table"
	errDelete  = "cannot delete external table"
	errSecrets = "cannot resolve connection strings"

	kind = "ExternalTable"
)

// Setup adds a controller that reconciles ExternalTable managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	k := base.Kind[*v1alpha1.ExternalTable]{GVK: v1alpha1.ExternalTableGroupVersionKind, Object: &v1alpha1.ExternalTable{}, List: &v1alpha1.ExternalTableList{}}
	return base.SetupGated(o, k, func() error {
		return base.Register(mgr, o, k, &connector{base: base.NewConnector(mgr, d)})
	})
}

type connector struct{ base *base.Connector }

func (c *connector) Connect(ctx context.Context, cr *v1alpha1.ExternalTable) (managed.TypedExternalClient[*v1alpha1.ExternalTable], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external{kc: kc, cache: c.base.Deps.Cache, kube: c.base.Kube}, nil
}

// NewExternal builds the external client directly (integration tests).
func NewExternal(kc kusto.Client, cache *snapshot.Cache, kube client.Client) managed.TypedExternalClient[*v1alpha1.ExternalTable] {
	return &external{kc: kc, cache: cache, kube: kube}
}

type external struct {
	kc    kusto.Client
	cache *snapshot.Cache
	kube  client.Client
}

// connectionStrings resolves the spec's connection strings; secret values are
// read from Secrets and must never leave this process except obfuscated.
func (e *external) connectionStrings(ctx context.Context, cr *v1alpha1.ExternalTable) ([]string, error) {
	out := make([]string, 0, len(cr.Spec.ForProvider.ConnectionStrings))
	for i, cs := range cr.Spec.ForProvider.ConnectionStrings {
		switch {
		case cs.Value != nil:
			out = append(out, strings.TrimSpace(*cs.Value))
		case cs.SecretKeyRef != nil:
			if e.kube == nil {
				return nil, errors.Errorf("connectionStrings[%d]: no Kubernetes client to read the secret", i)
			}
			ns := cs.SecretKeyRef.Namespace
			if ns == "" {
				ns = cr.GetNamespace()
			}
			s := &corev1.Secret{}
			if err := e.kube.Get(ctx, types.NamespacedName{Namespace: ns, Name: cs.SecretKeyRef.Name}, s); err != nil {
				return nil, errors.Wrapf(err, "connectionStrings[%d]: cannot get secret %s/%s", i, ns, cs.SecretKeyRef.Name)
			}
			v, ok := s.Data[cs.SecretKeyRef.Key]
			if !ok || len(v) == 0 {
				return nil, errors.Errorf("connectionStrings[%d]: secret %s/%s has no key %q", i, ns, cs.SecretKeyRef.Name, cs.SecretKeyRef.Key)
			}
			out = append(out, strings.TrimSpace(string(v)))
		default:
			return nil, errors.Errorf("connectionStrings[%d]: value or secretKeyRef is required", i)
		}
	}
	return out, nil
}

func (e *external) observe(ctx context.Context, db, name string, withSchema bool) (et.Observed, bool, error) { //nolint:gocyclo // cache/no-cache plus optional schema lookup.
	var o et.Observed
	var ok bool
	if e.cache.Disabled() {
		res, err := e.kc.Mgmt(ctx, db, et.Show(name))
		if err != nil {
			if kerrors.IsNotFound(err) {
				return et.Observed{}, false, nil
			}
			return et.Observed{}, false, err
		}
		o, ok = et.ParseOne(res, name)
	} else {
		all, err := snapshot.Load(ctx, e.cache, snapshot.Key{Endpoint: e.kc.Endpoint(), Database: db, Section: et.Section}, func(ctx context.Context) (map[string]et.Observed, error) {
			res, err := e.kc.Mgmt(ctx, db, et.ShowAll())
			if err != nil {
				return nil, err
			}
			return et.ParseRows(res), nil
		})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return et.Observed{}, false, nil
			}
			return et.Observed{}, false, err
		}
		o, ok = all[name]
	}
	if !ok {
		return et.Observed{}, false, nil
	}
	if withSchema {
		res, err := e.kc.Mgmt(ctx, db, et.ShowCslSchema(name))
		if err != nil {
			if kerrors.IsNotFound(err) {
				return et.Observed{}, false, nil
			}
			return et.Observed{}, false, err
		}
		if cols, ok := et.ParseCslSchema(res); ok {
			o.Columns, o.ColumnsLoaded = cols, true
		}
	}
	return o, true, nil
}

func (e *external) invalidate(db string) {
	e.cache.Invalidate(e.kc.Endpoint(), db, et.Section)
}

func (e *external) Observe(ctx context.Context, cr *v1alpha1.ExternalTable) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, kind, "observe")
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	css, err := e.connectionStrings(ctx, cr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errSecrets)
	}
	d := et.FromParams(name, cr.Spec.ForProvider, css)
	obs, ok, err := e.observe(ctx, db, name, len(d.Columns) > 0)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	if !ok {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	diff := et.Compare(d, obs)
	reasons := append([]string(nil), diff.Reasons...)
	if !base.TextUpToDate(cr, et.DesiredTexts(d), et.ObservedTexts(obs)) {
		reasons = append(reasons, "partitionBy or pathFormat differ")
	}
	if cr.GetAnnotations()[base.AnnotationSecretHash] != et.SecretHash(css) {
		reasons = append(reasons, "connection strings changed")
	}
	cr.Status.AtProvider = et.Observation(obs)
	cr.Status.SetConditions(xpv2.Available())
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: len(reasons) == 0, Diff: strings.Join(reasons, "; ")}, nil
}

// write runs create-or-alter, reads the table back and records the hashes.
func (e *external) write(ctx context.Context, cr *v1alpha1.ExternalTable) error {
	name, db := meta.GetExternalName(cr), cr.Spec.ForProvider.Database
	css, err := e.connectionStrings(ctx, cr)
	if err != nil {
		return errors.Wrap(err, errSecrets)
	}
	d := et.FromParams(name, cr.Spec.ForProvider, css)
	_, err = e.kc.Mgmt(ctx, db, et.BuildCreateOrAlter(d))
	e.invalidate(db)
	if err != nil {
		return err
	}
	meta.AddAnnotations(cr, map[string]string{base.AnnotationSecretHash: et.SecretHash(css)})
	res, err := e.kc.Mgmt(ctx, db, et.Show(name))
	if err != nil {
		return errors.Wrap(err, "cannot read external table back after write")
	}
	if obs, ok := et.ParseOne(res, name); ok {
		base.SetHashes(cr, et.DesiredTexts(d), et.ObservedTexts(obs))
	}
	return nil
}

func (e *external) Create(ctx context.Context, cr *v1alpha1.ExternalTable) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, kind, "create")
	if err := e.write(ctx, cr); err != nil && !kerrors.IsAlreadyExists(err) {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, cr *v1alpha1.ExternalTable) (managed.ExternalUpdate, error) {
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

func (e *external) Delete(ctx context.Context, cr *v1alpha1.ExternalTable) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, kind, "delete")
	cr.Status.SetConditions(xpv2.Deleting())
	db := cr.Spec.ForProvider.Database
	_, err := e.kc.Mgmt(ctx, db, et.BuildDelete(meta.GetExternalName(cr)))
	e.invalidate(db)
	if err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

func (e *external) Disconnect(_ context.Context) error { return nil }
