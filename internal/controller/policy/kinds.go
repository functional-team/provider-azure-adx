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

package policy

import (
	"fmt"
	"strconv"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"

	"github.com/functional-team/provider-azure-adx/apis/common"
	"github.com/functional-team/provider-azure-adx/apis/policy/v1alpha1"
	adxpolicy "github.com/functional-team/provider-azure-adx/internal/adx/policy"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/controller/base"
	"github.com/functional-team/provider-azure-adx/internal/timespan"
)

// SetupAll adds the controllers of every policy kind.
func SetupAll(mgr ctrl.Manager, o controller.Options, d base.Deps) error {
	for _, setup := range []func() error{
		func() error { return SetupKind(mgr, o, d, Retention()) },
		func() error { return SetupKind(mgr, o, d, Caching()) },
		func() error { return SetupKind(mgr, o, d, Update()) },
		func() error { return SetupKind(mgr, o, d, RowLevelSecurity()) },
		func() error { return SetupKind(mgr, o, d, IngestionBatching()) },
		func() error { return SetupKind(mgr, o, d, StreamingIngestion()) },
		func() error { return SetupKind(mgr, o, d, Merge()) },
		func() error { return SetupKind(mgr, o, d, Sharding()) },
		func() error { return SetupKind(mgr, o, d, Partitioning()) },
		func() error { return SetupKind(mgr, o, d, IngestionTime()) },
		func() error { return SetupKind(mgr, o, d, AutoDelete()) },
		func() error { return SetupKind(mgr, o, d, RestrictedViewAccess()) },
		func() error { return SetupKind(mgr, o, d, ExtentTagsRetention()) },
		func() error { return SetupKind(mgr, o, d, Encoding()) },
		func() error { return SetupKind(mgr, o, d, ManagedIdentity()) },
		func() error { return SetupKind(mgr, o, d, RowOrder()) },
		func() error { return SetupKind(mgr, o, d, QueryAcceleration()) },
	} {
		if err := setup(); err != nil {
			return err
		}
	}
	return nil
}

// --- helpers -----------------------------------------------------------------

