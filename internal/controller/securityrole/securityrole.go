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

// Package securityrole reconciles SecurityRole managed resources. The hard
// part, mapping what the user wrote onto the object ids Kusto reports, lives
// in internal/adx/securityrole; this package runs the commands and keeps the
// mapping in the status.
package securityrole

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

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/security/v1alpha1"
	sr "github.com/functional-team/provider-azure-adx/internal/adx/securityrole"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/kerrors"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
)

const (
	errObserve = "cannot observe security role"
	errCreate  = "cannot create security role"
	errUpdate  = "cannot update security role"
	errDelete  = "cannot delete security role"

	kind = "SecurityRole"

	// AnnotationOwned marks that this resource wrote the role assignment; an
	// empty authoritative role then still counts as existing.
	AnnotationOwned = "adx.functional.team/role-owned"
	// AnnotationBeforeIDs lists the principal ids of the role before the last
	// write; it lets Observe attribute a single new row to an unresolvable
	// spec entry.
	AnnotationBeforeIDs = "adx.functional.team/role-before-ids"
)

// Setup adds a controller that reconciles SecurityRole managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	k := base.Kind[*v1alpha1.SecurityRole]{GVK: v1alpha1.SecurityRoleGroupVersionKind, Object: &v1alpha1.SecurityRole{}, List: &v1alpha1.SecurityRoleList{}}
	return base.SetupGated(o, k, func() error {
		return base.Register(mgr, o, k, &connector{base: base.NewConnector(mgr, d)})
	})
}

type connector struct{ base *base.Connector }

func (c *connector) Connect(ctx context.Context, cr *v1alpha1.SecurityRole) (managed.TypedExternalClient[*v1alpha1.SecurityRole], error) {
	kc, err := c.base.Connect(ctx, cr)
	if err != nil {
		return nil, err
	}
	return &external{kc: kc, kube: c.base.Kube}, nil
}

// NewExternal builds the external client directly (integration tests).
func NewExternal(kc kusto.Client, kube client.Client) managed.TypedExternalClient[*v1alpha1.SecurityRole] {
	return &external{kc: kc, kube: kube}
}

type external struct {
	kc   kusto.Client
	kube client.Client
}

func entity(cr *v1alpha1.SecurityRole) (sr.Entity, error) {
	p := cr.Spec.ForProvider
	ent := sr.Entity{Kind: p.Entity.Kind, Database: p.Database}
	if p.Entity.Kind != common.EntityKindDatabase {
		if p.Entity.Name == nil || *p.Entity.Name == "" {
			return ent, errors.New("entity.name is empty; set it or wait for entity.nameRef to resolve")
		}
		ent.Name = *p.Entity.Name
	}
	return ent, nil
}

