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

package v1alpha1

import resource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"

// Interface checks live in a test file so that angryjet (which generates the
// methods asserted here) can type-check the package before they exist.
var (
	_ resource.ModernManaged = &RetentionPolicy{}
	_ resource.ManagedList   = &RetentionPolicyList{}
	_ resource.ModernManaged = &CachingPolicy{}
	_ resource.ManagedList   = &CachingPolicyList{}
	_ resource.ModernManaged = &UpdatePolicy{}
	_ resource.ManagedList   = &UpdatePolicyList{}
	_ resource.ModernManaged = &RowLevelSecurityPolicy{}
	_ resource.ManagedList   = &RowLevelSecurityPolicyList{}
	_ resource.ModernManaged = &IngestionBatchingPolicy{}
	_ resource.ManagedList   = &IngestionBatchingPolicyList{}
	_ resource.ModernManaged = &StreamingIngestionPolicy{}
	_ resource.ManagedList   = &StreamingIngestionPolicyList{}
	_ resource.ModernManaged = &MergePolicy{}
	_ resource.ManagedList   = &MergePolicyList{}
	_ resource.ModernManaged = &ShardingPolicy{}
	_ resource.ManagedList   = &ShardingPolicyList{}
	_ resource.ModernManaged = &PartitioningPolicy{}
	_ resource.ManagedList   = &PartitioningPolicyList{}
	_ resource.ModernManaged = &IngestionTimePolicy{}
	_ resource.ManagedList   = &IngestionTimePolicyList{}
	_ resource.ModernManaged = &AutoDeletePolicy{}
	_ resource.ManagedList   = &AutoDeletePolicyList{}
	_ resource.ModernManaged = &RestrictedViewAccessPolicy{}
	_ resource.ManagedList   = &RestrictedViewAccessPolicyList{}
	_ resource.ModernManaged = &ExtentTagsRetentionPolicy{}
	_ resource.ManagedList   = &ExtentTagsRetentionPolicyList{}
	_ resource.ModernManaged = &EncodingPolicy{}
	_ resource.ManagedList   = &EncodingPolicyList{}
	_ resource.ModernManaged = &ManagedIdentityPolicy{}
	_ resource.ManagedList   = &ManagedIdentityPolicyList{}
	_ resource.ModernManaged = &RowOrderPolicy{}
	_ resource.ManagedList   = &RowOrderPolicyList{}
	_ resource.ModernManaged = &QueryAccelerationPolicy{}
	_ resource.ManagedList   = &QueryAccelerationPolicyList{}
)
