package main

import "fmt"

func cleanCRD(filename string) {
	crd := ParseYaml(Read(filename))
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
	crd.WriteYaml(filename)
}

func patchWorkflowSpecTemplateFields(schema *obj, baseFields ...string) {
	patchTemplateFields(schema, append(baseFields, "templateDefaults")...)
	patchTemplateFields(schema, append(baseFields, "templates", "items")...)
}

func patchTemplateFields(schema *obj, baseFields ...string) {
	// container and script templates embed the k8s.io/api/core/v1/Container
	// struct, and kubebuilder marks the "name" field as required, but it's not actually required.
	schema.RemoveNestedField(append(baseFields, "properties", "container", "required")...)
	schema.RemoveNestedField(append(baseFields, "properties", "script", "required")...)
	stepFields := append(baseFields, "properties", "steps", "items")
	schema.CopyNestedField(append(stepFields, "properties", "steps"), stepFields)
}

// minimizeCRD generates a stripped-down CRD as a workaround for "Request entity too large: limit is 3145728" errors due to https://github.com/kubernetes/kubernetes/issues/82292.
func minimizeCRD(filename string) {
	data := Read(filename)
	shouldMinimize := false
	if len(data) > 512*1024 {
		fmt.Printf("Minimizing %s due to CRD size (%d) exceeding 512KB\n", filename, len(data))
		shouldMinimize = true
	}
	if !shouldMinimize {
		return
	}
	crd := ParseYaml(data)
	stripSpecAndStatusFields(crd)
	crd.WriteYaml(filename)
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
