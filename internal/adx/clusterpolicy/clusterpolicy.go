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

// Package clusterpolicy holds the domain logic shared by the cluster-level
// policy kinds: command shapes (".show/.alter/.delete cluster policy <name>")
// and the comparison, which reuses the database policy comparison so both
// levels behave the same ("only compare what the spec sets").
package clusterpolicy

import (
	"encoding/json"

	adxpolicy "github.com/functional-team/provider-azure-adx/internal/adx/policy"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
)

// Def describes one Kusto cluster policy type.
type Def struct {
	// Name is the policy name used in commands, e.g. "callout".
	Name string
	// Alter renders the alter command. nil means the JSON form
	// ".alter cluster policy <name> @'<json>'", or ".alter-merge" when Merge
	// is set.
	Alter func(desired any) (cmd.Command, error)
	// KQLFields are JSON keys whose string values are KQL text (two stage
	// comparison instead of exact equality).
	KQLFields []string
	// SetFields are JSON keys holding comma separated sets ("A, B").
	SetFields []string
	// HashOnly compares desired and observed only through the stage 2 hashes
	// (policies whose .show JSON shape is not verified).
	HashOnly bool
	// NoDelete marks policies Kusto cannot delete (capacity, query weak
	// consistency): deleting the managed resource leaves the policy as is.
	NoDelete bool
	// Merge marks policies Kusto only accepts via ".alter-merge". The capacity
	// policy is one: ".alter" is rejected with "The '.alter' command is not
	// supported for CapacityPolicy. Please use '.alter-merge' instead."
	// (observed against a real cluster on 2026-09-08).
	Merge bool
}

// ShowCmd is ".show cluster policy <name>".
func (d Def) ShowCmd() cmd.Command {
	return cmd.New(".show cluster policy ", d.Name)
}

// AlterCmd renders the alter command for desired.
func (d Def) AlterCmd(desired any) (cmd.Command, error) {
	if d.Alter != nil {
		return d.Alter(desired)
	}
	js, err := cmd.JSON(desired)
	if err != nil {
		return cmd.Command{}, err
	}
	return cmd.New(d.alterVerb(), " cluster policy ", d.Name, " ", js), nil
}

func (d Def) alterVerb() string {
	if d.Merge {
		return ".alter-merge"
	}
	return ".alter"
}

// DeleteCmd is ".delete cluster policy <name>"; ok is false for policies
// that cannot be deleted.
func (d Def) DeleteCmd() (cmd.Command, bool) {
	if d.NoDelete {
		return cmd.Command{}, false
	}
	return cmd.New(".delete cluster policy ", d.Name), true
}

// CompareWith compares desired against the observed policy JSON.
func (d Def) CompareWith(desired any, observed json.RawMessage) (adxpolicy.Result, error) {
	return adxpolicy.Def{Name: d.Name, KQLFields: d.KQLFields, SetFields: d.SetFields, HashOnly: d.HashOnly}.CompareWith(desired, observed)
}

// ParseShow returns the Policy JSON of a ".show cluster policy" result; ok is
// false when there is no row or the policy is null.
func ParseShow(res *kusto.Result) (json.RawMessage, bool) {
	return adxpolicy.ParseShow(res)
}
