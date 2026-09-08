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

// Package workloadgroup holds the domain logic for the WorkloadGroup kind.
package workloadgroup

import (
	"encoding/json"
	"strings"

	adxpolicy "github.com/functional-team/provider-azure-adx/internal/adx/policy"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
)

// Built-in workload groups that exist on every cluster and cannot be dropped.
const (
	BuiltInDefault  = "default"
	BuiltInInternal = "internal"
)

// IsBuiltIn reports whether name is a workload group Kusto refuses to drop.
func IsBuiltIn(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case BuiltInDefault, BuiltInInternal:
		return true
	}
	return false
}

// ShowCmd is ".show workload_group ['WG']".
func ShowCmd(name string) cmd.Command {
	return cmd.New(".show workload_group ", cmd.Ident(name))
}

// ShowAllCmd is ".show workload_groups".
func ShowAllCmd() cmd.Command { return cmd.New(".show workload_groups") }

// CreateOrAlterCmd is ".create-or-alter workload_group ['WG'] @'<json>'".
func CreateOrAlterCmd(name string, desired any) (cmd.Command, error) {
	js, err := cmd.JSON(desired)
	if err != nil {
		return cmd.Command{}, err
	}
	return cmd.New(".create-or-alter workload_group ", cmd.Ident(name), " ", js), nil
}

// DropCmd is ".drop workload_group ['WG']".
func DropCmd(name string) cmd.Command {
	return cmd.New(".drop workload_group ", cmd.Ident(name))
}

// ParseShow returns the workload group JSON for name from a ".show
// workload_group(s)" result (columns WorkloadGroupName, WorkloadGroup).
func ParseShow(res *kusto.Result, name string) (json.RawMessage, bool) {
	rows := res.Rows()
	for _, r := range rows {
		if r.String("WorkloadGroupName") == name {
			raw := r.Dynamic("WorkloadGroup")
			return raw, !adxpolicy.IsNull(raw)
		}
	}
	// Single-row result without a name column, or a renamed row.
	if len(rows) == 1 && !rows[0].Has("WorkloadGroupName") {
		raw := rows[0].Dynamic("WorkloadGroup")
		return raw, !adxpolicy.IsNull(raw)
	}
	return nil, false
}

// ParseShowAll returns name -> workload group JSON.
func ParseShowAll(res *kusto.Result) map[string]json.RawMessage {
	out := map[string]json.RawMessage{}
	for _, r := range res.Rows() {
		if n := r.String("WorkloadGroupName"); n != "" {
			out[n] = r.Dynamic("WorkloadGroup")
		}
	}
	return out
}

// Compare checks that every key of the desired JSON equals the observed one.
func Compare(desired any, observed json.RawMessage) (adxpolicy.Result, error) {
	return adxpolicy.Compare(desired, observed, nil, nil)
}
