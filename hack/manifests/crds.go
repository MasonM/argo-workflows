package main

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

func cleanCRD(filename string) {
	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		panic(err)
	}
	crd := make(obj)
	err = yaml.Unmarshal(data, &crd)
	if err != nil {
		panic(err)
	}
	crd.RemoveNestedField("status")
	crd.RemoveNestedField("metadata", "annotations")
	crd.RemoveNestedField("metadata", "creationTimestamp")
	schema := crd.OpenAPIV3Schema()
	switch crd.Name() {
	case "cronworkflows.argoproj.io":
		patchWorkflowSpecTemplateFields(&schema, "spec", "properties", "workflowSpec", "properties")
	case "clusterworkflowtemplates.argoproj.io", "workflows.argoproj.io", "workflowtemplates.argoproj.io":
		patchWorkflowSpecTemplateFields(&schema, "spec", "properties")
	}
	if crd.Name() == "workflows.argoproj.io" {
		patchTemplateFields(&schema, "status", "properties", "storedTemplates", "additionalProperties")
		patchWorkflowSpecTemplateFields(&schema, "status", "properties", "storedWorkflowTemplateSpec", "properties")
	}
	data, err = yaml.Marshal(crd)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filename, data, 0o600)
	if err != nil {
		panic(err)
	}
}

func patchWorkflowSpecTemplateFields(schema *obj, baseFields ...string) {
	patchTemplateFields(schema, append(baseFields, "templateDefaults")...)
	patchTemplateFields(schema, append(baseFields, "templates", "items")...)
}

func patchTemplateFields(schema *obj, baseFields ...string) {
	schema.SetNestedField([]string{"image"}, append(baseFields, "properties", "container", "required")...)
	schema.SetNestedField([]string{"image", "source"}, append(baseFields, "properties", "script", "required")...)
	stepFields := append(baseFields, "properties", "steps")
	schema.CopyNestedField(append(stepFields, "items", "properties", "steps"), stepFields)
}

// minimizeCRD generates a stripped-down CRD as a workaround for "Request entity too large: limit is 3145728" errors due to https://github.com/kubernetes/kubernetes/issues/82292.
func minimizeCRD(filename string) {
	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		panic(err)
	}

	shouldMinimize := false
	if len(data) > 512*1024 {
		fmt.Printf("Minimizing %s due to CRD size (%d) exceeding 512KB\n", filename, len(data))
		shouldMinimize = true
	}

	crd := make(obj)
	err = yaml.Unmarshal(data, &crd)
	if err != nil {
		panic(err)
	}

	if !shouldMinimize {
		return
	}

	stripSpecAndStatusFields(&crd)

	data, err = yaml.Marshal(crd)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filename, data, 0o600)
	if err != nil {
		panic(err)
	}
}

// stripSpecAndStatusFields strips the "spec" and "status" fields from the CRD, as those are usually the largest.
func stripSpecAndStatusFields(crd *obj) {
	schema := crd.OpenAPIV3Schema()
	preserveMarker := obj{"type": "object", "x-kubernetes-preserve-unknown-fields": true, "x-kubernetes-map-type": "atomic"}
	if _, ok := schema["spec"]; ok {
		schema["spec"] = preserveMarker
	}
	if _, ok := schema["status"]; ok {
		schema["status"] = preserveMarker
	}
}
