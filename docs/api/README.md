# API Reference

Packages:

- [adx.functional.team/v1alpha1](#adxfunctionalteamv1alpha1)
- [cluster.adx.functional.team/v1alpha1](#clusteradxfunctionalteamv1alpha1)
- [policy.adx.functional.team/v1alpha1](#policyadxfunctionalteamv1alpha1)
- [security.adx.functional.team/v1alpha1](#securityadxfunctionalteamv1alpha1)

# adx.functional.team/v1alpha1

Resource Types:

- [ClusterProviderConfig](#clusterproviderconfig)

- [ContinuousExport](#continuousexport)

- [EntityGroup](#entitygroup)

- [ExternalTable](#externaltable)

- [Function](#function)

- [IngestionMapping](#ingestionmapping)

- [MaterializedView](#materializedview)

- [ProviderConfig](#providerconfig)

- [ProviderConfigUsage](#providerconfigusage)

- [Table](#table)




## ClusterProviderConfig
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






A ClusterProviderConfig is the cluster-scoped variant of ProviderConfig. It
can be referenced by managed resources in any namespace.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ClusterProviderConfig</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#clusterproviderconfigspec">spec</a></b></td>
        <td>object</td>
        <td>
          A ProviderConfigSpec defines the desired state of a ProviderConfig.<br/>
          <br/>
            <i>Validations</i>:<li>self.credentials.source != 'None' || self.clusterUri.startsWith('http://'): credentials.source None is only allowed for http:// endpoints (Kusto emulator)</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#clusterproviderconfigstatus">status</a></b></td>
        <td>object</td>
        <td>
          A ProviderConfigStatus defines the status of a Provider.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterProviderConfig.spec
<sup><sup>[↩ Parent](#clusterproviderconfig)</sup></sup>



A ProviderConfigSpec defines the desired state of a ProviderConfig.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>clusterUri</b></td>
        <td>string</td>
        <td>
          ClusterURI is the Kusto query endpoint of the cluster, e.g.
https://mycluster.westeurope.kusto.windows.net. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: clusterUri is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#clusterproviderconfigspeccredentials">credentials</a></b></td>
        <td>object</td>
        <td>
          Credentials required to authenticate to this provider.<br/>
          <br/>
            <i>Validations</i>:<li>self.source != 'Secret' || has(self.secretRef): secretRef is required when source is Secret</li><li>self.source != 'ManagedIdentity' || has(self.managedIdentity): managedIdentity is required when source is ManagedIdentity</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>azureEnvironment</b></td>
        <td>enum</td>
        <td>
          AzureEnvironment selects the sovereign cloud. Only AzurePublicCloud is
tested; the others are wired through to the SDK but unverified.<br/>
          <br/>
            <i>Enum</i>: AzurePublicCloud, AzureChinaCloud, AzureUSGovernment<br/>
            <i>Default</i>: AzurePublicCloud<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterProviderConfig.spec.credentials
<sup><sup>[↩ Parent](#clusterproviderconfigspec)</sup></sup>



Credentials required to authenticate to this provider.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>source</b></td>
        <td>enum</td>
        <td>
          Source of the provider credentials.<br/>
          <br/>
            <i>Enum</i>: None, Secret, WorkloadIdentity, ManagedIdentity<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#clusterproviderconfigspeccredentialsmanagedidentity">managedIdentity</a></b></td>
        <td>object</td>
        <td>
          ManagedIdentity settings, used when source is ManagedIdentity.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#clusterproviderconfigspeccredentialssecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          SecretRef points to a Secret key holding a JSON document with the keys
clientId, clientSecret and tenantId. The format is compatible with the
credentials Secret of provider-upbound-azure (additional keys are ignored).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#clusterproviderconfigspeccredentialsworkloadidentity">workloadIdentity</a></b></td>
        <td>object</td>
        <td>
          WorkloadIdentity settings, used when source is WorkloadIdentity.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterProviderConfig.spec.credentials.managedIdentity
<sup><sup>[↩ Parent](#clusterproviderconfigspeccredentials)</sup></sup>



ManagedIdentity settings, used when source is ManagedIdentity.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type of the managed identity.<br/>
          <br/>
            <i>Enum</i>: SystemAssigned, UserAssigned<br/>
            <i>Default</i>: SystemAssigned<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>clientId</b></td>
        <td>string</td>
        <td>
          ClientID of a user-assigned identity. Mutually exclusive with resourceId.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resourceId</b></td>
        <td>string</td>
        <td>
          ResourceID of a user-assigned identity. Mutually exclusive with clientId.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterProviderConfig.spec.credentials.secretRef
<sup><sup>[↩ Parent](#clusterproviderconfigspeccredentials)</sup></sup>



SecretRef points to a Secret key holding a JSON document with the keys
clientId, clientSecret and tenantId. The format is compatible with the
credentials Secret of provider-upbound-azure (additional keys are ignored).

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key to select.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ClusterProviderConfig.spec.credentials.workloadIdentity
<sup><sup>[↩ Parent](#clusterproviderconfigspeccredentials)</sup></sup>



WorkloadIdentity settings, used when source is WorkloadIdentity.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>clientId</b></td>
        <td>string</td>
        <td>
          ClientID of the federated application.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenantId</b></td>
        <td>string</td>
        <td>
          TenantID of the federated application.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tokenFile</b></td>
        <td>string</td>
        <td>
          TokenFile is the path of the projected service account token.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterProviderConfig.status
<sup><sup>[↩ Parent](#clusterproviderconfig)</sup></sup>



A ProviderConfigStatus defines the status of a Provider.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#clusterproviderconfigstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>users</b></td>
        <td>integer</td>
        <td>
          Users of this provider configuration.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterProviderConfig.status.conditions[index]
<sup><sup>[↩ Parent](#clusterproviderconfigstatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ContinuousExport
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






A ContinuousExport periodically exports the results of a query to an
external table.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ContinuousExport</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#continuousexportspec">spec</a></b></td>
        <td>object</td>
        <td>
          A ContinuousExportSpec defines the desired state of a ContinuousExport.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#continuousexportstatus">status</a></b></td>
        <td>object</td>
        <td>
          A ContinuousExportStatus represents the observed state of a ContinuousExport.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec
<sup><sup>[↩ Parent](#continuousexport)</sup></sup>



A ContinuousExportSpec defines the desired state of a ContinuousExport.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#continuousexportspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          ContinuousExportParameters are the configurable fields of a ContinuousExport.<br/>
          <br/>
            <i>Validations</i>:<li>has(self.externalTable) || has(self.externalTableRef) || has(self.externalTableSelector): externalTable, externalTableRef or externalTableSelector is required</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.forProvider
<sup><sup>[↩ Parent](#continuousexportspec)</sup></sup>



ContinuousExportParameters are the configurable fields of a ContinuousExport.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the continuous export. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>intervalBetweenRuns</b></td>
        <td>string</td>
        <td>
          IntervalBetweenRuns between export runs, e.g. "1h".<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>query</b></td>
        <td>string</td>
        <td>
          Query is the KQL query whose results are exported.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>distributed</b></td>
        <td>boolean</td>
        <td>
          Distributed exports from multiple nodes concurrently.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>distribution</b></td>
        <td>enum</td>
        <td>
          Distribution hint: single, per_node or per_shard.<br/>
          <br/>
            <i>Enum</i>: single, per_node, per_shard<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>distributionKind</b></td>
        <td>string</td>
        <td>
          DistributionKind hint, e.g. default or uniform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled toggles the export. Defaults to true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>externalTable</b></td>
        <td>string</td>
        <td>
          ExternalTable is the target external table.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecforproviderexternaltableref">externalTableRef</a></b></td>
        <td>object</td>
        <td>
          ExternalTableRef references an ExternalTable managed resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecforproviderexternaltableselector">externalTableSelector</a></b></td>
        <td>object</td>
        <td>
          ExternalTableSelector selects an ExternalTable managed resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>forcedLatency</b></td>
        <td>string</td>
        <td>
          ForcedLatency delays the export to include late arriving data.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>managedIdentity</b></td>
        <td>string</td>
        <td>
          ManagedIdentity ("system" or an object id) used to write to storage.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the continuous export in Kusto. Written to the
crossplane.io/external-name annotation on the first reconcile if that
annotation is empty. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: name is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>overTables</b></td>
        <td>[]string</td>
        <td>
          OverTables are the fact tables the export cursor is scoped to. When
omitted Kusto derives them from the query.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecforproviderovertablesrefsindex">overTablesRefs</a></b></td>
        <td>[]object</td>
        <td>
          OverTablesRefs reference Table managed resources.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecforproviderovertablesselector">overTablesSelector</a></b></td>
        <td>object</td>
        <td>
          OverTablesSelector selects Table managed resources.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parquetRowGroupSize</b></td>
        <td>integer</td>
        <td>
          ParquetRowGroupSize for parquet exports.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sizeLimit</b></td>
        <td>integer</td>
        <td>
          SizeLimit of a single exported file in bytes.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.forProvider.externalTableRef
<sup><sup>[↩ Parent](#continuousexportspecforprovider)</sup></sup>



ExternalTableRef references an ExternalTable managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecforproviderexternaltablerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.forProvider.externalTableRef.policy
<sup><sup>[↩ Parent](#continuousexportspecforproviderexternaltableref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.forProvider.externalTableSelector
<sup><sup>[↩ Parent](#continuousexportspecforprovider)</sup></sup>



ExternalTableSelector selects an ExternalTable managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecforproviderexternaltableselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.forProvider.externalTableSelector.policy
<sup><sup>[↩ Parent](#continuousexportspecforproviderexternaltableselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.forProvider.overTablesRefs[index]
<sup><sup>[↩ Parent](#continuousexportspecforprovider)</sup></sup>



A NamespacedReference to a named object.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecforproviderovertablesrefsindexpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.forProvider.overTablesRefs[index].policy
<sup><sup>[↩ Parent](#continuousexportspecforproviderovertablesrefsindex)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.forProvider.overTablesSelector
<sup><sup>[↩ Parent](#continuousexportspecforprovider)</sup></sup>



OverTablesSelector selects Table managed resources.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportspecforproviderovertablesselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.forProvider.overTablesSelector.policy
<sup><sup>[↩ Parent](#continuousexportspecforproviderovertablesselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.providerConfigRef
<sup><sup>[↩ Parent](#continuousexportspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ContinuousExport.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#continuousexportspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ContinuousExport.status
<sup><sup>[↩ Parent](#continuousexport)</sup></sup>



A ContinuousExportStatus represents the observed state of a ContinuousExport.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#continuousexportstatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          ContinuousExportObservation are the observable fields of a ContinuousExport.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#continuousexportstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.status.atProvider
<sup><sup>[↩ Parent](#continuousexportstatus)</sup></sup>



ContinuousExportObservation are the observable fields of a ContinuousExport.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>exportedTo</b></td>
        <td>string</td>
        <td>
          ExportedTo is the ingestion time up to which data was exported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>externalTable</b></td>
        <td>string</td>
        <td>
          ExternalTable as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isDisabled</b></td>
        <td>boolean</td>
        <td>
          IsDisabled as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isRunning</b></td>
        <td>boolean</td>
        <td>
          IsRunning reports whether an export run is in progress.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastRunResult</b></td>
        <td>string</td>
        <td>
          LastRunResult of the export.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastRunTime</b></td>
        <td>string</td>
        <td>
          LastRunTime of the export.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>query</b></td>
        <td>string</td>
        <td>
          Query as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ContinuousExport.status.conditions[index]
<sup><sup>[↩ Parent](#continuousexportstatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## EntityGroup
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






An EntityGroup is a named set of Kusto entities usable with the macro-expand
operator (Tier 3).

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>EntityGroup</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#entitygroupspec">spec</a></b></td>
        <td>object</td>
        <td>
          An EntityGroupSpec defines the desired state of an EntityGroup.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#entitygroupstatus">status</a></b></td>
        <td>object</td>
        <td>
          An EntityGroupStatus represents the observed state of an EntityGroup.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EntityGroup.spec
<sup><sup>[↩ Parent](#entitygroup)</sup></sup>



An EntityGroupSpec defines the desired state of an EntityGroup.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#entitygroupspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          EntityGroupParameters are the configurable fields of an EntityGroup.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#entitygroupspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#entitygroupspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EntityGroup.spec.forProvider
<sup><sup>[↩ Parent](#entitygroupspec)</sup></sup>



EntityGroupParameters are the configurable fields of an EntityGroup.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity group. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>entities</b></td>
        <td>[]string</td>
        <td>
          Entities are raw KQL entity expressions such as
cluster('c').database('d') or database('d').table('t'). They are code
and are sent verbatim; order does not matter.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the entity group in Kusto. Written to the
crossplane.io/external-name annotation on the first reconcile if that
annotation is empty. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: name is immutable</li>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EntityGroup.spec.providerConfigRef
<sup><sup>[↩ Parent](#entitygroupspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### EntityGroup.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#entitygroupspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### EntityGroup.status
<sup><sup>[↩ Parent](#entitygroup)</sup></sup>



An EntityGroupStatus represents the observed state of an EntityGroup.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#entitygroupstatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          EntityGroupObservation are the observable fields of an EntityGroup.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#entitygroupstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EntityGroup.status.atProvider
<sup><sup>[↩ Parent](#entitygroupstatus)</sup></sup>



EntityGroupObservation are the observable fields of an EntityGroup.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entities</b></td>
        <td>[]string</td>
        <td>
          Entities as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EntityGroup.status.conditions[index]
<sup><sup>[↩ Parent](#entitygroupstatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ExternalTable
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






An ExternalTable is a Kusto external table over blob/ADLS storage or a
Delta Lake table. Secrets in connection strings are write-only; prefer
managed identity (";managed_identity=system"). Creating an external table
with a managed identity requires the AllDatabasesAdmin role.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ExternalTable</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#externaltablespec">spec</a></b></td>
        <td>object</td>
        <td>
          An ExternalTableSpec defines the desired state of an ExternalTable.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#externaltablestatus">status</a></b></td>
        <td>object</td>
        <td>
          An ExternalTableStatus represents the observed state of an ExternalTable.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.spec
<sup><sup>[↩ Parent](#externaltable)</sup></sup>



An ExternalTableSpec defines the desired state of an ExternalTable.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#externaltablespecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          ExternalTableParameters are the configurable fields of an ExternalTable.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind != 'Storage' || (has(self.columns) && size(self.columns) > 0): columns are required for kind Storage</li><li>self.kind != 'Storage' || has(self.dataFormat): dataFormat is required for kind Storage</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#externaltablespecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#externaltablespecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.spec.forProvider
<sup><sup>[↩ Parent](#externaltablespec)</sup></sup>



ExternalTableParameters are the configurable fields of an ExternalTable.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#externaltablespecforproviderconnectionstringsindex">connectionStrings</a></b></td>
        <td>[]object</td>
        <td>
          ConnectionStrings of the storage containers, at least one.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the external table. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#externaltablespecforprovidercolumnsindex">columns</a></b></td>
        <td>[]object</td>
        <td>
          Columns of the external table. Required for Storage; optional for Delta
(schema is inferred from the delta log when omitted).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>dataFormat</b></td>
        <td>enum</td>
        <td>
          DataFormat of the files (Storage kind).<br/>
          <br/>
            <i>Enum</i>: csv, tsv, tsve, psv, scsv, sohsv, json, multijson, avro, apacheavro, parquet, orc, w3clogfile, txt, raw, sstream<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the external table. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: kind is immutable</li>
            <i>Enum</i>: Storage, Delta<br/>
            <i>Default</i>: Storage<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the external table in Kusto. Written to the
crossplane.io/external-name annotation on the first reconcile if that
annotation is empty. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: name is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>partitionBy</b></td>
        <td>string</td>
        <td>
          PartitionBy is the raw Kusto partition expression, e.g.
"Date:datetime = bin(Timestamp, 1d)". Compared through the hash
mechanism because Kusto reports partitions as JSON.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pathFormat</b></td>
        <td>string</td>
        <td>
          PathFormat is the raw Kusto path format expression, e.g.
"datetime_pattern(\"yyyy/MM/dd\", Date)".<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#externaltablespecforproviderproperties">properties</a></b></td>
        <td>object</td>
        <td>
          Properties are the optional "with (...)" settings.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.spec.forProvider.connectionStrings[index]
<sup><sup>[↩ Parent](#externaltablespecforprovider)</sup></sup>



ExternalTableConnectionString is one storage connection string. Exactly one
of value and secretKeyRef must be set. Connection strings are always sent
obfuscated (h@'...') so Kusto hides everything after the first ';' in .show
output and command logs; only the URI part is compared and reported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#externaltablespecforproviderconnectionstringsindexsecretkeyref">secretKeyRef</a></b></td>
        <td>object</td>
        <td>
          SecretKeyRef reads the connection string (with SAS token or account
key) from a Secret. Write-only: a changed secret is detected through a
hash, the value never appears in status, events or logs.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Value is the connection string, e.g.
"https://acct.blob.core.windows.net/exports;managed_identity=system".
Use it for identity based access without secrets.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.spec.forProvider.connectionStrings[index].secretKeyRef
<sup><sup>[↩ Parent](#externaltablespecforproviderconnectionstringsindex)</sup></sup>



SecretKeyRef reads the connection string (with SAS token or account
key) from a Secret. Write-only: a changed secret is detected through a
hash, the value never appears in status, events or logs.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Key to select.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the secret; defaults to the namespace of the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.spec.forProvider.columns[index]
<sup><sup>[↩ Parent](#externaltablespecforprovider)</sup></sup>



Column is one column of a table schema.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the column.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type of the column.<br/>
          <br/>
            <i>Enum</i>: bool, boolean, datetime, date, dynamic, guid, uuid, uniqueid, int, long, real, double, decimal, string, timespan, time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring of the column.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.spec.forProvider.properties
<sup><sup>[↩ Parent](#externaltablespecforprovider)</sup></sup>



Properties are the optional "with (...)" settings.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>compressed</b></td>
        <td>boolean</td>
        <td>
          Compressed marks the files as compressed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>compressionType</b></td>
        <td>string</td>
        <td>
          CompressionType, e.g. gzip or snappy.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>docString</b></td>
        <td>string</td>
        <td>
          DocString of the external table.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>encoding</b></td>
        <td>string</td>
        <td>
          Encoding of text files, e.g. UTF8NoBOM.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>fileExtension</b></td>
        <td>string</td>
        <td>
          FileExtension of the files, e.g. ".parquet".<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>folder</b></td>
        <td>string</td>
        <td>
          Folder of the external table in the database tree.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>includeHeaders</b></td>
        <td>enum</td>
        <td>
          IncludeHeaders for delimited formats: None, All or FirstFile.<br/>
          <br/>
            <i>Enum</i>: None, All, FirstFile<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namePrefix</b></td>
        <td>string</td>
        <td>
          NamePrefix of exported files.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.spec.providerConfigRef
<sup><sup>[↩ Parent](#externaltablespec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ExternalTable.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#externaltablespec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ExternalTable.status
<sup><sup>[↩ Parent](#externaltable)</sup></sup>



An ExternalTableStatus represents the observed state of an ExternalTable.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#externaltablestatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          ExternalTableObservation are the observable fields of an ExternalTable.
Connection strings are reported as URI only (everything after the first
';' is a secret and is masked by Kusto as well).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#externaltablestatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.status.atProvider
<sup><sup>[↩ Parent](#externaltablestatus)</sup></sup>



ExternalTableObservation are the observable fields of an ExternalTable.
Connection strings are reported as URI only (everything after the first
';' is a secret and is masked by Kusto as well).

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#externaltablestatusatprovidercolumnsindex">columns</a></b></td>
        <td>[]object</td>
        <td>
          Columns as reported by the cluster (only loaded when the spec has columns).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>connectionStringUris</b></td>
        <td>[]string</td>
        <td>
          ConnectionStringURIs are the URI parts of the connection strings.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>dataFormat</b></td>
        <td>string</td>
        <td>
          DataFormat as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>docString</b></td>
        <td>string</td>
        <td>
          DocString as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>folder</b></td>
        <td>string</td>
        <td>
          Folder as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind as reported by the cluster (Blob, Adl, Delta, ...).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>partitions</b></td>
        <td>string</td>
        <td>
          Partitions is the partition definition JSON as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pathFormat</b></td>
        <td>string</td>
        <td>
          PathFormat as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.status.atProvider.columns[index]
<sup><sup>[↩ Parent](#externaltablestatusatprovider)</sup></sup>



ObservedColumn is a column as reported by the cluster.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the column.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of the column (canonical Kusto type).<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring of the column.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExternalTable.status.conditions[index]
<sup><sup>[↩ Parent](#externaltablestatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## Function
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






A Function is a stored Kusto function (or view). Whoever may write a
Function runs arbitrary KQL with the provider principal's permissions.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>Function</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#functionspec">spec</a></b></td>
        <td>object</td>
        <td>
          A FunctionSpec defines the desired state of a Function.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#functionstatus">status</a></b></td>
        <td>object</td>
        <td>
          A FunctionStatus represents the observed state of a Function.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Function.spec
<sup><sup>[↩ Parent](#function)</sup></sup>



A FunctionSpec defines the desired state of a Function.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#functionspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          FunctionParameters are the configurable fields of a Function.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#functionspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#functionspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Function.spec.forProvider
<sup><sup>[↩ Parent](#functionspec)</sup></sup>



FunctionParameters are the configurable fields of a Function.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>body</b></td>
        <td>string</td>
        <td>
          Body of the function without the outer curly braces. It is KQL code and
is sent verbatim; the provider adds the braces.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the function. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring of the function.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>folder</b></td>
        <td>string</td>
        <td>
          Folder of the function in the database tree.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the function in Kusto. Written to the crossplane.io/external-name
annotation on the first reconcile if that annotation is empty. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: name is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parameters</b></td>
        <td>string</td>
        <td>
          Parameters is the raw parameter list including parentheses, e.g.
"(limit:long = 100, T:(x:long))". Whitespace and type aliases are
normalized before comparison.<br/>
          <br/>
            <i>Default</i>: ()<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>skipValidation</b></td>
        <td>boolean</td>
        <td>
          SkipValidation skips semantic validation of the body on write. Write-only.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>view</b></td>
        <td>boolean</td>
        <td>
          View marks the function as a view (usable in wildcard unions). The flag
is write-only: Kusto does not report it in .show functions.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Function.spec.providerConfigRef
<sup><sup>[↩ Parent](#functionspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Function.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#functionspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Function.status
<sup><sup>[↩ Parent](#function)</sup></sup>



A FunctionStatus represents the observed state of a Function.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#functionstatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          FunctionObservation are the observable fields of a Function.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#functionstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Function.status.atProvider
<sup><sup>[↩ Parent](#functionstatus)</sup></sup>



FunctionObservation are the observable fields of a Function.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>body</b></td>
        <td>string</td>
        <td>
          Body as reported by the cluster (Kusto may reformat it).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>folder</b></td>
        <td>string</td>
        <td>
          Folder as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parameters</b></td>
        <td>string</td>
        <td>
          Parameters as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Function.status.conditions[index]
<sup><sup>[↩ Parent](#functionstatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## IngestionMapping
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






An IngestionMapping is a pre-created ingestion mapping of a table. Its
identity is (database, table, kind, name).

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>IngestionMapping</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingspec">spec</a></b></td>
        <td>object</td>
        <td>
          An IngestionMappingSpec defines the desired state of an IngestionMapping.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingstatus">status</a></b></td>
        <td>object</td>
        <td>
          An IngestionMappingStatus represents the observed state of an IngestionMapping.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.spec
<sup><sup>[↩ Parent](#ingestionmapping)</sup></sup>



An IngestionMappingSpec defines the desired state of an IngestionMapping.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ingestionmappingspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          IngestionMappingParameters are the configurable fields of an IngestionMapping.<br/>
          <br/>
            <i>Validations</i>:<li>has(self.table) || has(self.tableRef) || has(self.tableSelector): table, tableRef or tableSelector is required</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.spec.forProvider
<sup><sup>[↩ Parent](#ingestionmappingspec)</sup></sup>



IngestionMappingParameters are the configurable fields of an IngestionMapping.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the table. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind is the data format of the mapping. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: kind is immutable</li>
            <i>Enum</i>: csv, json, avro, apacheavro, parquet, orc, w3clogfile<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingspecforprovidermappingindex">mapping</a></b></td>
        <td>[]object</td>
        <td>
          Mapping is the ordered list of column mappings.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the mapping in Kusto. Written to the crossplane.io/external-name
annotation on the first reconcile if that annotation is empty. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: name is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>table</b></td>
        <td>string</td>
        <td>
          Table the mapping belongs to. Immutable once set.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: table is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingspecforprovidertableref">tableRef</a></b></td>
        <td>object</td>
        <td>
          TableRef references a Table managed resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingspecforprovidertableselector">tableSelector</a></b></td>
        <td>object</td>
        <td>
          TableSelector selects a Table managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.spec.forProvider.mapping[index]
<sup><sup>[↩ Parent](#ingestionmappingspecforprovider)</sup></sup>



IngestionMappingColumn maps one source field to a table column.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>column</b></td>
        <td>string</td>
        <td>
          Column is the target column of the table.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>dataType</b></td>
        <td>enum</td>
        <td>
          DataType of the target column. Optional; Kusto derives it from the
table when omitted.<br/>
          <br/>
            <i>Enum</i>: bool, boolean, datetime, date, dynamic, guid, uuid, uniqueid, int, long, real, double, decimal, string, timespan, time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>properties</b></td>
        <td>map[string]string</td>
        <td>
          Properties are the format specific mapping properties, e.g. path
("$.ts"), transform, ordinal, constValue or field. Keys are case
insensitive and written to Kusto with an upper-case first letter.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.spec.forProvider.tableRef
<sup><sup>[↩ Parent](#ingestionmappingspecforprovider)</sup></sup>



TableRef references a Table managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingspecforprovidertablerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.spec.forProvider.tableRef.policy
<sup><sup>[↩ Parent](#ingestionmappingspecforprovidertableref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.spec.forProvider.tableSelector
<sup><sup>[↩ Parent](#ingestionmappingspecforprovider)</sup></sup>



TableSelector selects a Table managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingspecforprovidertableselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.spec.forProvider.tableSelector.policy
<sup><sup>[↩ Parent](#ingestionmappingspecforprovidertableselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.spec.providerConfigRef
<sup><sup>[↩ Parent](#ingestionmappingspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### IngestionMapping.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#ingestionmappingspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### IngestionMapping.status
<sup><sup>[↩ Parent](#ingestionmapping)</sup></sup>



An IngestionMappingStatus represents the observed state of an IngestionMapping.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ingestionmappingstatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          IngestionMappingObservation are the observable fields of an IngestionMapping.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionmappingstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.status.atProvider
<sup><sup>[↩ Parent](#ingestionmappingstatus)</sup></sup>



IngestionMappingObservation are the observable fields of an IngestionMapping.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastUpdatedOn</b></td>
        <td>string</td>
        <td>
          LastUpdatedOn as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mapping</b></td>
        <td>string</td>
        <td>
          Mapping is the mapping JSON as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>table</b></td>
        <td>string</td>
        <td>
          Table as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionMapping.status.conditions[index]
<sup><sup>[↩ Parent](#ingestionmappingstatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## MaterializedView
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






A MaterializedView is a Kusto materialized view. With backfill the view is
created asynchronously; the resource stays in Creating until the backfill
operation completes. Deleting the managed resource drops the view.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>MaterializedView</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#materializedviewspec">spec</a></b></td>
        <td>object</td>
        <td>
          A MaterializedViewSpec defines the desired state of a MaterializedView.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#materializedviewstatus">status</a></b></td>
        <td>object</td>
        <td>
          A MaterializedViewStatus represents the observed state of a MaterializedView.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MaterializedView.spec
<sup><sup>[↩ Parent](#materializedview)</sup></sup>



A MaterializedViewSpec defines the desired state of a MaterializedView.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#materializedviewspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          MaterializedViewParameters are the configurable fields of a MaterializedView.<br/>
          <br/>
            <i>Validations</i>:<li>has(self.sourceTable) || has(self.sourceTableRef) || has(self.sourceTableSelector) || has(self.sourceMaterializedView): one of sourceTable, sourceTableRef, sourceTableSelector or sourceMaterializedView is required</li><li>!(has(self.sourceMaterializedView) && (has(self.sourceTable) || has(self.sourceTableRef) || has(self.sourceTableSelector))): sourceMaterializedView and sourceTable are mutually exclusive</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#materializedviewspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#materializedviewspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MaterializedView.spec.forProvider
<sup><sup>[↩ Parent](#materializedviewspec)</sup></sup>



MaterializedViewParameters are the configurable fields of a MaterializedView.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the view. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>query</b></td>
        <td>string</td>
        <td>
          Query is the KQL aggregation query over the source. It is sent verbatim.
Kusto only accepts limited changes to an existing view's query (no
change of the group-by expressions, column names or types); rejected
changes surface as Synced=False with the Kusto message.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>allowMaterializedViewsWithoutRowLevelSecurity</b></td>
        <td>boolean</td>
        <td>
          AllowMaterializedViewsWithoutRowLevelSecurity allows creating the view
over a source table with a row level security policy. Write-only.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>autoUpdateSchema</b></td>
        <td>boolean</td>
        <td>
          AutoUpdateSchema propagates source schema changes to the view.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>backfill</b></td>
        <td>boolean</td>
        <td>
          Backfill materializes existing source data. Create-only: the view is
created asynchronously and the backfill operation is polled until it
completes.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: backfill is create-only</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>concurrency</b></td>
        <td>integer</td>
        <td>
          Concurrency of the backfill. Create-only.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: concurrency is create-only</li>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>dimensionTables</b></td>
        <td>[]string</td>
        <td>
          DimensionTables are tables joined in the query that must not be treated
as fact tables. Write-only: Kusto does not report them.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring of the view.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>effectiveDateTime</b></td>
        <td>string</td>
        <td>
          EffectiveDateTime (ISO 8601) limits backfill to records ingested after
this time. Create-only.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: effectiveDateTime is create-only</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled toggles materialization (.enable/.disable materialized-view).
Defaults to true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>folder</b></td>
        <td>string</td>
        <td>
          Folder of the view in the database tree.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lookback</b></td>
        <td>string</td>
        <td>
          Lookback limits the period of source data considered during
materialization, e.g. "6h".<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lookbackColumn</b></td>
        <td>string</td>
        <td>
          LookbackColumn is the datetime column the lookback applies to. Kusto
does not allow changing it after it has been set.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxSourceRecordsForSingleIngest</b></td>
        <td>integer</td>
        <td>
          MaxSourceRecordsForSingleIngest caps the records per backfill ingest
operation. Create-only.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: maxSourceRecordsForSingleIngest is create-only</li>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the materialized view in Kusto. Written to the
crossplane.io/external-name annotation on the first reconcile if that
annotation is empty. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: name is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sourceMaterializedView</b></td>
        <td>string</td>
        <td>
          SourceMaterializedView is another materialized view used as source
(materialized view over materialized view). Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: sourceMaterializedView is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sourceTable</b></td>
        <td>string</td>
        <td>
          SourceTable the view is defined over. Immutable once set; changing the
source means a new view.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: sourceTable is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#materializedviewspecforprovidersourcetableref">sourceTableRef</a></b></td>
        <td>object</td>
        <td>
          SourceTableRef references a Table managed resource as source.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#materializedviewspecforprovidersourcetableselector">sourceTableSelector</a></b></td>
        <td>object</td>
        <td>
          SourceTableSelector selects a Table managed resource as source.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>updateExtentsCreationTime</b></td>
        <td>boolean</td>
        <td>
          UpdateExtentsCreationTime sets the creation time of backfilled extents
to the ingestion time of the source records. Create-only.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: updateExtentsCreationTime is create-only</li>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MaterializedView.spec.forProvider.sourceTableRef
<sup><sup>[↩ Parent](#materializedviewspecforprovider)</sup></sup>



SourceTableRef references a Table managed resource as source.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#materializedviewspecforprovidersourcetablerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MaterializedView.spec.forProvider.sourceTableRef.policy
<sup><sup>[↩ Parent](#materializedviewspecforprovidersourcetableref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MaterializedView.spec.forProvider.sourceTableSelector
<sup><sup>[↩ Parent](#materializedviewspecforprovider)</sup></sup>



SourceTableSelector selects a Table managed resource as source.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#materializedviewspecforprovidersourcetableselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MaterializedView.spec.forProvider.sourceTableSelector.policy
<sup><sup>[↩ Parent](#materializedviewspecforprovidersourcetableselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MaterializedView.spec.providerConfigRef
<sup><sup>[↩ Parent](#materializedviewspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### MaterializedView.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#materializedviewspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### MaterializedView.status
<sup><sup>[↩ Parent](#materializedview)</sup></sup>



A MaterializedViewStatus represents the observed state of a MaterializedView.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#materializedviewstatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          MaterializedViewObservation are the observable fields of a MaterializedView.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#materializedviewstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MaterializedView.status.atProvider
<sup><sup>[↩ Parent](#materializedviewstatus)</sup></sup>



MaterializedViewObservation are the observable fields of a MaterializedView.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>effectiveDateTime</b></td>
        <td>string</td>
        <td>
          EffectiveDateTime as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>folder</b></td>
        <td>string</td>
        <td>
          Folder as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isEnabled</b></td>
        <td>boolean</td>
        <td>
          IsEnabled as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isHealthy</b></td>
        <td>boolean</td>
        <td>
          IsHealthy as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastRunResult</b></td>
        <td>string</td>
        <td>
          LastRunResult of the materialization as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lookback</b></td>
        <td>string</td>
        <td>
          Lookback as reported by the cluster (.NET timespan format).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>operationId</b></td>
        <td>string</td>
        <td>
          OperationID of a running backfill operation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>operationState</b></td>
        <td>string</td>
        <td>
          OperationState of the backfill operation (InProgress, Scheduled, ...).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>query</b></td>
        <td>string</td>
        <td>
          Query as reported by the cluster (Kusto may reformat it).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sourceTable</b></td>
        <td>string</td>
        <td>
          SourceTable as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MaterializedView.status.conditions[index]
<sup><sup>[↩ Parent](#materializedviewstatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ProviderConfig
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






A ProviderConfig configures how the provider connects to one Azure Data
Explorer cluster: the endpoint plus the identity used for management commands.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ProviderConfig</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#providerconfigspec">spec</a></b></td>
        <td>object</td>
        <td>
          A ProviderConfigSpec defines the desired state of a ProviderConfig.<br/>
          <br/>
            <i>Validations</i>:<li>self.credentials.source != 'None' || self.clusterUri.startsWith('http://'): credentials.source None is only allowed for http:// endpoints (Kusto emulator)</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#providerconfigstatus">status</a></b></td>
        <td>object</td>
        <td>
          A ProviderConfigStatus defines the status of a Provider.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProviderConfig.spec
<sup><sup>[↩ Parent](#providerconfig)</sup></sup>



A ProviderConfigSpec defines the desired state of a ProviderConfig.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>clusterUri</b></td>
        <td>string</td>
        <td>
          ClusterURI is the Kusto query endpoint of the cluster, e.g.
https://mycluster.westeurope.kusto.windows.net. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: clusterUri is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#providerconfigspeccredentials">credentials</a></b></td>
        <td>object</td>
        <td>
          Credentials required to authenticate to this provider.<br/>
          <br/>
            <i>Validations</i>:<li>self.source != 'Secret' || has(self.secretRef): secretRef is required when source is Secret</li><li>self.source != 'ManagedIdentity' || has(self.managedIdentity): managedIdentity is required when source is ManagedIdentity</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>azureEnvironment</b></td>
        <td>enum</td>
        <td>
          AzureEnvironment selects the sovereign cloud. Only AzurePublicCloud is
tested; the others are wired through to the SDK but unverified.<br/>
          <br/>
            <i>Enum</i>: AzurePublicCloud, AzureChinaCloud, AzureUSGovernment<br/>
            <i>Default</i>: AzurePublicCloud<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProviderConfig.spec.credentials
<sup><sup>[↩ Parent](#providerconfigspec)</sup></sup>



Credentials required to authenticate to this provider.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>source</b></td>
        <td>enum</td>
        <td>
          Source of the provider credentials.<br/>
          <br/>
            <i>Enum</i>: None, Secret, WorkloadIdentity, ManagedIdentity<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#providerconfigspeccredentialsmanagedidentity">managedIdentity</a></b></td>
        <td>object</td>
        <td>
          ManagedIdentity settings, used when source is ManagedIdentity.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#providerconfigspeccredentialssecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          SecretRef points to a Secret key holding a JSON document with the keys
clientId, clientSecret and tenantId. The format is compatible with the
credentials Secret of provider-upbound-azure (additional keys are ignored).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#providerconfigspeccredentialsworkloadidentity">workloadIdentity</a></b></td>
        <td>object</td>
        <td>
          WorkloadIdentity settings, used when source is WorkloadIdentity.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProviderConfig.spec.credentials.managedIdentity
<sup><sup>[↩ Parent](#providerconfigspeccredentials)</sup></sup>



ManagedIdentity settings, used when source is ManagedIdentity.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type of the managed identity.<br/>
          <br/>
            <i>Enum</i>: SystemAssigned, UserAssigned<br/>
            <i>Default</i>: SystemAssigned<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>clientId</b></td>
        <td>string</td>
        <td>
          ClientID of a user-assigned identity. Mutually exclusive with resourceId.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resourceId</b></td>
        <td>string</td>
        <td>
          ResourceID of a user-assigned identity. Mutually exclusive with clientId.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProviderConfig.spec.credentials.secretRef
<sup><sup>[↩ Parent](#providerconfigspeccredentials)</sup></sup>



SecretRef points to a Secret key holding a JSON document with the keys
clientId, clientSecret and tenantId. The format is compatible with the
credentials Secret of provider-upbound-azure (additional keys are ignored).

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key to select.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ProviderConfig.spec.credentials.workloadIdentity
<sup><sup>[↩ Parent](#providerconfigspeccredentials)</sup></sup>



WorkloadIdentity settings, used when source is WorkloadIdentity.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>clientId</b></td>
        <td>string</td>
        <td>
          ClientID of the federated application.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenantId</b></td>
        <td>string</td>
        <td>
          TenantID of the federated application.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tokenFile</b></td>
        <td>string</td>
        <td>
          TokenFile is the path of the projected service account token.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProviderConfig.status
<sup><sup>[↩ Parent](#providerconfig)</sup></sup>



A ProviderConfigStatus defines the status of a Provider.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#providerconfigstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>users</b></td>
        <td>integer</td>
        <td>
          Users of this provider configuration.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProviderConfig.status.conditions[index]
<sup><sup>[↩ Parent](#providerconfigstatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ProviderConfigUsage
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






A ProviderConfigUsage indicates that a resource is using a ProviderConfig or a
ClusterProviderConfig. There is deliberately no cluster scoped usage type, Usages always live in
the namespace of the MR that created them, and record which kind of config they refer to.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ProviderConfigUsage</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#providerconfigusageproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference to the provider config being used.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#providerconfigusageresourceref">resourceRef</a></b></td>
        <td>object</td>
        <td>
          ResourceReference to the managed resource using the provider config.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ProviderConfigUsage.providerConfigRef
<sup><sup>[↩ Parent](#providerconfigusage)</sup></sup>



ProviderConfigReference to the provider config being used.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ProviderConfigUsage.resourceRef
<sup><sup>[↩ Parent](#providerconfigusage)</sup></sup>



ResourceReference to the managed resource using the provider config.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          APIVersion of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>uid</b></td>
        <td>string</td>
        <td>
          UID of the referenced object.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## Table
<sup><sup>[↩ Parent](#adxfunctionalteamv1alpha1 )</sup></sup>






A Table is a Kusto table with a managed schema. Deleting the managed
resource drops the table and all of its data (Crossplane default deletion
semantics); use deletionPolicy Orphan or managementPolicies to keep it.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>Table</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#tablespec">spec</a></b></td>
        <td>object</td>
        <td>
          A TableSpec defines the desired state of a Table.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#tablestatus">status</a></b></td>
        <td>object</td>
        <td>
          A TableStatus represents the observed state of a Table.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Table.spec
<sup><sup>[↩ Parent](#table)</sup></sup>



A TableSpec defines the desired state of a Table.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#tablespecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          TableParameters are the configurable fields of a Table.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#tablespecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#tablespecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Table.spec.forProvider
<sup><sup>[↩ Parent](#tablespec)</sup></sup>



TableParameters are the configurable fields of a Table.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#tablespecforprovidercolumnsindex">columns</a></b></td>
        <td>[]object</td>
        <td>
          Columns of the table in order.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the table. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring of the table.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>folder</b></td>
        <td>string</td>
        <td>
          Folder of the table in the database tree.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the table in Kusto. Written to the crossplane.io/external-name
annotation on the first reconcile if that annotation is empty; the
annotation remains the source of truth. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: name is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>schemaUpdateMode</b></td>
        <td>enum</td>
        <td>
          SchemaUpdateMode selects Merge (default, add-only, safe) or Replace
(authoritative, drops columns and their data).<br/>
          <br/>
            <i>Enum</i>: Merge, Replace<br/>
            <i>Default</i>: Merge<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Table.spec.forProvider.columns[index]
<sup><sup>[↩ Parent](#tablespecforprovider)</sup></sup>



Column is one column of a table schema.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the column.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type of the column.<br/>
          <br/>
            <i>Enum</i>: bool, boolean, datetime, date, dynamic, guid, uuid, uniqueid, int, long, real, double, decimal, string, timespan, time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring of the column.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Table.spec.providerConfigRef
<sup><sup>[↩ Parent](#tablespec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Table.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#tablespec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Table.status
<sup><sup>[↩ Parent](#table)</sup></sup>



A TableStatus represents the observed state of a Table.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#tablestatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          TableObservation are the observable fields of a Table.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#tablestatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Table.status.atProvider
<sup><sup>[↩ Parent](#tablestatus)</sup></sup>



TableObservation are the observable fields of a Table.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#tablestatusatprovidercolumnsindex">columns</a></b></td>
        <td>[]object</td>
        <td>
          Columns as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>driftColumns</b></td>
        <td>[]string</td>
        <td>
          DriftColumns exist in the cluster but not in the spec. They are kept in
Merge mode and would be dropped in Replace mode.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>folder</b></td>
        <td>string</td>
        <td>
          Folder as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Table.status.atProvider.columns[index]
<sup><sup>[↩ Parent](#tablestatusatprovider)</sup></sup>



ObservedColumn is a column as reported by the cluster.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the column.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of the column (canonical Kusto type).<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>docstring</b></td>
        <td>string</td>
        <td>
          Docstring of the column.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Table.status.conditions[index]
<sup><sup>[↩ Parent](#tablestatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

# cluster.adx.functional.team/v1alpha1

Resource Types:

- [CalloutPolicy](#calloutpolicy)

- [CapacityPolicy](#capacitypolicy)

- [ClusterManagedIdentityPolicy](#clustermanagedidentitypolicy)

- [MultiDatabaseAdminsPolicy](#multidatabaseadminspolicy)

- [QueryWeakConsistencyPolicy](#queryweakconsistencypolicy)

- [RequestClassificationPolicy](#requestclassificationpolicy)

- [SandboxPolicy](#sandboxpolicy)

- [WorkloadGroup](#workloadgroup)




## CalloutPolicy
<sup><sup>[↩ Parent](#clusteradxfunctionalteamv1alpha1 )</sup></sup>






A CalloutPolicy controls which external endpoints (SQL, Cosmos DB, external data, sandboxes, ...) the cluster may call. The list is authoritative for the cluster.
There is exactly one callout policy per cluster, so one managed resource
per cluster (ProviderConfig) is expected; two managed resources on the same
cluster policy fight over it and are a user error. The provider principal
needs the AllDatabasesAdmin cluster role for cluster policies.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>cluster.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>CalloutPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#calloutpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A CalloutPolicySpec defines the desired state of a CalloutPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#calloutpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A CalloutPolicyStatus represents the observed state of a CalloutPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CalloutPolicy.spec
<sup><sup>[↩ Parent](#calloutpolicy)</sup></sup>



A CalloutPolicySpec defines the desired state of a CalloutPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#calloutpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          CalloutPolicyParameters are the configurable fields of a CalloutPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#calloutpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#calloutpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CalloutPolicy.spec.forProvider
<sup><sup>[↩ Parent](#calloutpolicyspec)</sup></sup>



CalloutPolicyParameters are the configurable fields of a CalloutPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#calloutpolicyspecforprovidercalloutsindex">callouts</a></b></td>
        <td>[]object</td>
        <td>
          Callouts is the complete list of callout rules.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CalloutPolicy.spec.forProvider.callouts[index]
<sup><sup>[↩ Parent](#calloutpolicyspecforprovider)</sup></sup>



CalloutRule allows or denies callouts of one type to URIs matching a regex.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>calloutType</b></td>
        <td>enum</td>
        <td>
          CalloutType of the rule.<br/>
          <br/>
            <i>Enum</i>: kusto, sql, cosmosdb, external_data, azure_digital_twins, sandbox_artifacts, webapi, mysql, postgresql, genevametrics, azure_openai<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>calloutUriRegex</b></td>
        <td>string</td>
        <td>
          CalloutURIRegex matches the target URIs.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>canCall</b></td>
        <td>boolean</td>
        <td>
          CanCall allows (true) or denies (false) matching callouts.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### CalloutPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#calloutpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### CalloutPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#calloutpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### CalloutPolicy.status
<sup><sup>[↩ Parent](#calloutpolicy)</sup></sup>



A CalloutPolicyStatus represents the observed state of a CalloutPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#calloutpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#calloutpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CalloutPolicy.status.atProvider
<sup><sup>[↩ Parent](#calloutpolicystatus)</sup></sup>



ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by ".show cluster policy <name>".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CalloutPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#calloutpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## CapacityPolicy
<sup><sup>[↩ Parent](#clusteradxfunctionalteamv1alpha1 )</sup></sup>






A CapacityPolicy tunes the concurrency limits of cluster operations (ingestion, merge, export, materialized views, ...). The policy always exists on a cluster, so the managed resource only ever alters it; deleting the managed resource leaves the last applied values in place (Kusto has no delete for this policy).
There is exactly one capacity policy per cluster, so one managed resource
per cluster (ProviderConfig) is expected; two managed resources on the same
cluster policy fight over it and are a user error. The provider principal
needs the AllDatabasesAdmin cluster role for cluster policies.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>cluster.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>CapacityPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#capacitypolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A CapacityPolicySpec defines the desired state of a CapacityPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#capacitypolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A CapacityPolicyStatus represents the observed state of a CapacityPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CapacityPolicy.spec
<sup><sup>[↩ Parent](#capacitypolicy)</sup></sup>



A CapacityPolicySpec defines the desired state of a CapacityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#capacitypolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          CapacityPolicyParameters are the configurable fields of a CapacityPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#capacitypolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#capacitypolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CapacityPolicy.spec.forProvider
<sup><sup>[↩ Parent](#capacitypolicyspec)</sup></sup>



CapacityPolicyParameters are the configurable fields of a CapacityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>policy</b></td>
        <td>JSON</td>
        <td>
          Policy is the capacity policy JSON as documented by Kusto (IngestionCapacity, ExtentsMergeCapacity, ExportCapacity, MaterializedViewsCapacity, ...). Only the keys present here are compared with the cluster; unknown keys are passed through.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### CapacityPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#capacitypolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### CapacityPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#capacitypolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### CapacityPolicy.status
<sup><sup>[↩ Parent](#capacitypolicy)</sup></sup>



A CapacityPolicyStatus represents the observed state of a CapacityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#capacitypolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#capacitypolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CapacityPolicy.status.atProvider
<sup><sup>[↩ Parent](#capacitypolicystatus)</sup></sup>



ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by ".show cluster policy <name>".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CapacityPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#capacitypolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ClusterManagedIdentityPolicy
<sup><sup>[↩ Parent](#clusteradxfunctionalteamv1alpha1 )</sup></sup>






A ClusterManagedIdentityPolicy allows managed identities to be used for specific usages at the cluster level (all databases). The list is authoritative for the cluster.
There is exactly one managed_identity policy per cluster, so one managed resource
per cluster (ProviderConfig) is expected; two managed resources on the same
cluster policy fight over it and are a user error. The provider principal
needs the AllDatabasesAdmin cluster role for cluster policies.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>cluster.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ClusterManagedIdentityPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#clustermanagedidentitypolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A ClusterManagedIdentityPolicySpec defines the desired state of a ClusterManagedIdentityPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#clustermanagedidentitypolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A ClusterManagedIdentityPolicyStatus represents the observed state of a ClusterManagedIdentityPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterManagedIdentityPolicy.spec
<sup><sup>[↩ Parent](#clustermanagedidentitypolicy)</sup></sup>



A ClusterManagedIdentityPolicySpec defines the desired state of a ClusterManagedIdentityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#clustermanagedidentitypolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          ClusterManagedIdentityPolicyParameters are the configurable fields of a ClusterManagedIdentityPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#clustermanagedidentitypolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#clustermanagedidentitypolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterManagedIdentityPolicy.spec.forProvider
<sup><sup>[↩ Parent](#clustermanagedidentitypolicyspec)</sup></sup>



ClusterManagedIdentityPolicyParameters are the configurable fields of a ClusterManagedIdentityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#clustermanagedidentitypolicyspecforprovideridentitiesindex">identities</a></b></td>
        <td>[]object</td>
        <td>
          Identities is the complete list of allowed identities.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterManagedIdentityPolicy.spec.forProvider.identities[index]
<sup><sup>[↩ Parent](#clustermanagedidentitypolicyspecforprovider)</sup></sup>



ClusterManagedIdentityEntry allows one identity for usages.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>allowedUsages</b></td>
        <td>[]string</td>
        <td>
          AllowedUsages of the identity (e.g. NativeIngestion, ExternalTable, DataConnection, AutomatedFlows, SandboxArtifacts, AzureAI, All).<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>objectId</b></td>
        <td>string</td>
        <td>
          ObjectID of the managed identity, or "system".<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ClusterManagedIdentityPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#clustermanagedidentitypolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ClusterManagedIdentityPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#clustermanagedidentitypolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ClusterManagedIdentityPolicy.status
<sup><sup>[↩ Parent](#clustermanagedidentitypolicy)</sup></sup>



A ClusterManagedIdentityPolicyStatus represents the observed state of a ClusterManagedIdentityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#clustermanagedidentitypolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#clustermanagedidentitypolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterManagedIdentityPolicy.status.atProvider
<sup><sup>[↩ Parent](#clustermanagedidentitypolicystatus)</sup></sup>



ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by ".show cluster policy <name>".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ClusterManagedIdentityPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#clustermanagedidentitypolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## MultiDatabaseAdminsPolicy
<sup><sup>[↩ Parent](#clusteradxfunctionalteamv1alpha1 )</sup></sup>






A MultiDatabaseAdminsPolicy lists principals that may administer several databases without holding AllDatabasesAdmin. The policy JSON shape is not verified against a cluster (spike), so the comparison relies on the hash annotations.
There is exactly one multidatabaseadmins policy per cluster, so one managed resource
per cluster (ProviderConfig) is expected; two managed resources on the same
cluster policy fight over it and are a user error. The provider principal
needs the AllDatabasesAdmin cluster role for cluster policies.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>cluster.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>MultiDatabaseAdminsPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#multidatabaseadminspolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A MultiDatabaseAdminsPolicySpec defines the desired state of a MultiDatabaseAdminsPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#multidatabaseadminspolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A MultiDatabaseAdminsPolicyStatus represents the observed state of a MultiDatabaseAdminsPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MultiDatabaseAdminsPolicy.spec
<sup><sup>[↩ Parent](#multidatabaseadminspolicy)</sup></sup>



A MultiDatabaseAdminsPolicySpec defines the desired state of a MultiDatabaseAdminsPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#multidatabaseadminspolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          MultiDatabaseAdminsPolicyParameters are the configurable fields of a MultiDatabaseAdminsPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#multidatabaseadminspolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#multidatabaseadminspolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MultiDatabaseAdminsPolicy.spec.forProvider
<sup><sup>[↩ Parent](#multidatabaseadminspolicyspec)</sup></sup>



MultiDatabaseAdminsPolicyParameters are the configurable fields of a MultiDatabaseAdminsPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>principals</b></td>
        <td>[]string</td>
        <td>
          Principals is the complete list of allowed principals.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MultiDatabaseAdminsPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#multidatabaseadminspolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### MultiDatabaseAdminsPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#multidatabaseadminspolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### MultiDatabaseAdminsPolicy.status
<sup><sup>[↩ Parent](#multidatabaseadminspolicy)</sup></sup>



A MultiDatabaseAdminsPolicyStatus represents the observed state of a MultiDatabaseAdminsPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#multidatabaseadminspolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#multidatabaseadminspolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MultiDatabaseAdminsPolicy.status.atProvider
<sup><sup>[↩ Parent](#multidatabaseadminspolicystatus)</sup></sup>



ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by ".show cluster policy <name>".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MultiDatabaseAdminsPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#multidatabaseadminspolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## QueryWeakConsistencyPolicy
<sup><sup>[↩ Parent](#clusteradxfunctionalteamv1alpha1 )</sup></sup>






A QueryWeakConsistencyPolicy configures weak consistency query nodes. The policy always exists on a cluster; deleting the managed resource leaves the last applied values in place (Kusto has no delete for this policy).
There is exactly one query_weak_consistency policy per cluster, so one managed resource
per cluster (ProviderConfig) is expected; two managed resources on the same
cluster policy fight over it and are a user error. The provider principal
needs the AllDatabasesAdmin cluster role for cluster policies.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>cluster.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>QueryWeakConsistencyPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#queryweakconsistencypolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A QueryWeakConsistencyPolicySpec defines the desired state of a QueryWeakConsistencyPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#queryweakconsistencypolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A QueryWeakConsistencyPolicyStatus represents the observed state of a QueryWeakConsistencyPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryWeakConsistencyPolicy.spec
<sup><sup>[↩ Parent](#queryweakconsistencypolicy)</sup></sup>



A QueryWeakConsistencyPolicySpec defines the desired state of a QueryWeakConsistencyPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#queryweakconsistencypolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          QueryWeakConsistencyPolicyParameters are the configurable fields of a QueryWeakConsistencyPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryweakconsistencypolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryweakconsistencypolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryWeakConsistencyPolicy.spec.forProvider
<sup><sup>[↩ Parent](#queryweakconsistencypolicyspec)</sup></sup>



QueryWeakConsistencyPolicyParameters are the configurable fields of a QueryWeakConsistencyPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enableMetadataPrefetch</b></td>
        <td>boolean</td>
        <td>
          EnableMetadataPrefetch prefetches metadata on weak consistency nodes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maximumLagAllowedInMinutes</b></td>
        <td>integer</td>
        <td>
          MaximumLagAllowedInMinutes of metadata (-1 for default).<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maximumNumberOfNodes</b></td>
        <td>integer</td>
        <td>
          MaximumNumberOfNodes for weak consistency (-1 for default).<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>minimumNumberOfNodes</b></td>
        <td>integer</td>
        <td>
          MinimumNumberOfNodes for weak consistency (-1 for default).<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>percentageOfNodes</b></td>
        <td>integer</td>
        <td>
          PercentageOfNodes that serve weakly consistent queries (-1 for default).<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>refreshPeriodInSeconds</b></td>
        <td>integer</td>
        <td>
          RefreshPeriodInSeconds of metadata (-1 for default).<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>superSlackerNumberOfNodesThreshold</b></td>
        <td>integer</td>
        <td>
          SuperSlackerNumberOfNodesThreshold (-1 for default).<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryWeakConsistencyPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#queryweakconsistencypolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### QueryWeakConsistencyPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#queryweakconsistencypolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### QueryWeakConsistencyPolicy.status
<sup><sup>[↩ Parent](#queryweakconsistencypolicy)</sup></sup>



A QueryWeakConsistencyPolicyStatus represents the observed state of a QueryWeakConsistencyPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#queryweakconsistencypolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryweakconsistencypolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryWeakConsistencyPolicy.status.atProvider
<sup><sup>[↩ Parent](#queryweakconsistencypolicystatus)</sup></sup>



ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by ".show cluster policy <name>".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryWeakConsistencyPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#queryweakconsistencypolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## RequestClassificationPolicy
<sup><sup>[↩ Parent](#clusteradxfunctionalteamv1alpha1 )</sup></sup>






A RequestClassificationPolicy routes requests to workload groups through a classification function (KQL). The query is compared with the two stage KQL normalization.
There is exactly one request_classification policy per cluster, so one managed resource
per cluster (ProviderConfig) is expected; two managed resources on the same
cluster policy fight over it and are a user error. The provider principal
needs the AllDatabasesAdmin cluster role for cluster policies.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>cluster.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>RequestClassificationPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#requestclassificationpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A RequestClassificationPolicySpec defines the desired state of a RequestClassificationPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#requestclassificationpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A RequestClassificationPolicyStatus represents the observed state of a RequestClassificationPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RequestClassificationPolicy.spec
<sup><sup>[↩ Parent](#requestclassificationpolicy)</sup></sup>



A RequestClassificationPolicySpec defines the desired state of a RequestClassificationPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#requestclassificationpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          RequestClassificationPolicyParameters are the configurable fields of a RequestClassificationPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#requestclassificationpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#requestclassificationpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RequestClassificationPolicy.spec.forProvider
<sup><sup>[↩ Parent](#requestclassificationpolicyspec)</sup></sup>



RequestClassificationPolicyParameters are the configurable fields of a RequestClassificationPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>query</b></td>
        <td>string</td>
        <td>
          Query is the classification function body, e.g. iff(request_properties.current_principal == "aadapp=...", "Ingestion", "default").<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled toggles the policy. Defaults to true.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RequestClassificationPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#requestclassificationpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RequestClassificationPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#requestclassificationpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RequestClassificationPolicy.status
<sup><sup>[↩ Parent](#requestclassificationpolicy)</sup></sup>



A RequestClassificationPolicyStatus represents the observed state of a RequestClassificationPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#requestclassificationpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#requestclassificationpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RequestClassificationPolicy.status.atProvider
<sup><sup>[↩ Parent](#requestclassificationpolicystatus)</sup></sup>



ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by ".show cluster policy <name>".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RequestClassificationPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#requestclassificationpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## SandboxPolicy
<sup><sup>[↩ Parent](#clusteradxfunctionalteamv1alpha1 )</sup></sup>






A SandboxPolicy configures the sandboxes used by the python() and r() plugins. The list is authoritative for the cluster.
There is exactly one sandbox policy per cluster, so one managed resource
per cluster (ProviderConfig) is expected; two managed resources on the same
cluster policy fight over it and are a user error. The provider principal
needs the AllDatabasesAdmin cluster role for cluster policies.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>cluster.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>SandboxPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#sandboxpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A SandboxPolicySpec defines the desired state of a SandboxPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#sandboxpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A SandboxPolicyStatus represents the observed state of a SandboxPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SandboxPolicy.spec
<sup><sup>[↩ Parent](#sandboxpolicy)</sup></sup>



A SandboxPolicySpec defines the desired state of a SandboxPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#sandboxpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          SandboxPolicyParameters are the configurable fields of a SandboxPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#sandboxpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#sandboxpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SandboxPolicy.spec.forProvider
<sup><sup>[↩ Parent](#sandboxpolicyspec)</sup></sup>



SandboxPolicyParameters are the configurable fields of a SandboxPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#sandboxpolicyspecforprovidersandboxesindex">sandboxes</a></b></td>
        <td>[]object</td>
        <td>
          Sandboxes is the complete list of sandbox settings.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SandboxPolicy.spec.forProvider.sandboxes[index]
<sup><sup>[↩ Parent](#sandboxpolicyspecforprovider)</sup></sup>



SandboxRule configures one sandbox kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>sandboxKind</b></td>
        <td>enum</td>
        <td>
          SandboxKind of the rule.<br/>
          <br/>
            <i>Enum</i>: PythonExecution, RExecution<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>initializeOnStartup</b></td>
        <td>boolean</td>
        <td>
          InitializeOnStartup creates sandboxes eagerly on node startup.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxCpuPerSandbox</b></td>
        <td>integer</td>
        <td>
          MaxCpuPerSandbox in percent of a core.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxMemoryMbPerSandbox</b></td>
        <td>integer</td>
        <td>
          MaxMemoryMbPerSandbox in MB.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxNumberOfSandboxes</b></td>
        <td>integer</td>
        <td>
          MaxNumberOfSandboxes per node.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>virtualMachineSize</b></td>
        <td>string</td>
        <td>
          VirtualMachineSize the sandbox VM size the settings apply to; empty for all sizes.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SandboxPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#sandboxpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### SandboxPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#sandboxpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### SandboxPolicy.status
<sup><sup>[↩ Parent](#sandboxpolicy)</sup></sup>



A SandboxPolicyStatus represents the observed state of a SandboxPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#sandboxpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#sandboxpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SandboxPolicy.status.atProvider
<sup><sup>[↩ Parent](#sandboxpolicystatus)</sup></sup>



ClusterPolicyObservation is the observed state shared by all cluster
policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by ".show cluster policy <name>".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SandboxPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#sandboxpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## WorkloadGroup
<sup><sup>[↩ Parent](#clusteradxfunctionalteamv1alpha1 )</sup></sup>






A WorkloadGroup is a Kusto workload group (request limits, rate limits,
queuing). The built-in groups "default" and "internal" can be altered but
not dropped: deleting their managed resource leaves them in place. The
provider principal needs the AllDatabasesAdmin cluster role.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>cluster.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>WorkloadGroup</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#workloadgroupspec">spec</a></b></td>
        <td>object</td>
        <td>
          A WorkloadGroupSpec defines the desired state of a WorkloadGroup.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#workloadgroupstatus">status</a></b></td>
        <td>object</td>
        <td>
          A WorkloadGroupStatus represents the observed state of a WorkloadGroup.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### WorkloadGroup.spec
<sup><sup>[↩ Parent](#workloadgroup)</sup></sup>



A WorkloadGroupSpec defines the desired state of a WorkloadGroup.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#workloadgroupspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          WorkloadGroupParameters are the configurable fields of a WorkloadGroup.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#workloadgroupspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#workloadgroupspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### WorkloadGroup.spec.forProvider
<sup><sup>[↩ Parent](#workloadgroupspec)</sup></sup>



WorkloadGroupParameters are the configurable fields of a WorkloadGroup.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>workloadGroup</b></td>
        <td>JSON</td>
        <td>
          WorkloadGroup is the workload group JSON as documented by Kusto:
RequestLimitsPolicy, RequestRateLimitPolicies,
RequestRateLimitsEnforcementPolicy and RequestQueuingPolicy. Only the
keys present here are compared with the cluster; unknown keys are
passed through unchanged.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the workload group in Kusto. Written to the
crossplane.io/external-name annotation on the first reconcile if that
annotation is empty. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: name is immutable</li>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### WorkloadGroup.spec.providerConfigRef
<sup><sup>[↩ Parent](#workloadgroupspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### WorkloadGroup.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#workloadgroupspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### WorkloadGroup.status
<sup><sup>[↩ Parent](#workloadgroup)</sup></sup>



A WorkloadGroupStatus represents the observed state of a WorkloadGroup.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#workloadgroupstatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          WorkloadGroupObservation are the observable fields of a WorkloadGroup.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#workloadgroupstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### WorkloadGroup.status.atProvider
<sup><sup>[↩ Parent](#workloadgroupstatus)</sup></sup>



WorkloadGroupObservation are the observable fields of a WorkloadGroup.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>workloadGroup</b></td>
        <td>string</td>
        <td>
          WorkloadGroup is the workload group JSON as reported by the cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### WorkloadGroup.status.conditions[index]
<sup><sup>[↩ Parent](#workloadgroupstatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

# policy.adx.functional.team/v1alpha1

Resource Types:

- [AutoDeletePolicy](#autodeletepolicy)

- [CachingPolicy](#cachingpolicy)

- [EncodingPolicy](#encodingpolicy)

- [ExtentTagsRetentionPolicy](#extenttagsretentionpolicy)

- [IngestionBatchingPolicy](#ingestionbatchingpolicy)

- [IngestionTimePolicy](#ingestiontimepolicy)

- [ManagedIdentityPolicy](#managedidentitypolicy)

- [MergePolicy](#mergepolicy)

- [PartitioningPolicy](#partitioningpolicy)

- [QueryAccelerationPolicy](#queryaccelerationpolicy)

- [RestrictedViewAccessPolicy](#restrictedviewaccesspolicy)

- [RetentionPolicy](#retentionpolicy)

- [RowLevelSecurityPolicy](#rowlevelsecuritypolicy)

- [RowOrderPolicy](#roworderpolicy)

- [ShardingPolicy](#shardingpolicy)

- [StreamingIngestionPolicy](#streamingingestionpolicy)

- [UpdatePolicy](#updatepolicy)




## AutoDeletePolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A AutoDeletePolicy drops the table at an expiry date.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy auto_delete"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AutoDeletePolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A AutoDeletePolicySpec defines the desired state of a AutoDeletePolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A AutoDeletePolicyStatus represents the observed state of a AutoDeletePolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.spec
<sup><sup>[↩ Parent](#autodeletepolicy)</sup></sup>



A AutoDeletePolicySpec defines the desired state of a AutoDeletePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#autodeletepolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          AutoDeletePolicyParameters are the configurable fields of a AutoDeletePolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table']: AutoDeletePolicy can only target Table</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.spec.forProvider
<sup><sup>[↩ Parent](#autodeletepolicyspec)</sup></sup>



AutoDeletePolicyParameters are the configurable fields of a AutoDeletePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>expiryDate</b></td>
        <td>string</td>
        <td>
          ExpiryDate (ISO 8601 date or datetime) after which the table is dropped.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>deleteIfNotEmpty</b></td>
        <td>boolean</td>
        <td>
          DeleteIfNotEmpty also drops the table when it still holds data.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#autodeletepolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#autodeletepolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#autodeletepolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#autodeletepolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#autodeletepolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#autodeletepolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#autodeletepolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.status
<sup><sup>[↩ Parent](#autodeletepolicy)</sup></sup>



A AutoDeletePolicyStatus represents the observed state of a AutoDeletePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#autodeletepolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#autodeletepolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.status.atProvider
<sup><sup>[↩ Parent](#autodeletepolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AutoDeletePolicy.status.conditions[index]
<sup><sup>[↩ Parent](#autodeletepolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## CachingPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A CachingPolicy sets the hot cache period and optional hot windows.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy caching"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>CachingPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#cachingpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A CachingPolicySpec defines the desired state of a CachingPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#cachingpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A CachingPolicyStatus represents the observed state of a CachingPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.spec
<sup><sup>[↩ Parent](#cachingpolicy)</sup></sup>



A CachingPolicySpec defines the desired state of a CachingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#cachingpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          CachingPolicyParameters are the configurable fields of a CachingPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'MaterializedView', 'Database']: CachingPolicy can only target Table, MaterializedView, Database</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#cachingpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#cachingpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.spec.forProvider
<sup><sup>[↩ Parent](#cachingpolicyspec)</sup></sup>



CachingPolicyParameters are the configurable fields of a CachingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#cachingpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hot</b></td>
        <td>string</td>
        <td>
          Hot is the hot cache span applied to data and index, e.g. "31d".<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#cachingpolicyspecforproviderhotwindowsindex">hotWindows</a></b></td>
        <td>[]object</td>
        <td>
          HotWindows are additional datetime ranges kept in hot cache.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#cachingpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#cachingpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#cachingpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#cachingpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#cachingpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#cachingpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#cachingpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#cachingpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#cachingpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.spec.forProvider.hotWindows[index]
<sup><sup>[↩ Parent](#cachingpolicyspecforprovider)</sup></sup>



HotWindow is a datetime range kept in the hot cache.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>maxValue</b></td>
        <td>string</td>
        <td>
          MaxValue is the exclusive end, ISO 8601.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>minValue</b></td>
        <td>string</td>
        <td>
          MinValue is the inclusive start, ISO 8601 (e.g. 2026-01-01T00:00:00Z).<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### CachingPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#cachingpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### CachingPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#cachingpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### CachingPolicy.status
<sup><sup>[↩ Parent](#cachingpolicy)</sup></sup>



A CachingPolicyStatus represents the observed state of a CachingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#cachingpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#cachingpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.status.atProvider
<sup><sup>[↩ Parent](#cachingpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### CachingPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#cachingpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## EncodingPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A EncodingPolicy sets the encoding policy type of a database, table or single column.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy encoding"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>EncodingPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#encodingpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A EncodingPolicySpec defines the desired state of a EncodingPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#encodingpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A EncodingPolicyStatus represents the observed state of a EncodingPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.spec
<sup><sup>[↩ Parent](#encodingpolicy)</sup></sup>



A EncodingPolicySpec defines the desired state of a EncodingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#encodingpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          EncodingPolicyParameters are the configurable fields of a EncodingPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'Database']: EncodingPolicy can only target Table, Database</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#encodingpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#encodingpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.spec.forProvider
<sup><sup>[↩ Parent](#encodingpolicyspec)</sup></sup>



EncodingPolicyParameters are the configurable fields of a EncodingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#encodingpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type is the encoding policy type, e.g. Identifier, BigObject32 or Vector16.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>column</b></td>
        <td>string</td>
        <td>
          Column applies the policy to one column of the table instead of the whole entity.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#encodingpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#encodingpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#encodingpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#encodingpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#encodingpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#encodingpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#encodingpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#encodingpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#encodingpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#encodingpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### EncodingPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#encodingpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### EncodingPolicy.status
<sup><sup>[↩ Parent](#encodingpolicy)</sup></sup>



A EncodingPolicyStatus represents the observed state of a EncodingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#encodingpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#encodingpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.status.atProvider
<sup><sup>[↩ Parent](#encodingpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### EncodingPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#encodingpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ExtentTagsRetentionPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A ExtentTagsRetentionPolicy removes extent tags with a prefix after a retention period.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy extent_tags_retention"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ExtentTagsRetentionPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A ExtentTagsRetentionPolicySpec defines the desired state of a ExtentTagsRetentionPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A ExtentTagsRetentionPolicyStatus represents the observed state of a ExtentTagsRetentionPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec
<sup><sup>[↩ Parent](#extenttagsretentionpolicy)</sup></sup>



A ExtentTagsRetentionPolicySpec defines the desired state of a ExtentTagsRetentionPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#extenttagsretentionpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          ExtentTagsRetentionPolicyParameters are the configurable fields of a ExtentTagsRetentionPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'Database']: ExtentTagsRetentionPolicy can only target Table, Database</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec.forProvider
<sup><sup>[↩ Parent](#extenttagsretentionpolicyspec)</sup></sup>



ExtentTagsRetentionPolicyParameters are the configurable fields of a ExtentTagsRetentionPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicyspecforproviderrulesindex">rules</a></b></td>
        <td>[]object</td>
        <td>
          Rules is the complete list of tag retention rules.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#extenttagsretentionpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#extenttagsretentionpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#extenttagsretentionpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#extenttagsretentionpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#extenttagsretentionpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec.forProvider.rules[index]
<sup><sup>[↩ Parent](#extenttagsretentionpolicyspecforprovider)</sup></sup>



ExtentTagsRetentionRule drops tags with a prefix after a period.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>retentionPeriod</b></td>
        <td>string</td>
        <td>
          RetentionPeriod after which matching tags are removed.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>tagPrefix</b></td>
        <td>string</td>
        <td>
          TagPrefix of the tags to drop, e.g. "drop-by:".<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#extenttagsretentionpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#extenttagsretentionpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.status
<sup><sup>[↩ Parent](#extenttagsretentionpolicy)</sup></sup>



A ExtentTagsRetentionPolicyStatus represents the observed state of a ExtentTagsRetentionPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#extenttagsretentionpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#extenttagsretentionpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.status.atProvider
<sup><sup>[↩ Parent](#extenttagsretentionpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ExtentTagsRetentionPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#extenttagsretentionpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## IngestionBatchingPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A IngestionBatchingPolicy tunes queued ingestion batching.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy ingestionbatching"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>IngestionBatchingPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A IngestionBatchingPolicySpec defines the desired state of a IngestionBatchingPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A IngestionBatchingPolicyStatus represents the observed state of a IngestionBatchingPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.spec
<sup><sup>[↩ Parent](#ingestionbatchingpolicy)</sup></sup>



A IngestionBatchingPolicySpec defines the desired state of a IngestionBatchingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ingestionbatchingpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          IngestionBatchingPolicyParameters are the configurable fields of a IngestionBatchingPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'Database']: IngestionBatchingPolicy can only target Table, Database</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.spec.forProvider
<sup><sup>[↩ Parent](#ingestionbatchingpolicyspec)</sup></sup>



IngestionBatchingPolicyParameters are the configurable fields of a IngestionBatchingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>maximumBatchingTimeSpan</b></td>
        <td>string</td>
        <td>
          MaximumBatchingTimeSpan is the maximum delay before a batch is sealed (10s to 15m).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maximumNumberOfItems</b></td>
        <td>integer</td>
        <td>
          MaximumNumberOfItems is the maximum number of blobs per batch.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maximumRawDataSizeMB</b></td>
        <td>integer</td>
        <td>
          MaximumRawDataSizeMB is the maximum raw size of a batch in MB.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#ingestionbatchingpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#ingestionbatchingpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#ingestionbatchingpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#ingestionbatchingpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#ingestionbatchingpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#ingestionbatchingpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#ingestionbatchingpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.status
<sup><sup>[↩ Parent](#ingestionbatchingpolicy)</sup></sup>



A IngestionBatchingPolicyStatus represents the observed state of a IngestionBatchingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ingestionbatchingpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestionbatchingpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.status.atProvider
<sup><sup>[↩ Parent](#ingestionbatchingpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionBatchingPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#ingestionbatchingpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## IngestionTimePolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A IngestionTimePolicy adds the hidden ingestion_time() column.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy ingestiontime"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>IngestionTimePolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A IngestionTimePolicySpec defines the desired state of a IngestionTimePolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A IngestionTimePolicyStatus represents the observed state of a IngestionTimePolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.spec
<sup><sup>[↩ Parent](#ingestiontimepolicy)</sup></sup>



A IngestionTimePolicySpec defines the desired state of a IngestionTimePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ingestiontimepolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          IngestionTimePolicyParameters are the configurable fields of a IngestionTimePolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table']: IngestionTimePolicy can only target Table</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.spec.forProvider
<sup><sup>[↩ Parent](#ingestiontimepolicyspec)</sup></sup>



IngestionTimePolicyParameters are the configurable fields of a IngestionTimePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled toggles the policy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#ingestiontimepolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#ingestiontimepolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#ingestiontimepolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#ingestiontimepolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#ingestiontimepolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#ingestiontimepolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#ingestiontimepolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.status
<sup><sup>[↩ Parent](#ingestiontimepolicy)</sup></sup>



A IngestionTimePolicyStatus represents the observed state of a IngestionTimePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ingestiontimepolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ingestiontimepolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.status.atProvider
<sup><sup>[↩ Parent](#ingestiontimepolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### IngestionTimePolicy.status.conditions[index]
<sup><sup>[↩ Parent](#ingestiontimepolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ManagedIdentityPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A ManagedIdentityPolicy allows managed identities to be used for specific usages in a database.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy managed_identity"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ManagedIdentityPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A ManagedIdentityPolicySpec defines the desired state of a ManagedIdentityPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A ManagedIdentityPolicyStatus represents the observed state of a ManagedIdentityPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec
<sup><sup>[↩ Parent](#managedidentitypolicy)</sup></sup>



A ManagedIdentityPolicySpec defines the desired state of a ManagedIdentityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#managedidentitypolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          ManagedIdentityPolicyParameters are the configurable fields of a ManagedIdentityPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Database']: ManagedIdentityPolicy can only target Database</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec.forProvider
<sup><sup>[↩ Parent](#managedidentitypolicyspec)</sup></sup>



ManagedIdentityPolicyParameters are the configurable fields of a ManagedIdentityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicyspecforprovideridentitiesindex">identities</a></b></td>
        <td>[]object</td>
        <td>
          Identities is the complete list of allowed identities.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#managedidentitypolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#managedidentitypolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#managedidentitypolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#managedidentitypolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#managedidentitypolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec.forProvider.identities[index]
<sup><sup>[↩ Parent](#managedidentitypolicyspecforprovider)</sup></sup>



ManagedIdentityEntry allows one identity for usages.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>allowedUsages</b></td>
        <td>[]string</td>
        <td>
          AllowedUsages of the identity.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>objectId</b></td>
        <td>string</td>
        <td>
          ObjectID of the managed identity, or "system".<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#managedidentitypolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#managedidentitypolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.status
<sup><sup>[↩ Parent](#managedidentitypolicy)</sup></sup>



A ManagedIdentityPolicyStatus represents the observed state of a ManagedIdentityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#managedidentitypolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managedidentitypolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.status.atProvider
<sup><sup>[↩ Parent](#managedidentitypolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagedIdentityPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#managedidentitypolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## MergePolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A MergePolicy controls extent merging.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy merge"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>MergePolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#mergepolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A MergePolicySpec defines the desired state of a MergePolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#mergepolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A MergePolicyStatus represents the observed state of a MergePolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.spec
<sup><sup>[↩ Parent](#mergepolicy)</sup></sup>



A MergePolicySpec defines the desired state of a MergePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#mergepolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          MergePolicyParameters are the configurable fields of a MergePolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'MaterializedView', 'Database']: MergePolicy can only target Table, MaterializedView, Database</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#mergepolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#mergepolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.spec.forProvider
<sup><sup>[↩ Parent](#mergepolicyspec)</sup></sup>



MergePolicyParameters are the configurable fields of a MergePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#mergepolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>allowMerge</b></td>
        <td>boolean</td>
        <td>
          AllowMerge enables merge operations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>allowRebuild</b></td>
        <td>boolean</td>
        <td>
          AllowRebuild enables rebuild operations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#mergepolicyspecforproviderlookback">lookback</a></b></td>
        <td>object</td>
        <td>
          Lookback limits which extents are considered.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>loopPeriod</b></td>
        <td>string</td>
        <td>
          LoopPeriod is the maximum time between merge iterations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxExtentsToMerge</b></td>
        <td>integer</td>
        <td>
          MaxExtentsToMerge caps the number of extents per merge.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxRangeInHours</b></td>
        <td>integer</td>
        <td>
          MaxRangeInHours caps the creation time span of merged extents.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>originalSizeMBUpperBoundForMerge</b></td>
        <td>integer</td>
        <td>
          OriginalSizeMBUpperBoundForMerge caps the original size of a merged extent.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>rowCountUpperBoundForMerge</b></td>
        <td>integer</td>
        <td>
          RowCountUpperBoundForMerge caps the row count of a merged extent.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#mergepolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#mergepolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#mergepolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#mergepolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#mergepolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#mergepolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#mergepolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#mergepolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#mergepolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.spec.forProvider.lookback
<sup><sup>[↩ Parent](#mergepolicyspecforprovider)</sup></sup>



Lookback limits which extents are considered.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of lookback.<br/>
          <br/>
            <i>Enum</i>: Default, All, HotCache, Custom<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>customPeriod</b></td>
        <td>string</td>
        <td>
          CustomPeriod is the lookback period for kind Custom.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#mergepolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### MergePolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#mergepolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### MergePolicy.status
<sup><sup>[↩ Parent](#mergepolicy)</sup></sup>



A MergePolicyStatus represents the observed state of a MergePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#mergepolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#mergepolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.status.atProvider
<sup><sup>[↩ Parent](#mergepolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### MergePolicy.status.conditions[index]
<sup><sup>[↩ Parent](#mergepolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## PartitioningPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A PartitioningPolicy partitions extents by hash or uniform range keys.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy partitioning"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>PartitioningPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A PartitioningPolicySpec defines the desired state of a PartitioningPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A PartitioningPolicyStatus represents the observed state of a PartitioningPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec
<sup><sup>[↩ Parent](#partitioningpolicy)</sup></sup>



A PartitioningPolicySpec defines the desired state of a PartitioningPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#partitioningpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          PartitioningPolicyParameters are the configurable fields of a PartitioningPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'MaterializedView']: PartitioningPolicy can only target Table, MaterializedView</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.forProvider
<sup><sup>[↩ Parent](#partitioningpolicyspec)</sup></sup>



PartitioningPolicyParameters are the configurable fields of a PartitioningPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>effectiveDateTime</b></td>
        <td>string</td>
        <td>
          EffectiveDateTime applies the policy to extents created after this ISO 8601 datetime.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspecforproviderpartitionkeysindex">partitionKeys</a></b></td>
        <td>[]object</td>
        <td>
          PartitionKeys are the partition keys (at most one hash and one uniform range key).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#partitioningpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#partitioningpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#partitioningpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#partitioningpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#partitioningpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.forProvider.partitionKeys[index]
<sup><sup>[↩ Parent](#partitioningpolicyspecforprovider)</sup></sup>



PartitionKey is one partition key.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>columnName</b></td>
        <td>string</td>
        <td>
          ColumnName of the key.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the key.<br/>
          <br/>
            <i>Enum</i>: Hash, UniformRange<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicyspecforproviderpartitionkeysindexproperties">properties</a></b></td>
        <td>object</td>
        <td>
          Properties of the key.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.forProvider.partitionKeys[index].properties
<sup><sup>[↩ Parent](#partitioningpolicyspecforproviderpartitionkeysindex)</sup></sup>



Properties of the key.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>function</b></td>
        <td>enum</td>
        <td>
          Function is the hash function (Hash keys).<br/>
          <br/>
            <i>Enum</i>: XxHash64, StringHash<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxPartitionCount</b></td>
        <td>integer</td>
        <td>
          MaxPartitionCount is the number of hash partitions (Hash keys).<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>overrideCreationTime</b></td>
        <td>boolean</td>
        <td>
          OverrideCreationTime sets the extent creation time to the range start.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>partitionAssignmentMode</b></td>
        <td>enum</td>
        <td>
          PartitionAssignmentMode for Hash keys.<br/>
          <br/>
            <i>Enum</i>: Default, Uniform<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>rangeSize</b></td>
        <td>string</td>
        <td>
          RangeSize of a partition (UniformRange keys).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reference</b></td>
        <td>string</td>
        <td>
          Reference datetime (ISO 8601) that anchors the ranges (UniformRange keys).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>seed</b></td>
        <td>integer</td>
        <td>
          Seed of the hash function (Hash keys).<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#partitioningpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### PartitioningPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#partitioningpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### PartitioningPolicy.status
<sup><sup>[↩ Parent](#partitioningpolicy)</sup></sup>



A PartitioningPolicyStatus represents the observed state of a PartitioningPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#partitioningpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#partitioningpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.status.atProvider
<sup><sup>[↩ Parent](#partitioningpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PartitioningPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#partitioningpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## QueryAccelerationPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A QueryAccelerationPolicy accelerates queries over an external delta table by caching recent data (Tier 3).
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy query_acceleration"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>QueryAccelerationPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A QueryAccelerationPolicySpec defines the desired state of a QueryAccelerationPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A QueryAccelerationPolicyStatus represents the observed state of a QueryAccelerationPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.spec
<sup><sup>[↩ Parent](#queryaccelerationpolicy)</sup></sup>



A QueryAccelerationPolicySpec defines the desired state of a QueryAccelerationPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#queryaccelerationpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          QueryAccelerationPolicyParameters are the configurable fields of a QueryAccelerationPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['ExternalTable']: QueryAccelerationPolicy can only target ExternalTable</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.spec.forProvider
<sup><sup>[↩ Parent](#queryaccelerationpolicyspec)</sup></sup>



QueryAccelerationPolicyParameters are the configurable fields of a QueryAccelerationPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled toggles the policy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hot</b></td>
        <td>string</td>
        <td>
          Hot is the period of data kept accelerated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxAge</b></td>
        <td>string</td>
        <td>
          MaxAge caps how stale accelerated data may be.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#queryaccelerationpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#queryaccelerationpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#queryaccelerationpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#queryaccelerationpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#queryaccelerationpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#queryaccelerationpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#queryaccelerationpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.status
<sup><sup>[↩ Parent](#queryaccelerationpolicy)</sup></sup>



A QueryAccelerationPolicyStatus represents the observed state of a QueryAccelerationPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#queryaccelerationpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#queryaccelerationpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.status.atProvider
<sup><sup>[↩ Parent](#queryaccelerationpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### QueryAccelerationPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#queryaccelerationpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## RestrictedViewAccessPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A RestrictedViewAccessPolicy restricts a table to principals with the UnrestrictedViewer role.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy restricted_view_access"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>RestrictedViewAccessPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A RestrictedViewAccessPolicySpec defines the desired state of a RestrictedViewAccessPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A RestrictedViewAccessPolicyStatus represents the observed state of a RestrictedViewAccessPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.spec
<sup><sup>[↩ Parent](#restrictedviewaccesspolicy)</sup></sup>



A RestrictedViewAccessPolicySpec defines the desired state of a RestrictedViewAccessPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#restrictedviewaccesspolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          RestrictedViewAccessPolicyParameters are the configurable fields of a RestrictedViewAccessPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table']: RestrictedViewAccessPolicy can only target Table</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.spec.forProvider
<sup><sup>[↩ Parent](#restrictedviewaccesspolicyspec)</sup></sup>



RestrictedViewAccessPolicyParameters are the configurable fields of a RestrictedViewAccessPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled toggles the policy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#restrictedviewaccesspolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#restrictedviewaccesspolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#restrictedviewaccesspolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#restrictedviewaccesspolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#restrictedviewaccesspolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#restrictedviewaccesspolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#restrictedviewaccesspolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.status
<sup><sup>[↩ Parent](#restrictedviewaccesspolicy)</sup></sup>



A RestrictedViewAccessPolicyStatus represents the observed state of a RestrictedViewAccessPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#restrictedviewaccesspolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#restrictedviewaccesspolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.status.atProvider
<sup><sup>[↩ Parent](#restrictedviewaccesspolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RestrictedViewAccessPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#restrictedviewaccesspolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## RetentionPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A RetentionPolicy controls how long data is kept (soft delete) and whether it is recoverable. Deleting the managed resource removes the policy from the entity, which falls back to the inherited policy.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy retention"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>RetentionPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#retentionpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A RetentionPolicySpec defines the desired state of a RetentionPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#retentionpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A RetentionPolicyStatus represents the observed state of a RetentionPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.spec
<sup><sup>[↩ Parent](#retentionpolicy)</sup></sup>



A RetentionPolicySpec defines the desired state of a RetentionPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#retentionpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          RetentionPolicyParameters are the configurable fields of a RetentionPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'MaterializedView', 'Database']: RetentionPolicy can only target Table, MaterializedView, Database</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#retentionpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#retentionpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.spec.forProvider
<sup><sup>[↩ Parent](#retentionpolicyspec)</sup></sup>



RetentionPolicyParameters are the configurable fields of a RetentionPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#retentionpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>recoverability</b></td>
        <td>enum</td>
        <td>
          Recoverability enables recovery of deleted data for 14 days.<br/>
          <br/>
            <i>Enum</i>: Enabled, Disabled<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>softDeletePeriod</b></td>
        <td>string</td>
        <td>
          SoftDeletePeriod is how long data is kept before soft deletion, e.g. "365d" or "1000000d" for unlimited.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#retentionpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#retentionpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#retentionpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#retentionpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#retentionpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#retentionpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#retentionpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#retentionpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#retentionpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#retentionpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RetentionPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#retentionpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RetentionPolicy.status
<sup><sup>[↩ Parent](#retentionpolicy)</sup></sup>



A RetentionPolicyStatus represents the observed state of a RetentionPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#retentionpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#retentionpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.status.atProvider
<sup><sup>[↩ Parent](#retentionpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RetentionPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#retentionpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## RowLevelSecurityPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A RowLevelSecurityPolicy restricts the rows a principal can read through a filter query.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy row_level_security"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>RowLevelSecurityPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A RowLevelSecurityPolicySpec defines the desired state of a RowLevelSecurityPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A RowLevelSecurityPolicyStatus represents the observed state of a RowLevelSecurityPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.spec
<sup><sup>[↩ Parent](#rowlevelsecuritypolicy)</sup></sup>



A RowLevelSecurityPolicySpec defines the desired state of a RowLevelSecurityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#rowlevelsecuritypolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          RowLevelSecurityPolicyParameters are the configurable fields of a RowLevelSecurityPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'MaterializedView']: RowLevelSecurityPolicy can only target Table, MaterializedView</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.spec.forProvider
<sup><sup>[↩ Parent](#rowlevelsecuritypolicyspec)</sup></sup>



RowLevelSecurityPolicyParameters are the configurable fields of a RowLevelSecurityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>query</b></td>
        <td>string</td>
        <td>
          Query is the KQL filter query (typically a function call) that yields the visible rows.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled toggles the policy. Defaults to true.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#rowlevelsecuritypolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#rowlevelsecuritypolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#rowlevelsecuritypolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#rowlevelsecuritypolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#rowlevelsecuritypolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#rowlevelsecuritypolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#rowlevelsecuritypolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.status
<sup><sup>[↩ Parent](#rowlevelsecuritypolicy)</sup></sup>



A RowLevelSecurityPolicyStatus represents the observed state of a RowLevelSecurityPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#rowlevelsecuritypolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#rowlevelsecuritypolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.status.atProvider
<sup><sup>[↩ Parent](#rowlevelsecuritypolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowLevelSecurityPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#rowlevelsecuritypolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## RowOrderPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A RowOrderPolicy orders rows inside extents by columns (Tier 3, unverified against a cluster).
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy roworder"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>RowOrderPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#roworderpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A RowOrderPolicySpec defines the desired state of a RowOrderPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#roworderpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A RowOrderPolicyStatus represents the observed state of a RowOrderPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec
<sup><sup>[↩ Parent](#roworderpolicy)</sup></sup>



A RowOrderPolicySpec defines the desired state of a RowOrderPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#roworderpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          RowOrderPolicyParameters are the configurable fields of a RowOrderPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'MaterializedView']: RowOrderPolicy can only target Table, MaterializedView</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#roworderpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#roworderpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec.forProvider
<sup><sup>[↩ Parent](#roworderpolicyspec)</sup></sup>



RowOrderPolicyParameters are the configurable fields of a RowOrderPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#roworderpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#roworderpolicyspecforprovidercolumnsindex">columns</a></b></td>
        <td>[]object</td>
        <td>
          Columns and their sort direction.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#roworderpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#roworderpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#roworderpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#roworderpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#roworderpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#roworderpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#roworderpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#roworderpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#roworderpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec.forProvider.columns[index]
<sup><sup>[↩ Parent](#roworderpolicyspecforprovider)</sup></sup>



RowOrderColumn is one sort column.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>direction</b></td>
        <td>enum</td>
        <td>
          Direction of the sort.<br/>
          <br/>
            <i>Enum</i>: asc, desc<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the column.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#roworderpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RowOrderPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#roworderpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### RowOrderPolicy.status
<sup><sup>[↩ Parent](#roworderpolicy)</sup></sup>



A RowOrderPolicyStatus represents the observed state of a RowOrderPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#roworderpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#roworderpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.status.atProvider
<sup><sup>[↩ Parent](#roworderpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RowOrderPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#roworderpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ShardingPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A ShardingPolicy controls extent (shard) sizing.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy sharding"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ShardingPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#shardingpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A ShardingPolicySpec defines the desired state of a ShardingPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#shardingpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A ShardingPolicyStatus represents the observed state of a ShardingPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.spec
<sup><sup>[↩ Parent](#shardingpolicy)</sup></sup>



A ShardingPolicySpec defines the desired state of a ShardingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#shardingpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          ShardingPolicyParameters are the configurable fields of a ShardingPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'MaterializedView', 'Database']: ShardingPolicy can only target Table, MaterializedView, Database</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#shardingpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#shardingpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.spec.forProvider
<sup><sup>[↩ Parent](#shardingpolicyspec)</sup></sup>



ShardingPolicyParameters are the configurable fields of a ShardingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#shardingpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>maxExtentSizeInMb</b></td>
        <td>integer</td>
        <td>
          MaxExtentSizeInMb caps the compressed size per extent.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxOriginalSizeInMb</b></td>
        <td>integer</td>
        <td>
          MaxOriginalSizeInMb caps the original size per extent.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxRowCount</b></td>
        <td>integer</td>
        <td>
          MaxRowCount caps the rows per extent.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>shardEngineMaxExtentSizeInMb</b></td>
        <td>integer</td>
        <td>
          ShardEngineMaxExtentSizeInMb caps the compressed size per shard engine extent.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>shardEngineMaxOriginalSizeInMb</b></td>
        <td>integer</td>
        <td>
          ShardEngineMaxOriginalSizeInMb caps the original size per shard engine extent.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>shardEngineMaxRowCount</b></td>
        <td>integer</td>
        <td>
          ShardEngineMaxRowCount caps rows per extent created by the shard engine.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#shardingpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#shardingpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#shardingpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#shardingpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#shardingpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#shardingpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#shardingpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#shardingpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#shardingpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#shardingpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ShardingPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#shardingpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ShardingPolicy.status
<sup><sup>[↩ Parent](#shardingpolicy)</sup></sup>



A ShardingPolicyStatus represents the observed state of a ShardingPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#shardingpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#shardingpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.status.atProvider
<sup><sup>[↩ Parent](#shardingpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ShardingPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#shardingpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## StreamingIngestionPolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A StreamingIngestionPolicy enables streaming ingestion.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy streamingingestion"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>StreamingIngestionPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A StreamingIngestionPolicySpec defines the desired state of a StreamingIngestionPolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A StreamingIngestionPolicyStatus represents the observed state of a StreamingIngestionPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.spec
<sup><sup>[↩ Parent](#streamingingestionpolicy)</sup></sup>



A StreamingIngestionPolicySpec defines the desired state of a StreamingIngestionPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#streamingingestionpolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          StreamingIngestionPolicyParameters are the configurable fields of a StreamingIngestionPolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table', 'Database']: StreamingIngestionPolicy can only target Table, Database</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.spec.forProvider
<sup><sup>[↩ Parent](#streamingingestionpolicyspec)</sup></sup>



StreamingIngestionPolicyParameters are the configurable fields of a StreamingIngestionPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled toggles streaming ingestion.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hintAllocatedRate</b></td>
        <td>string</td>
        <td>
          HintAllocatedRate hints the expected ingestion rate in GB per hour (decimal number as string).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#streamingingestionpolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#streamingingestionpolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#streamingingestionpolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#streamingingestionpolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#streamingingestionpolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#streamingingestionpolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#streamingingestionpolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.status
<sup><sup>[↩ Parent](#streamingingestionpolicy)</sup></sup>



A StreamingIngestionPolicyStatus represents the observed state of a StreamingIngestionPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#streamingingestionpolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#streamingingestionpolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.status.atProvider
<sup><sup>[↩ Parent](#streamingingestionpolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### StreamingIngestionPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#streamingingestionpolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## UpdatePolicy
<sup><sup>[↩ Parent](#policyadxfunctionalteamv1alpha1 )</sup></sup>






A UpdatePolicy attaches update policies (query-driven ingestion from a source table) to a target table. The list is authoritative for the table.
Only fields set in the spec are compared with the cluster; unset fields keep
the Kusto defaults. Deleting the managed resource deletes the policy on the
entity (".delete ... policy update"), which restores inheritance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>policy.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>UpdatePolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          A UpdatePolicySpec defines the desired state of a UpdatePolicy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#updatepolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          A UpdatePolicyStatus represents the observed state of a UpdatePolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec
<sup><sup>[↩ Parent](#updatepolicy)</sup></sup>



A UpdatePolicySpec defines the desired state of a UpdatePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#updatepolicyspecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          UpdatePolicyParameters are the configurable fields of a UpdatePolicy.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind in ['Table']: UpdatePolicy can only target Table</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider
<sup><sup>[↩ Parent](#updatepolicyspec)</sup></sup>



UpdatePolicyParameters are the configurable fields of a UpdatePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderupdatesindex">updates</a></b></td>
        <td>[]object</td>
        <td>
          Updates is the complete list of update policies of the table.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.entity
<sup><sup>[↩ Parent](#updatepolicyspecforprovider)</sup></sup>



Entity the policy applies to. Which kinds are allowed depends on the
policy type and is validated per kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#updatepolicyspecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#updatepolicyspecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#updatepolicyspecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#updatepolicyspecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.updates[index]
<sup><sup>[↩ Parent](#updatepolicyspecforprovider)</sup></sup>



UpdatePolicyEntry is one update policy of a table.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>query</b></td>
        <td>string</td>
        <td>
          Query is the KQL query (or function call) producing rows for the target table.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>source</b></td>
        <td>string</td>
        <td>
          Source is the table whose ingestion triggers the query.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled toggles the policy. Defaults to true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>managedIdentity</b></td>
        <td>string</td>
        <td>
          ManagedIdentity ("system" or an object id) runs the query with a managed identity.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>propagateIngestionProperties</b></td>
        <td>boolean</td>
        <td>
          PropagateIngestionProperties copies extent tags and creation time.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderupdatesindexsourceref">sourceRef</a></b></td>
        <td>object</td>
        <td>
          SourceRef references a Table managed resource as source.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderupdatesindexsourceselector">sourceSelector</a></b></td>
        <td>object</td>
        <td>
          SourceSelector selects a Table managed resource as source.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>transactional</b></td>
        <td>boolean</td>
        <td>
          Transactional fails the source ingestion when the update fails.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.updates[index].sourceRef
<sup><sup>[↩ Parent](#updatepolicyspecforproviderupdatesindex)</sup></sup>



SourceRef references a Table managed resource as source.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderupdatesindexsourcerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.updates[index].sourceRef.policy
<sup><sup>[↩ Parent](#updatepolicyspecforproviderupdatesindexsourceref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.updates[index].sourceSelector
<sup><sup>[↩ Parent](#updatepolicyspecforproviderupdatesindex)</sup></sup>



SourceSelector selects a Table managed resource as source.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicyspecforproviderupdatesindexsourceselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.forProvider.updates[index].sourceSelector.policy
<sup><sup>[↩ Parent](#updatepolicyspecforproviderupdatesindexsourceselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.providerConfigRef
<sup><sup>[↩ Parent](#updatepolicyspec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### UpdatePolicy.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#updatepolicyspec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### UpdatePolicy.status
<sup><sup>[↩ Parent](#updatepolicy)</sup></sup>



A UpdatePolicyStatus represents the observed state of a UpdatePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#updatepolicystatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          PolicyObservation is the observed state shared by all policy kinds.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#updatepolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.status.atProvider
<sup><sup>[↩ Parent](#updatepolicystatus)</sup></sup>



PolicyObservation is the observed state shared by all policy kinds.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entity</b></td>
        <td>string</td>
        <td>
          Entity the policy was read from, in Kusto notation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          Policy is the policy JSON as reported by the cluster for this entity
(not the effective/inherited policy).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### UpdatePolicy.status.conditions[index]
<sup><sup>[↩ Parent](#updatepolicystatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

# security.adx.functional.team/v1alpha1

Resource Types:

- [SecurityRole](#securityrole)




## SecurityRole
<sup><sup>[↩ Parent](#securityadxfunctionalteamv1alpha1 )</sup></sup>






A SecurityRole assigns principals to a role of a database, table, external
table, materialized view or function. Deleting an Authoritative resource
removes all principals of that role; deleting an Additive one removes only
the principals it added.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>security.adx.functional.team/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>SecurityRole</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#securityrolespec">spec</a></b></td>
        <td>object</td>
        <td>
          A SecurityRoleSpec defines the desired state of a SecurityRole.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#securityrolestatus">status</a></b></td>
        <td>object</td>
        <td>
          A SecurityRoleStatus represents the observed state of a SecurityRole.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.spec
<sup><sup>[↩ Parent](#securityrole)</sup></sup>



A SecurityRoleSpec defines the desired state of a SecurityRole.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#securityrolespecforprovider">forProvider</a></b></td>
        <td>object</td>
        <td>
          SecurityRoleParameters are the configurable fields of a SecurityRole.<br/>
          <br/>
            <i>Validations</i>:<li>self.entity.kind != 'Table' || self.role in ['admins','ingestors']: tables only have the admins and ingestors roles</li><li>!(self.entity.kind in ['ExternalTable','MaterializedView','Function']) || self.role == 'admins': external tables, materialized views and functions only have the admins role</li><li>self.entity.kind == oldSelf.entity.kind: entity.kind is immutable</li><li>!has(oldSelf.entity.name) || !has(self.entity.name) || self.entity.name == oldSelf.entity.name: entity.name is immutable once set</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>managementPolicies</b></td>
        <td>[]enum</td>
        <td>
          THIS IS A BETA FIELD. It is on by default but can be opted out
through a Crossplane feature flag.
ManagementPolicies specify the array of actions Crossplane is allowed to
take on the managed and external resources.
See the design doc for more information: https://github.com/crossplane/crossplane/blob/499895a25d1a1a0ba1604944ef98ac7a1a71f197/design/design-doc-observe-only-resources.md?plain=1#L223
and this one: https://github.com/crossplane/crossplane/blob/444267e84783136daa93568b364a5f01228cacbe/design/one-pager-ignore-changes.md<br/>
          <br/>
            <i>Enum</i>: Observe, Create, Update, Delete, LateInitialize, *<br/>
            <i>Default</i>: [*]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#securityrolespecproviderconfigref">providerConfigRef</a></b></td>
        <td>object</td>
        <td>
          ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.<br/>
          <br/>
            <i>Default</i>: map[kind:ClusterProviderConfig name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#securityrolespecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.spec.forProvider
<sup><sup>[↩ Parent](#securityrolespec)</sup></sup>



SecurityRoleParameters are the configurable fields of a SecurityRole.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Database that holds the entity. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: database is immutable</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#securityrolespecforproviderentity">entity</a></b></td>
        <td>object</td>
        <td>
          Entity the role belongs to: Database, Table, ExternalTable,
MaterializedView or Function.<br/>
          <br/>
            <i>Validations</i>:<li>self.kind == 'Database' || has(self.name) || has(self.nameRef) || has(self.nameSelector): name, nameRef or nameSelector is required unless kind is Database</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>enum</td>
        <td>
          Role to manage. Immutable.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: role is immutable</li>
            <i>Enum</i>: admins, users, viewers, unrestrictedviewers, ingestors, monitors<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Description recorded with the role assignment.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>enum</td>
        <td>
          Mode is Authoritative (default; the spec is the complete principal
list of the role) or Additive (only the listed principals are managed).<br/>
          <br/>
            <i>Enum</i>: Authoritative, Additive<br/>
            <i>Default</i>: Authoritative<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>principals</b></td>
        <td>[]string</td>
        <td>
          Principals in Kusto notation, e.g. "aaduser=alice@contoso.com",
"aadapp=<appId>;<tenantId>", "aadgroup=<objectId>;<tenantId>". An empty
list in Authoritative mode removes every principal from the role.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.spec.forProvider.entity
<sup><sup>[↩ Parent](#securityrolespecforprovider)</sup></sup>



Entity the role belongs to: Database, Table, ExternalTable,
MaterializedView or Function.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the target entity.<br/>
          <br/>
            <i>Enum</i>: Database, Table, MaterializedView, ExternalTable, Function<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the target entity (Kusto name). Ignored for kind Database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#securityrolespecforproviderentitynameref">nameRef</a></b></td>
        <td>object</td>
        <td>
          NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#securityrolespecforproviderentitynameselector">nameSelector</a></b></td>
        <td>object</td>
        <td>
          NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.spec.forProvider.entity.nameRef
<sup><sup>[↩ Parent](#securityrolespecforproviderentity)</sup></sup>



NameRef references a managed resource whose external name is used as
the target. The referenced kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#securityrolespecforproviderentitynamerefpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for referencing.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.spec.forProvider.entity.nameRef.policy
<sup><sup>[↩ Parent](#securityrolespecforproviderentitynameref)</sup></sup>



Policies for referencing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.spec.forProvider.entity.nameSelector
<sup><sup>[↩ Parent](#securityrolespecforproviderentity)</sup></sup>



NameSelector selects a managed resource whose external name is used as
the target. The selected kind follows from kind.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchControllerRef</b></td>
        <td>boolean</td>
        <td>
          MatchControllerRef ensures an object with the same controller reference
as the selecting object is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          MatchLabels ensures an object with matching labels is selected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace for the selector<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#securityrolespecforproviderentitynameselectorpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Policies for selection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.spec.forProvider.entity.nameSelector.policy
<sup><sup>[↩ Parent](#securityrolespecforproviderentitynameselector)</sup></sup>



Policies for selection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resolution</b></td>
        <td>enum</td>
        <td>
          Resolution specifies whether resolution of this reference is required.
The default is 'Required', which means the reconcile will fail if the
reference cannot be resolved. 'Optional' means this reference will be
a no-op if it cannot be resolved.<br/>
          <br/>
            <i>Enum</i>: Required, Optional<br/>
            <i>Default</i>: Required<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolve</b></td>
        <td>enum</td>
        <td>
          Resolve specifies when this reference should be resolved. The default
is 'IfNotPresent', which will attempt to resolve the reference only when
the corresponding field is not present. Use 'Always' to resolve the
reference on every reconcile.<br/>
          <br/>
            <i>Enum</i>: Always, IfNotPresent<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.spec.providerConfigRef
<sup><sup>[↩ Parent](#securityrolespec)</sup></sup>



ProviderConfigReference specifies how the provider that will be used to
create, observe, update, and delete this managed resource should be
configured.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### SecurityRole.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#securityrolespec)</sup></sup>



WriteConnectionSecretToReference specifies the namespace and name of a
Secret to which any connection details for this managed resource should
be written. Connection details frequently include the endpoint, username,
and password required to connect to the managed resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### SecurityRole.status
<sup><sup>[↩ Parent](#securityrole)</sup></sup>



A SecurityRoleStatus represents the observed state of a SecurityRole.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#securityrolestatusatprovider">atProvider</a></b></td>
        <td>object</td>
        <td>
          SecurityRoleObservation are the observable fields of a SecurityRole.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#securityrolestatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastHandledReconcileAt</b></td>
        <td>string</td>
        <td>
          LastHandledReconcileAt holds the value of the most recent
reconcile-requested-at annotation token that the controller has
processed. Users can compare this to the annotation to determine
whether a reconcile request has been handled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the latest metadata.generation
which resulted in either a ready state, or stalled due to error
it can not recover from without human intervention.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.status.atProvider
<sup><sup>[↩ Parent](#securityrolestatus)</sup></sup>



SecurityRoleObservation are the observable fields of a SecurityRole.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>principals</b></td>
        <td>[]string</td>
        <td>
          Principals are the FQNs currently assigned to the role in the cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#securityrolestatusatproviderresolvedprincipalsindex">resolvedPrincipals</a></b></td>
        <td>[]object</td>
        <td>
          ResolvedPrincipals maps spec entries to the object ids Kusto reports.
The provider compares sets of object ids, so the notation used in the
spec (UPN, app id, object id) does not matter.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>unresolved</b></td>
        <td>[]string</td>
        <td>
          Unresolved lists spec entries that could not be mapped to a principal
reported by the cluster (see the provider documentation).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.status.atProvider.resolvedPrincipals[index]
<sup><sup>[↩ Parent](#securityrolestatusatprovider)</sup></sup>



ResolvedPrincipal maps one spec entry to the principal Kusto reports.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>objectId</b></td>
        <td>string</td>
        <td>
          ObjectID is the Entra object id Kusto reports for it.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>spec</b></td>
        <td>string</td>
        <td>
          Spec is the principal string from the spec.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>displayName</b></td>
        <td>string</td>
        <td>
          DisplayName Kusto reports.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>fqn</b></td>
        <td>string</td>
        <td>
          FQN is the fully qualified principal name Kusto reports.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of the principal as reported (e.g. AAD User).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SecurityRole.status.conditions[index]
<sup><sup>[↩ Parent](#securityrolestatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
