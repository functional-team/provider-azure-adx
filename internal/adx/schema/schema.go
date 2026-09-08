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

// Package schema loads and parses ".show database schema as json", the one
// command that describes all tables, functions, materialized views, external
// tables and entity groups of a database. It backs the batch-observe cache
// section shared by those kinds.
package schema

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/functional-team/provider-azure-adx/internal/clients/kusto"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/cmd"
	"github.com/functional-team/provider-azure-adx/internal/clients/kusto/snapshot"
)

// Section is the snapshot cache section name.
const Section = "schema"

// Column is a column in the schema JSON.
type Column struct {
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	CslType   string `json:"CslType"`
	DocString string `json:"DocString"`
}

// Table is a table in the schema JSON.
type Table struct {
	Name           string   `json:"Name"`
	OrderedColumns []Column `json:"OrderedColumns"`
	Folder         string   `json:"Folder"`
	DocString      string   `json:"DocString"`
}

// FunctionParameter is one input parameter of a function; tabular parameters
// carry Columns instead of CslType.
type FunctionParameter struct {
	Name            string   `json:"Name"`
	CslType         string   `json:"CslType"`
	CslDefaultValue string   `json:"CslDefaultValue"`
	Columns         []Column `json:"Columns"`
}

// Function is a stored function in the schema JSON.
type Function struct {
	Name            string              `json:"Name"`
	InputParameters []FunctionParameter `json:"InputParameters"`
	Body            string              `json:"Body"`
	Folder          string              `json:"Folder"`
	DocString       string              `json:"DocString"`
	FunctionKind    string              `json:"FunctionKind"`
	OutputColumns   []Column            `json:"OutputColumns"`
}

// MaterializedView is a materialized view in the schema JSON.
type MaterializedView struct {
	Name           string   `json:"Name"`
	OrderedColumns []Column `json:"OrderedColumns"`
	Query          string   `json:"Query"`
	SourceTable    string   `json:"SourceTable"`
	Folder         string   `json:"Folder"`
	DocString      string   `json:"DocString"`
}

// ExternalTable is an external table in the schema JSON.
type ExternalTable struct {
	Name           string   `json:"Name"`
	OrderedColumns []Column `json:"OrderedColumns"`
	Folder         string   `json:"Folder"`
	DocString      string   `json:"DocString"`
}

// EntityGroup is an entity group in the schema JSON.
type EntityGroup struct {
	Name     string   `json:"Name"`
	Entities []string `json:"Entities"`
}

// Database is the parsed schema of one database.
type Database struct {
	Name              string                      `json:"Name"`
	Tables            map[string]Table            `json:"Tables"`
	MaterializedViews map[string]MaterializedView `json:"MaterializedViews"`
	Functions         map[string]Function         `json:"Functions"`
	ExternalTables    map[string]ExternalTable    `json:"ExternalTables"`
	EntityGroups      map[string]EntityGroup      `json:"EntityGroups"`
}

type envelope struct {
	Databases map[string]*Database `json:"Databases"`
}

// Parse parses the DatabaseSchema JSON document. With several databases the
// one named db is returned; with exactly one, that one.
func Parse(raw []byte, db string) (*Database, error) {
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("cannot parse database schema JSON: %w", err)
	}
	if d, ok := env.Databases[db]; ok && d != nil {
		return normalize(d), nil
	}
	if len(env.Databases) == 1 {
		for _, d := range env.Databases {
			if d != nil {
				return normalize(d), nil
			}
		}
	}
	return nil, fmt.Errorf("database %q not found in schema JSON", db)
}

func normalize(d *Database) *Database {
	if d.Tables == nil {
		d.Tables = map[string]Table{}
	}
	if d.MaterializedViews == nil {
		d.MaterializedViews = map[string]MaterializedView{}
	}
	if d.Functions == nil {
		d.Functions = map[string]Function{}
	}
	if d.ExternalTables == nil {
		d.ExternalTables = map[string]ExternalTable{}
	}
	if d.EntityGroups == nil {
		d.EntityGroups = map[string]EntityGroup{}
	}
	return d
}

// ShowCmd is ".show database ['DB'] schema as json".
func ShowCmd(db string) cmd.Command {
	return cmd.New(".show database ", cmd.Ident(db), " schema as json")
}

// Load runs ShowCmd and parses the result.
func Load(ctx context.Context, kc kusto.Client, db string) (*Database, error) {
	res, err := kc.Mgmt(ctx, db, ShowCmd(db))
	if err != nil {
		return nil, err
	}
	rows := res.Rows()
	if len(rows) == 0 {
		return nil, fmt.Errorf("schema command returned no rows for database %q", db)
	}
	raw := rows[0].String("DatabaseSchema")
	if raw == "" && len(rows[0].Dynamic("DatabaseSchema")) > 0 {
		raw = string(rows[0].Dynamic("DatabaseSchema"))
	}
	if raw == "" {
		return nil, fmt.Errorf("schema command returned an empty DatabaseSchema for %q", db)
	}
	return Parse([]byte(raw), db)
}

// Key is the cache key of the schema section for db.
func Key(endpoint, db string) snapshot.Key {
	return snapshot.Key{Endpoint: endpoint, Database: db, Section: Section}
}

// Cached returns the schema of db through the snapshot cache.
func Cached(ctx context.Context, c *snapshot.Cache, kc kusto.Client, db string) (*Database, error) {
	return snapshot.Load(ctx, c, Key(kc.Endpoint(), db), func(ctx context.Context) (*Database, error) {
		return Load(ctx, kc, db)
	})
}

// Invalidate drops the cached schema of db.
func Invalidate(c *snapshot.Cache, kc kusto.Client, db string) {
	c.Invalidate(kc.Endpoint(), db, Section)
}
