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

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"

	"github.com/functional-team/provider-azure-adx/internal/controller/base"
	"github.com/functional-team/provider-azure-adx/internal/controller/clusterpolicy"
	"github.com/functional-team/provider-azure-adx/internal/controller/config"
	"github.com/functional-team/provider-azure-adx/internal/controller/continuousexport"
	"github.com/functional-team/provider-azure-adx/internal/controller/entitygroup"
	"github.com/functional-team/provider-azure-adx/internal/controller/externaltable"
	"github.com/functional-team/provider-azure-adx/internal/controller/function"
	"github.com/functional-team/provider-azure-adx/internal/controller/ingestionmapping"
	"github.com/functional-team/provider-azure-adx/internal/controller/materializedview"
	"github.com/functional-team/provider-azure-adx/internal/controller/policy"
	"github.com/functional-team/provider-azure-adx/internal/controller/securityrole"
	"github.com/functional-team/provider-azure-adx/internal/controller/table"
	"github.com/functional-team/provider-azure-adx/internal/controller/workloadgroup"
)

// SetupGated creates all provider controllers with safe-start support and adds
// them to the supplied manager.
func SetupGated(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	if err := config.Setup(mgr, o); err != nil {
		return err
	}
	for _, setup := range []func(ctrl.Manager, controller.Options, base.Deps) error{
		table.Setup,
		function.Setup,
		materializedview.Setup,
		externaltable.Setup,
		continuousexport.Setup,
		ingestionmapping.Setup,
		entitygroup.Setup,
		policy.SetupAll,
		securityrole.Setup,
		clusterpolicy.SetupAll,
		workloadgroup.Setup,
	} {
		if err := setup(mgr, o, d); err != nil {
			return err
		}
	}
	return nil
}