func ts(t *common.Timespan) (*string, error) {
	if t == nil {
		return nil, nil //nolint:nilnil // nil means "not set" here.
	}
	s, err := tsReq(*t)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func tsReq(t common.Timespan) (string, error) {
	ticks, err := timespan.Parse(string(t))
	if err != nil {
		return "", err
	}
	return ticks.DotNet(), nil
}

func tsKQL(t common.Timespan) (string, error) {
	ticks, err := timespan.Parse(string(t))
	if err != nil {
		return "", err
	}
	return ticks.KQL(), nil
}

func boolOr(b *bool, def bool) bool {
	if b == nil {
		return def
	}
	return *b
}

// enabledAlter renders ".alter <entity> policy <name> true|false".
func enabledAlter(name string) func(e adxpolicy.Entity, desired any) (cmd.Command, error) {
	return func(e adxpolicy.Entity, desired any) (cmd.Command, error) {
		ent, err := e.Render()
		if err != nil {
			return cmd.Command{}, err
		}
		d, ok := desired.(enabledPolicy)
		if !ok {
			return cmd.Command{}, fmt.Errorf("unexpected desired type %T", desired)
		}
		return cmd.New(".alter ", ent, " policy ", name, " ", cmd.Bool(d.IsEnabled)), nil
	}
}

type enabledPolicy struct {
	IsEnabled bool `json:"IsEnabled"`
}

// --- kinds -------------------------------------------------------------------

type retentionJSON struct {
	SoftDeletePeriod *string `json:"SoftDeletePeriod,omitempty"`
	Recoverability   *string `json:"Recoverability,omitempty"`
}

// Retention defines the RetentionPolicy kind.
func Retention() Def[*v1alpha1.RetentionPolicy] {
	return Def[*v1alpha1.RetentionPolicy]{
		Kind:   base.Kind[*v1alpha1.RetentionPolicy]{GVK: v1alpha1.RetentionPolicyGroupVersionKind, Object: &v1alpha1.RetentionPolicy{}, List: &v1alpha1.RetentionPolicyList{}},
		Policy: adxpolicy.Def{Name: "retention"},
		Desired: func(cr *v1alpha1.RetentionPolicy) (any, error) {
			p := cr.Spec.ForProvider
			sd, err := ts(p.SoftDeletePeriod)
			if err != nil {
				return nil, fmt.Errorf("softDeletePeriod: %w", err)
			}
			return retentionJSON{SoftDeletePeriod: sd, Recoverability: p.Recoverability}, nil
		},
	}
}

type hotWindowJSON struct {
	MinValue string `json:"MinValue"`
	MaxValue string `json:"MaxValue"`
}

type cachingJSON struct {
	DataHotSpan  string          `json:"DataHotSpan"`
	IndexHotSpan string          `json:"IndexHotSpan"`
	HotWindows   []hotWindowJSON `json:"HotWindows"`
}

// Caching defines the CachingPolicy kind. Kusto has no JSON alter form for the
// caching policy, so the "hot = <timespan>, hot_window = ..." form is used.
func Caching() Def[*v1alpha1.CachingPolicy] {
	return Def[*v1alpha1.CachingPolicy]{
		Kind: base.Kind[*v1alpha1.CachingPolicy]{GVK: v1alpha1.CachingPolicyGroupVersionKind, Object: &v1alpha1.CachingPolicy{}, List: &v1alpha1.CachingPolicyList{}},
		Policy: adxpolicy.Def{Name: "caching", Alter: func(e adxpolicy.Entity, desired any) (cmd.Command, error) {
			ent, err := e.Render()
			if err != nil {
				return cmd.Command{}, err
			}
			d, ok := desired.(cachingJSON)
			if !ok {
				return cmd.Command{}, fmt.Errorf("unexpected desired type %T", desired)
			}
			hot, err := timespan.Parse(d.DataHotSpan)
			if err != nil {
				return cmd.Command{}, err
			}
			parts := []string{".alter ", ent, " policy caching hot = ", hot.KQL()}
			for _, w := range d.HotWindows {
				parts = append(parts, ", hot_window = datetime(", w.MinValue, ") .. datetime(", w.MaxValue, ")")
			}
			return cmd.New(parts...), nil
		}},
		Desired: func(cr *v1alpha1.CachingPolicy) (any, error) {
			p := cr.Spec.ForProvider
			hot, err := tsReq(p.Hot)
			if err != nil {
				return nil, fmt.Errorf("hot: %w", err)
			}
			d := cachingJSON{DataHotSpan: hot, IndexHotSpan: hot, HotWindows: []hotWindowJSON{}}
			for _, w := range p.HotWindows {
				if !isDatetime(w.MinValue) {
					return nil, fmt.Errorf("hotWindows.minValue %q is not an ISO 8601 datetime", w.MinValue)
				}
				if !isDatetime(w.MaxValue) {
					return nil, fmt.Errorf("hotWindows.maxValue %q is not an ISO 8601 datetime", w.MaxValue)
				}
				d.HotWindows = append(d.HotWindows, hotWindowJSON{MinValue: w.MinValue, MaxValue: w.MaxValue})
			}
			return d, nil
		},
	}
}

// isDatetime is a cheap ISO 8601 plausibility check (YYYY-MM-DD prefix).
func isDatetime(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 10 || s[4] != '-' || s[7] != '-' {
		return false
	}
	for _, c := range s[:4] + s[5:7] + s[8:10] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

type updateJSON struct {
	IsEnabled                    bool    `json:"IsEnabled"`
	Source                       string  `json:"Source"`
	Query                        string  `json:"Query"`
	IsTransactional              *bool   `json:"IsTransactional,omitempty"`
	PropagateIngestionProperties *bool   `json:"PropagateIngestionProperties,omitempty"`
	ManagedIdentity              *string `json:"ManagedIdentity,omitempty"`
}

// Update defines the UpdatePolicy kind (array, Query is KQL).
func Update() Def[*v1alpha1.UpdatePolicy] {
	return Def[*v1alpha1.UpdatePolicy]{
		Kind:   base.Kind[*v1alpha1.UpdatePolicy]{GVK: v1alpha1.UpdatePolicyGroupVersionKind, Object: &v1alpha1.UpdatePolicy{}, List: &v1alpha1.UpdatePolicyList{}},
		Policy: adxpolicy.Def{Name: "update", KQLFields: []string{"Query"}},
		Desired: func(cr *v1alpha1.UpdatePolicy) (any, error) {
			out := make([]updateJSON, 0, len(cr.Spec.ForProvider.Updates))
			for i, u := range cr.Spec.ForProvider.Updates {
				if strings.TrimSpace(u.Source) == "" {
					return nil, fmt.Errorf("updates[%d].source is empty (sourceRef not resolved yet?)", i)
				}
				out = append(out, updateJSON{IsEnabled: boolOr(u.Enabled, true), Source: u.Source, Query: u.Query, IsTransactional: u.Transactional, PropagateIngestionProperties: u.PropagateIngestionProperties, ManagedIdentity: u.ManagedIdentity})
			}
			return out, nil
		},
	}
}

type rlsJSON struct {
	IsEnabled bool   `json:"IsEnabled"`
	Query     string `json:"Query"`
}

// RowLevelSecurity defines the RowLevelSecurityPolicy kind
// (".alter table T policy row_level_security enable|disable "query"").
func RowLevelSecurity() Def[*v1alpha1.RowLevelSecurityPolicy] {
	return Def[*v1alpha1.RowLevelSecurityPolicy]{
		Kind: base.Kind[*v1alpha1.RowLevelSecurityPolicy]{GVK: v1alpha1.RowLevelSecurityPolicyGroupVersionKind, Object: &v1alpha1.RowLevelSecurityPolicy{}, List: &v1alpha1.RowLevelSecurityPolicyList{}},
		Policy: adxpolicy.Def{Name: "row_level_security", KQLFields: []string{"Query"}, Alter: func(e adxpolicy.Entity, desired any) (cmd.Command, error) {
			ent, err := e.Render()
			if err != nil {
				return cmd.Command{}, err
			}
			d, ok := desired.(rlsJSON)
			if !ok {
				return cmd.Command{}, fmt.Errorf("unexpected desired type %T", desired)
			}
			mode := "enable"
			if !d.IsEnabled {
				mode = "disable"
			}
			return cmd.New(".alter ", ent, " policy row_level_security ", mode, " ", cmd.Str(d.Query)), nil
		}},
		Desired: func(cr *v1alpha1.RowLevelSecurityPolicy) (any, error) {
			return rlsJSON{IsEnabled: boolOr(cr.Spec.ForProvider.Enabled, true), Query: cr.Spec.ForProvider.Query}, nil
		},
	}
}

type ingestionBatchingJSON struct {
	MaximumBatchingTimeSpan *string `json:"MaximumBatchingTimeSpan,omitempty"`
	MaximumNumberOfItems    *int64  `json:"MaximumNumberOfItems,omitempty"`
	MaximumRawDataSizeMB    *int64  `json:"MaximumRawDataSizeMB,omitempty"`
}

// IngestionBatching defines the IngestionBatchingPolicy kind.
func IngestionBatching() Def[*v1alpha1.IngestionBatchingPolicy] {
	return Def[*v1alpha1.IngestionBatchingPolicy]{
		Kind:   base.Kind[*v1alpha1.IngestionBatchingPolicy]{GVK: v1alpha1.IngestionBatchingPolicyGroupVersionKind, Object: &v1alpha1.IngestionBatchingPolicy{}, List: &v1alpha1.IngestionBatchingPolicyList{}},
		Policy: adxpolicy.Def{Name: "ingestionbatching"},
		Desired: func(cr *v1alpha1.IngestionBatchingPolicy) (any, error) {
			p := cr.Spec.ForProvider
			t, err := ts(p.MaximumBatchingTimeSpan)
			if err != nil {
				return nil, fmt.Errorf("maximumBatchingTimeSpan: %w", err)
			}
			return ingestionBatchingJSON{MaximumBatchingTimeSpan: t, MaximumNumberOfItems: p.MaximumNumberOfItems, MaximumRawDataSizeMB: p.MaximumRawDataSizeMB}, nil
		},
	}
}

type streamingIngestionJSON struct {
	IsEnabled         bool     `json:"IsEnabled"`
	HintAllocatedRate *float64 `json:"HintAllocatedRate,omitempty"`
}

// StreamingIngestion defines the StreamingIngestionPolicy kind.
func StreamingIngestion() Def[*v1alpha1.StreamingIngestionPolicy] {
	return Def[*v1alpha1.StreamingIngestionPolicy]{
		Kind:   base.Kind[*v1alpha1.StreamingIngestionPolicy]{GVK: v1alpha1.StreamingIngestionPolicyGroupVersionKind, Object: &v1alpha1.StreamingIngestionPolicy{}, List: &v1alpha1.StreamingIngestionPolicyList{}},
		Policy: adxpolicy.Def{Name: "streamingingestion"},
		Desired: func(cr *v1alpha1.StreamingIngestionPolicy) (any, error) {
			p := cr.Spec.ForProvider
			d := streamingIngestionJSON{IsEnabled: p.Enabled}
			if p.HintAllocatedRate != nil {
				f, err := strconv.ParseFloat(*p.HintAllocatedRate, 64)
				if err != nil {
					return nil, fmt.Errorf("hintAllocatedRate: %w", err)
				}
				d.HintAllocatedRate = &f
			}
			return d, nil
		},
	}
}

type mergeLookbackJSON struct {
	Kind         string  `json:"Kind"`
	CustomPeriod *string `json:"CustomPeriod,omitempty"`
}

type mergeJSON struct {
	RowCountUpperBoundForMerge       *int64             `json:"RowCountUpperBoundForMerge,omitempty"`
	OriginalSizeMBUpperBoundForMerge *int64             `json:"OriginalSizeMBUpperBoundForMerge,omitempty"`
	MaxExtentsToMerge                *int64             `json:"MaxExtentsToMerge,omitempty"`
	LoopPeriod                       *string            `json:"LoopPeriod,omitempty"`
	MaxRangeInHours                  *int64             `json:"MaxRangeInHours,omitempty"`
	AllowRebuild                     *bool              `json:"AllowRebuild,omitempty"`
	AllowMerge                       *bool              `json:"AllowMerge,omitempty"`
	Lookback                         *mergeLookbackJSON `json:"Lookback,omitempty"`
}

// Merge defines the MergePolicy kind.
func Merge() Def[*v1alpha1.MergePolicy] {
	return Def[*v1alpha1.MergePolicy]{
		Kind:   base.Kind[*v1alpha1.MergePolicy]{GVK: v1alpha1.MergePolicyGroupVersionKind, Object: &v1alpha1.MergePolicy{}, List: &v1alpha1.MergePolicyList{}},
		Policy: adxpolicy.Def{Name: "merge"},
		Desired: func(cr *v1alpha1.MergePolicy) (any, error) {
			p := cr.Spec.ForProvider
			lp, err := ts(p.LoopPeriod)
			if err != nil {
				return nil, fmt.Errorf("loopPeriod: %w", err)
			}
			d := mergeJSON{RowCountUpperBoundForMerge: p.RowCountUpperBoundForMerge, OriginalSizeMBUpperBoundForMerge: p.OriginalSizeMBUpperBoundForMerge, MaxExtentsToMerge: p.MaxExtentsToMerge, LoopPeriod: lp, MaxRangeInHours: p.MaxRangeInHours, AllowRebuild: p.AllowRebuild, AllowMerge: p.AllowMerge}
			if p.Lookback != nil {
				cp, err := ts(p.Lookback.CustomPeriod)
				if err != nil {
					return nil, fmt.Errorf("lookback.customPeriod: %w", err)
				}
				d.Lookback = &mergeLookbackJSON{Kind: p.Lookback.Kind, CustomPeriod: cp}
			}
			return d, nil
		},
	}
}

type shardingJSON struct {
	MaxRowCount                    *int64 `json:"MaxRowCount,omitempty"`
	MaxExtentSizeInMb              *int64 `json:"MaxExtentSizeInMb,omitempty"`
	MaxOriginalSizeInMb            *int64 `json:"MaxOriginalSizeInMb,omitempty"`
	ShardEngineMaxRowCount         *int64 `json:"ShardEngineMaxRowCount,omitempty"`
	ShardEngineMaxExtentSizeInMb   *int64 `json:"ShardEngineMaxExtentSizeInMb,omitempty"`
	ShardEngineMaxOriginalSizeInMb *int64 `json:"ShardEngineMaxOriginalSizeInMb,omitempty"`
}

// Sharding defines the ShardingPolicy kind.
func Sharding() Def[*v1alpha1.ShardingPolicy] {
	return Def[*v1alpha1.ShardingPolicy]{
		Kind:   base.Kind[*v1alpha1.ShardingPolicy]{GVK: v1alpha1.ShardingPolicyGroupVersionKind, Object: &v1alpha1.ShardingPolicy{}, List: &v1alpha1.ShardingPolicyList{}},
		Policy: adxpolicy.Def{Name: "sharding"},
		Desired: func(cr *v1alpha1.ShardingPolicy) (any, error) {
			p := cr.Spec.ForProvider
			return shardingJSON{MaxRowCount: p.MaxRowCount, MaxExtentSizeInMb: p.MaxExtentSizeInMb, MaxOriginalSizeInMb: p.MaxOriginalSizeInMb, ShardEngineMaxRowCount: p.ShardEngineMaxRowCount, ShardEngineMaxExtentSizeInMb: p.ShardEngineMaxExtentSizeInMb, ShardEngineMaxOriginalSizeInMb: p.ShardEngineMaxOriginalSizeInMb}, nil
		},
	}
}

type partitionKeyPropertiesJSON struct {
	Function                *string `json:"Function,omitempty"`
	MaxPartitionCount       *int64  `json:"MaxPartitionCount,omitempty"`
	Seed                    *int64  `json:"Seed,omitempty"`
	PartitionAssignmentMode *string `json:"PartitionAssignmentMode,omitempty"`
	Reference               *string `json:"Reference,omitempty"`
	RangeSize               *string `json:"RangeSize,omitempty"`
	OverrideCreationTime    *bool   `json:"OverrideCreationTime,omitempty"`
}

type partitionKeyJSON struct {
	ColumnName string                     `json:"ColumnName"`
	Kind       string                     `json:"Kind"`
	Properties partitionKeyPropertiesJSON `json:"Properties"`
}

type partitioningJSON struct {
	PartitionKeys     []partitionKeyJSON `json:"PartitionKeys"`
	EffectiveDateTime *string            `json:"EffectiveDateTime,omitempty"`
}

// Partitioning defines the PartitioningPolicy kind.
func Partitioning() Def[*v1alpha1.PartitioningPolicy] {
	return Def[*v1alpha1.PartitioningPolicy]{
		Kind:   base.Kind[*v1alpha1.PartitioningPolicy]{GVK: v1alpha1.PartitioningPolicyGroupVersionKind, Object: &v1alpha1.PartitioningPolicy{}, List: &v1alpha1.PartitioningPolicyList{}},
		Policy: adxpolicy.Def{Name: "partitioning"},
		Desired: func(cr *v1alpha1.PartitioningPolicy) (any, error) {
			p := cr.Spec.ForProvider
			d := partitioningJSON{EffectiveDateTime: p.EffectiveDateTime, PartitionKeys: make([]partitionKeyJSON, 0, len(p.PartitionKeys))}
			for i, k := range p.PartitionKeys {
				rs, err := ts(k.Properties.RangeSize)
				if err != nil {
					return nil, fmt.Errorf("partitionKeys[%d].properties.rangeSize: %w", i, err)
				}
				d.PartitionKeys = append(d.PartitionKeys, partitionKeyJSON{ColumnName: k.ColumnName, Kind: k.Kind, Properties: partitionKeyPropertiesJSON{
					Function: k.Properties.Function, MaxPartitionCount: k.Properties.MaxPartitionCount, Seed: k.Properties.Seed, PartitionAssignmentMode: k.Properties.PartitionAssignmentMode,
					Reference: k.Properties.Reference, RangeSize: rs, OverrideCreationTime: k.Properties.OverrideCreationTime,
				}})
			}
			return d, nil
		},
	}
}

// IngestionTime defines the IngestionTimePolicy kind (".alter table T policy ingestiontime true").
func IngestionTime() Def[*v1alpha1.IngestionTimePolicy] {
	return Def[*v1alpha1.IngestionTimePolicy]{
		Kind:   base.Kind[*v1alpha1.IngestionTimePolicy]{GVK: v1alpha1.IngestionTimePolicyGroupVersionKind, Object: &v1alpha1.IngestionTimePolicy{}, List: &v1alpha1.IngestionTimePolicyList{}},
		Policy: adxpolicy.Def{Name: "ingestiontime", Alter: enabledAlter("ingestiontime")},
		Desired: func(cr *v1alpha1.IngestionTimePolicy) (any, error) {
			return enabledPolicy{IsEnabled: cr.Spec.ForProvider.Enabled}, nil
		},
	}
}

type autoDeleteJSON struct {
	ExpiryDate       string `json:"ExpiryDate"`
	DeleteIfNotEmpty *bool  `json:"DeleteIfNotEmpty,omitempty"`
}

// AutoDelete defines the AutoDeletePolicy kind.
func AutoDelete() Def[*v1alpha1.AutoDeletePolicy] {
	return Def[*v1alpha1.AutoDeletePolicy]{
		Kind:   base.Kind[*v1alpha1.AutoDeletePolicy]{GVK: v1alpha1.AutoDeletePolicyGroupVersionKind, Object: &v1alpha1.AutoDeletePolicy{}, List: &v1alpha1.AutoDeletePolicyList{}},
		Policy: adxpolicy.Def{Name: "auto_delete"},
		Desired: func(cr *v1alpha1.AutoDeletePolicy) (any, error) {
			p := cr.Spec.ForProvider
			if !isDatetime(p.ExpiryDate) {
				return nil, fmt.Errorf("expiryDate %q is not an ISO 8601 date", p.ExpiryDate)
			}
			return autoDeleteJSON{ExpiryDate: p.ExpiryDate, DeleteIfNotEmpty: p.DeleteIfNotEmpty}, nil
		},
	}
}

// RestrictedViewAccess defines the RestrictedViewAccessPolicy kind.
func RestrictedViewAccess() Def[*v1alpha1.RestrictedViewAccessPolicy] {
	return Def[*v1alpha1.RestrictedViewAccessPolicy]{
		Kind:   base.Kind[*v1alpha1.RestrictedViewAccessPolicy]{GVK: v1alpha1.RestrictedViewAccessPolicyGroupVersionKind, Object: &v1alpha1.RestrictedViewAccessPolicy{}, List: &v1alpha1.RestrictedViewAccessPolicyList{}},
		Policy: adxpolicy.Def{Name: "restricted_view_access", Alter: enabledAlter("restricted_view_access")},
		Desired: func(cr *v1alpha1.RestrictedViewAccessPolicy) (any, error) {
			return enabledPolicy{IsEnabled: cr.Spec.ForProvider.Enabled}, nil
		},
	}
}

type extentTagsRetentionJSON struct {
	TagPrefix       string `json:"TagPrefix"`
	RetentionPeriod string `json:"RetentionPeriod"`
}

// ExtentTagsRetention defines the ExtentTagsRetentionPolicy kind (array).
func ExtentTagsRetention() Def[*v1alpha1.ExtentTagsRetentionPolicy] {
	return Def[*v1alpha1.ExtentTagsRetentionPolicy]{
		Kind:   base.Kind[*v1alpha1.ExtentTagsRetentionPolicy]{GVK: v1alpha1.ExtentTagsRetentionPolicyGroupVersionKind, Object: &v1alpha1.ExtentTagsRetentionPolicy{}, List: &v1alpha1.ExtentTagsRetentionPolicyList{}},
		Policy: adxpolicy.Def{Name: "extent_tags_retention"},
		Desired: func(cr *v1alpha1.ExtentTagsRetentionPolicy) (any, error) {
			out := make([]extentTagsRetentionJSON, 0, len(cr.Spec.ForProvider.Rules))
			for i, r := range cr.Spec.ForProvider.Rules {
				rp, err := tsReq(r.RetentionPeriod)
				if err != nil {
					return nil, fmt.Errorf("rules[%d].retentionPeriod: %w", i, err)
				}
				out = append(out, extentTagsRetentionJSON{TagPrefix: r.TagPrefix, RetentionPeriod: rp})
			}
			return out, nil
		},
	}
}

type encodingJSON struct {
	Type string `json:"Type"`
}

// Encoding defines the EncodingPolicy kind (".alter <entity> policy encoding
// type = 'X'"). The .show JSON shape is not verified, so comparison relies on
// the hash annotations.
func Encoding() Def[*v1alpha1.EncodingPolicy] {
	return Def[*v1alpha1.EncodingPolicy]{
		Kind: base.Kind[*v1alpha1.EncodingPolicy]{GVK: v1alpha1.EncodingPolicyGroupVersionKind, Object: &v1alpha1.EncodingPolicy{}, List: &v1alpha1.EncodingPolicyList{}},
		Policy: adxpolicy.Def{Name: "encoding", HashOnly: true, NoBatch: true, Alter: func(e adxpolicy.Entity, desired any) (cmd.Command, error) {
			ent, err := e.Render()
			if err != nil {
				return cmd.Command{}, err
			}
			d, ok := desired.(encodingJSON)
			if !ok {
				return cmd.Command{}, fmt.Errorf("unexpected desired type %T", desired)
			}
			return cmd.New(".alter ", ent, " policy encoding type = ", cmd.Str(d.Type)), nil
		}},
		Desired: func(cr *v1alpha1.EncodingPolicy) (any, error) {
			return encodingJSON{Type: cr.Spec.ForProvider.Type}, nil
		},
		Column: func(cr *v1alpha1.EncodingPolicy) string {
			if cr.Spec.ForProvider.Column == nil {
				return ""
			}
			return *cr.Spec.ForProvider.Column
		},
	}
}

type managedIdentityJSON struct {
	ObjectID      string `json:"ObjectId"`
	AllowedUsages string `json:"AllowedUsages"`
}

// ManagedIdentity defines the (database level) ManagedIdentityPolicy kind.
func ManagedIdentity() Def[*v1alpha1.ManagedIdentityPolicy] {
	return Def[*v1alpha1.ManagedIdentityPolicy]{
		Kind:   base.Kind[*v1alpha1.ManagedIdentityPolicy]{GVK: v1alpha1.ManagedIdentityPolicyGroupVersionKind, Object: &v1alpha1.ManagedIdentityPolicy{}, List: &v1alpha1.ManagedIdentityPolicyList{}},
		Policy: adxpolicy.Def{Name: "managed_identity", SetFields: []string{"AllowedUsages"}},
		Desired: func(cr *v1alpha1.ManagedIdentityPolicy) (any, error) {
			out := make([]managedIdentityJSON, 0, len(cr.Spec.ForProvider.Identities))
			for _, id := range cr.Spec.ForProvider.Identities {
				out = append(out, managedIdentityJSON{ObjectID: id.ObjectID, AllowedUsages: strings.Join(id.AllowedUsages, ", ")})
			}
			return out, nil
		},
	}
}

type rowOrderJSON struct {
	Column    string `json:"Column"`
	Direction string `json:"Direction"`
}

// RowOrder defines the RowOrderPolicy kind (".alter table T policy roworder
// (a asc, b desc)"). Tier 3, unverified .show shape: hash comparison.
func RowOrder() Def[*v1alpha1.RowOrderPolicy] {
	return Def[*v1alpha1.RowOrderPolicy]{
		Kind: base.Kind[*v1alpha1.RowOrderPolicy]{GVK: v1alpha1.RowOrderPolicyGroupVersionKind, Object: &v1alpha1.RowOrderPolicy{}, List: &v1alpha1.RowOrderPolicyList{}},
		Policy: adxpolicy.Def{Name: "roworder", HashOnly: true, NoBatch: true, Alter: func(e adxpolicy.Entity, desired any) (cmd.Command, error) {
			ent, err := e.Render()
			if err != nil {
				return cmd.Command{}, err
			}
			d, ok := desired.([]rowOrderJSON)
			if !ok {
				return cmd.Command{}, fmt.Errorf("unexpected desired type %T", desired)
			}
			items := make([]string, 0, len(d))
			for _, c := range d {
				items = append(items, cmd.Ident(c.Column)+" "+c.Direction)
			}
			return cmd.New(".alter ", ent, " policy roworder ", cmd.List(items)), nil
		}},
		Desired: func(cr *v1alpha1.RowOrderPolicy) (any, error) {
			out := make([]rowOrderJSON, 0, len(cr.Spec.ForProvider.Columns))
			for _, c := range cr.Spec.ForProvider.Columns {
				out = append(out, rowOrderJSON{Column: c.Name, Direction: strings.ToLower(c.Direction)})
			}
			return out, nil
		},
	}
}

type queryAccelerationJSON struct {
	IsEnabled bool    `json:"IsEnabled"`
	Hot       *string `json:"Hot,omitempty"`
	MaxAge    *string `json:"MaxAge,omitempty"`
}

// QueryAcceleration defines the QueryAccelerationPolicy kind (external delta tables).
func QueryAcceleration() Def[*v1alpha1.QueryAccelerationPolicy] {
	return Def[*v1alpha1.QueryAccelerationPolicy]{
		Kind:   base.Kind[*v1alpha1.QueryAccelerationPolicy]{GVK: v1alpha1.QueryAccelerationPolicyGroupVersionKind, Object: &v1alpha1.QueryAccelerationPolicy{}, List: &v1alpha1.QueryAccelerationPolicyList{}},
		Policy: adxpolicy.Def{Name: "query_acceleration", NoBatch: true},
		Desired: func(cr *v1alpha1.QueryAccelerationPolicy) (any, error) {
			p := cr.Spec.ForProvider
			hot, err := ts(p.Hot)
			if err != nil {
				return nil, fmt.Errorf("hot: %w", err)
			}
			maxAge, err := ts(p.MaxAge)
			if err != nil {
				return nil, fmt.Errorf("maxAge: %w", err)
			}
			return queryAccelerationJSON{IsEnabled: p.Enabled, Hot: hot, MaxAge: maxAge}, nil
		},
	}
}

var _ = tsKQL
