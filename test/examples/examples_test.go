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

// Package examples validates every manifest under examples/ against the
// generated CRD schemas, so a typo in an example or a field renamed in the API
// fails the unit tests instead of the first kubectl apply.
package examples

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/validation"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

type crdVersion struct {
	validator validation.SchemaValidator
	scope     apiextensionsv1.ResourceScope
}

func loadCRDs(t *testing.T) map[schema.GroupVersionKind]crdVersion {
	t.Helper()
	files, err := filepath.Glob("../../package/crds/*.yaml")
	if err != nil || len(files) == 0 {
		t.Fatalf("no CRDs found (run go generate): %v", err)
	}
	out := map[schema.GroupVersionKind]crdVersion{}
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var crd apiextensionsv1.CustomResourceDefinition
		if err := yaml.Unmarshal(raw, &crd); err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		for _, v := range crd.Spec.Versions {
			if v.Schema == nil || v.Schema.OpenAPIV3Schema == nil {
				continue
			}
			internal := &apiextensions.JSONSchemaProps{}
			if err := apiextensionsv1.Convert_v1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(v.Schema.OpenAPIV3Schema, internal, nil); err != nil {
				t.Fatalf("%s: %v", f, err)
			}
			val, _, err := validation.NewSchemaValidator(internal)
			if err != nil {
				t.Fatalf("%s: %v", f, err)
			}
			out[schema.GroupVersionKind{Group: crd.Spec.Group, Version: v.Name, Kind: crd.Spec.Names.Kind}] = crdVersion{validator: val, scope: crd.Spec.Scope}
		}
	}
	return out
}

func TestExamplesMatchCRDs(t *testing.T) {
	crds := loadCRDs(t)
	var files []string
	if err := filepath.WalkDir("../../examples", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no examples found")
	}
	covered := map[schema.GroupVersionKind]bool{}
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		dec := utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), 4096)
		for i := 0; ; i++ {
			var obj unstructured.Unstructured
			if err := dec.Decode(&obj); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				t.Fatalf("%s: document %d: %v", f, i, err)
			}
			if len(obj.Object) == 0 {
				continue
			}
			gvk := obj.GroupVersionKind()
			if !strings.HasSuffix(gvk.Group, "adx.functional.team") {
				continue // Secrets, Namespaces, XRDs, Compositions: not ours
			}
			crd, ok := crds[gvk]
			if !ok {
				t.Errorf("%s: no CRD for %s", f, gvk)
				continue
			}
			covered[gvk] = true
			if crd.scope == apiextensionsv1.NamespaceScoped && obj.GetNamespace() == "" {
				t.Errorf("%s: %s %s is namespaced but has no metadata.namespace", f, gvk.Kind, obj.GetName())
			}
			if errs := validation.ValidateCustomResource(nil, obj.Object, crd.validator); len(errs) > 0 {
				t.Errorf("%s: %s %s does not validate against its CRD:\n%v", f, gvk.Kind, obj.GetName(), errs.ToAggregate())
			}
		}
	}
	// Every managed resource kind needs an example (docs/tech-implement.md 13).
	for gvk := range crds {
		if strings.HasPrefix(gvk.Kind, "ProviderConfig") || gvk.Kind == "ClusterProviderConfig" {
			continue
		}
		if !covered[gvk] {
			t.Errorf("kind %s has no example under examples/", gvk)
		}
	}
}
