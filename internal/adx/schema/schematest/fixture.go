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

// Package schematest holds recorded ".show database schema as json" output
// for tests of the kinds that read the schema section.
package schematest

// DatabaseJSON is a schema JSON in the shape returned by the emulator and the
// service (tables, functions, materialized views, external tables).
const DatabaseJSON = `{"Databases":{"Telemetry":{"Name":"Telemetry","Tables":{"RawEvents":{"Name":"RawEvents","OrderedColumns":[{"Name":"Timestamp","Type":"System.DateTime","CslType":"datetime"},{"Name":"Payload","Type":"System.Object","CslType":"dynamic","DocString":"Raw JSON body"}],"Folder":"Raw","DocString":"Landing table"},"Dim":{"Name":"Dim","OrderedColumns":[{"Name":"Id","Type":"System.String","CslType":"string"}]}},"ExternalTables":{"Exports":{"Name":"Exports","OrderedColumns":[{"Name":"A","Type":"System.String","CslType":"string"}],"Folder":"External","DocString":""}},"MaterializedViews":{"LatestEvents":{"Name":"LatestEvents","OrderedColumns":[{"Name":"DeviceId","Type":"System.String","CslType":"string"}],"Query":"RawEvents | summarize arg_max(Timestamp, *) by DeviceId","SourceTable":"RawEvents","Folder":"Views","DocString":""}},"Functions":{"TakeSome":{"Name":"TakeSome","InputParameters":[{"Name":"limit","CslType":"long","CslDefaultValue":"100"},{"Name":"T","Columns":[{"Name":"x","CslType":"long"}]}],"Body":"{ RawEvents | take limit }","Folder":"Parsing","DocString":"doc","FunctionKind":"UnknownFunction","OutputColumns":[]}},"EntityGroups":{"EG":["cluster('c').database('d')"]},"DatabaseAccessMode":"ReadWrite","MajorVersion":12,"MinorVersion":3,"PrettyName":"Telemetry"}}}`
