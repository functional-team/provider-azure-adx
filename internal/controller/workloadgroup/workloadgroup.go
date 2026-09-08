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

// Package workloadgroup reconciles WorkloadGroup managed resources.
package workloadgroup

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/functional-team/provider-azure-adx/apis/cluster/v1alpha1"
	adxpolicy "github.com/functional-team/provider-azure-adx/internal/adx/policy"
	"github.com/functional-team/provider-azure-adx/internal/adx/workloadgroup"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve = "cannot observe workload group"
	errCreate  = "cannot create workload group"
	errUpdate  = "cannot update workload group"
	errDelete  = "cannot delete workload group"

	kind = "WorkloadGroup"
)

// Setup adds a controller that reconciles WorkloadGroup managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	k := base.Kind[*v1alpha1.WorkloadGroup]{GVK: v1alpha1.WorkloadGroupGroupVersionKind, Object: &v1alpha1.WorkloadGroup{}, List: &v1alpha1.WorkloadGroupList{}}
	return base.SetupGated(o, k, func() error {
		return base.Register(mgr, o, k, &connector{base: base.NewConnector(mgr, d)})
	})
}

type connector struct{ base *base.Connector }

func (c *connector) Connect(ctx context.Context, cr *v1alpha1.WorkloadGroup) (managed.TypedExternalClient[*v1alpha1.WorkloadGroup], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external{kc: kc}, nil
}

type external struct {
	kc kusto.Client
}

func desired(cr *v1alpha1.WorkloadGroup) (any, error) {
	raw := strings.TrimSpace(string(cr.Spec.ForProvider.WorkloadGroup.Raw))
	if raw == "" || raw == "null" {
		return nil, fmt.Errorf("workloadGroup is empty")
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return nil, fmt.Errorf("workloadGroup is not valid JSON: %w", err)
	}
	if _, ok := v.(map[string]any); !ok {
		return nil, fmt.Errorf("workloadGroup must be a JSON object")
	}
	return v, nil
}

func (e *external) observe(ctx context.Context, name string) (json.RawMessage, bool, error) {
	res, err := e.kc.Mgmt(ctx, "", workloadgroup.ShowCmd(name))
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	raw, ok := workloadgroup.ParseShow(res, name)
	return raw, ok, nil
}

func (e *external) Observe(ctx context.Context, cr *v1alpha1.WorkloadGroup) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, kind, "observe")
	name := meta.GetExternalName(cr)
	d, err := desired(cr)
	if err != nil {
		return managed.ExternalObservation{}, kerrors.NewBlocked(kerrors.ReasonInvalidSpec, "%v", err)
	}
	raw, ok, err := e.observe(ctx, name)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	if !ok {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	res, err := workloadgroup.Compare(d, raw)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	cr.Status.AtProvider = v1alpha1.WorkloadGroupObservation{WorkloadGroup: adxpolicy.Compact(raw)}
	cr.Status.SetConditions(xpv2.Available())
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: res.Equal, Diff: res.Diff}, nil
}

func (e *external) write(ctx context.Context, cr *v1alpha1.WorkloadGroup) error {
	d, err := desired(cr)
	if err != nil {
		return kerrors.NewBlocked(kerrors.ReasonInvalidSpec, "%v", err)
	}
	c, err := workloadgroup.CreateOrAlterCmd(meta.GetExternalName(cr), d)
	if err != nil {
		return err
	}
	_, err = e.kc.Mgmt(ctx, "", c)
	return err
}

func (e *external) Create(ctx context.Context, cr *v1alpha1.WorkloadGroup) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, kind, "create")
	if err := e.write(ctx, cr); err != nil && !kerrors.IsAlreadyExists(err) {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, cr *v1alpha1.WorkloadGroup) (managed.ExternalUpdate, error) {
	ctx = kusto.WithOp(ctx, kind, "update")
	if err := e.write(ctx, cr); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdate)
	}
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, cr *v1alpha1.WorkloadGroup) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, kind, "delete")
	cr.Status.SetConditions(xpv2.Deleting())
	name := meta.GetExternalName(cr)
	if workloadgroup.IsBuiltIn(name) {
		// "default" and "internal" exist on every cluster and cannot be
		// dropped; the last applied settings stay in place.
		return managed.ExternalDelete{}, nil
	}
	if _, err := e.kc.Mgmt(ctx, "", workloadgroup.DropCmd(name)); err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

func (e *external) Disconnect(_ context.Context) error { return nil }