func principals(cr *v1alpha1.SecurityRole) []string {
	out := make([]string, 0, len(cr.Spec.ForProvider.Principals))
	for _, p := range cr.Spec.ForProvider.Principals {
		if s := strings.TrimSpace(string(p)); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func additive(cr *v1alpha1.SecurityRole) bool {
	return cr.Spec.ForProvider.Mode == v1alpha1.ModeAdditive
}

// rows reads the principals of the entity filtered to the role.
func (e *external) rows(ctx context.Context, ent sr.Entity, role string) ([]sr.Row, error) {
	show, err := sr.ShowPrincipals(ent)
	if err != nil {
		return nil, err
	}
	res, err := e.kc.Mgmt(ctx, ent.Database, show)
	if err != nil {
		return nil, err
	}
	return sr.FilterRole(sr.ParsePrincipals(res), role), nil
}

func splitIDs(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func (e *external) evaluate(cr *v1alpha1.SecurityRole, rows []sr.Row) sr.Evaluation {
	ann := cr.GetAnnotations()
	applied, _ := base.Hashes(cr)
	spec := principals(cr)
	return sr.Evaluate(sr.Input{
		Additive:      additive(cr),
		Spec:          spec,
		Resolved:      sr.FromObservation(cr.Status.AtProvider),
		Rows:          rows,
		BeforeIDs:     splitIDs(ann[AnnotationBeforeIDs]),
		SpecUnchanged: applied != "" && applied == sr.SpecHash(spec),
		Owned:         ann[AnnotationOwned] == "true",
	})
}

func (e *external) Observe(ctx context.Context, cr *v1alpha1.SecurityRole) (managed.ExternalObservation, error) {
	ctx = kusto.WithOp(ctx, kind, "observe")
	ent, err := entity(cr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	rows, err := e.rows(ctx, ent, string(cr.Spec.ForProvider.Role))
	if err != nil {
		if kerrors.IsNotFound(err) {
			return managed.ExternalObservation{ResourceExists: false}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errObserve)
	}
	ev := e.evaluate(cr, rows)
	cr.Status.AtProvider = sr.Observation(ev, rows)
	if !ev.Exists {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	cr.Status.SetConditions(xpv2.Available())
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: ev.UpToDate, Diff: ev.Diff}, nil
}

// write applies the spec: Authoritative sets the whole list, Additive adds
// missing and drops stale principals. The pre-write ids and the spec hash are
// recorded as annotations so Observe can resolve the newly added principals.
func (e *external) write(ctx context.Context, cr *v1alpha1.SecurityRole) error { //nolint:gocyclo // authoritative vs additive write plan.
	ent, err := entity(cr)
	if err != nil {
		return err
	}
	role := string(cr.Spec.ForProvider.Role)
	spec := principals(cr)
	before, err := e.rows(ctx, ent, role)
	if err != nil && !kerrors.IsNotFound(err) {
		return err
	}
	var cmds []func() error
	if additive(cr) {
		ev := e.evaluate(cr, before)
		if len(ev.ToAdd) > 0 {
			c, err := sr.BuildAdd(ent, role, ev.ToAdd, cr.Spec.ForProvider.Description)
			if err != nil {
				return err
			}
			cmds = append(cmds, func() error { _, err := e.kc.Mgmt(ctx, ent.Database, c); return err })
		}
		if len(ev.ToDrop) > 0 {
			c, err := sr.BuildDrop(ent, role, ev.ToDrop)
			if err != nil {
				return err
			}
			cmds = append(cmds, func() error { _, err := e.kc.Mgmt(ctx, ent.Database, c); return err })
		}
	} else {
		c, err := sr.BuildSet(ent, role, spec, cr.Spec.ForProvider.Description)
		if err != nil {
			return err
		}
		cmds = append(cmds, func() error { _, err := e.kc.Mgmt(ctx, ent.Database, c); return err })
	}
	for _, run := range cmds {
		if err := run(); err != nil {
			return err
		}
	}
	meta.AddAnnotations(cr, map[string]string{
		AnnotationOwned:            "true",
		AnnotationBeforeIDs:        strings.Join(sr.IDs(before), ","),
		base.AnnotationAppliedHash: sr.SpecHash(spec),
	})
	return nil
}

func (e *external) Create(ctx context.Context, cr *v1alpha1.SecurityRole) (managed.ExternalCreation, error) {
	ctx = kusto.WithOp(ctx, kind, "create")
	if err := e.write(ctx, cr); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreate)
	}
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, cr *v1alpha1.SecurityRole) (managed.ExternalUpdate, error) {
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

// Delete removes what this resource manages: Authoritative sets the role to
// none, Additive drops only the principals it resolved.
func (e *external) Delete(ctx context.Context, cr *v1alpha1.SecurityRole) (managed.ExternalDelete, error) {
	ctx = kusto.WithOp(ctx, kind, "delete")
	cr.Status.SetConditions(xpv2.Deleting())
	ent, err := entity(cr)
	if err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	role := string(cr.Spec.ForProvider.Role)
	var command cmd.Command
	if additive(cr) {
		own := ownPrincipals(cr)
		if len(own) == 0 {
			return managed.ExternalDelete{}, nil
		}
		command, err = sr.BuildDrop(ent, role, own)
	} else {
		command, err = sr.BuildSet(ent, role, nil, nil)
	}
	if err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	if _, err := e.kc.Mgmt(ctx, ent.Database, command); err != nil && !kerrors.IsNotFound(err) {
		return managed.ExternalDelete{}, errors.Wrap(err, errDelete)
	}
	return managed.ExternalDelete{}, nil
}

// ownPrincipals returns the principals this resource resolved (FQN, or the
// spec string when Kusto reported no FQN).
func ownPrincipals(cr *v1alpha1.SecurityRole) []string {
	var own []string
	for _, r := range cr.Status.AtProvider.ResolvedPrincipals {
		switch {
		case r.FQN != "":
			own = append(own, r.FQN)
		case r.Spec != "":
			own = append(own, r.Spec)
		}
	}
	return own
}

func (e *external) Disconnect(_ context.Context) error { return nil }
