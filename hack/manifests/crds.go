package main

import (
	"fmt"
	"os"
	"path/filepath"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func cleanCRD(filename string) {
	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		panic(err)
	}
	crd := apiext.CustomResourceDefinition{}
	err = yaml.Unmarshal(data, &crd)
	if err != nil {
		panic(err)
	}
	crd.Status = apiext.CustomResourceDefinitionStatus{}
	crd.Annotations = nil
	crd.CreationTimestamp = metav1.Time{}
	schema := crd.Spec.Versions[0].Schema.OpenAPIV3Schema
	switch crd.Name {
	case "cronworkflows.argoproj.io":
		specProperties := schema.Properties["Spec"].Properties["workflowSpec"].Properties
		patchWorkflowSpecTemplateFields(&specProperties)
	case "clusterworkflowtemplates.argoproj.io", "workflows.argoproj.io", "workflowtemplates.argoproj.io":
		specProperties := schema.Properties["Spec"].Properties
		patchWorkflowSpecTemplateFields(&specProperties)
	}
	if crd.Name == "workflows.argoproj.io" {
		statusProperties := schema.Properties["status"].Properties
		storedTemplates := statusProperties["storedTemplates"].AdditionalProperties.Schema
		patchTemplateFields(storedTemplates)
		storedWorkflowTemplateSpec := statusProperties["storedWorkflowTemplateSpec"].Properties
		patchWorkflowSpecTemplateFields(&storedWorkflowTemplateSpec)
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

func patchWorkflowSpecTemplateFields(specProperties *map[string]apiext.JSONSchemaProps) {
	for _, properties := range []apiext.JSONSchemaProps{(*specProperties)["templateDefaults"], *(*specProperties)["template"].Items.Schema} {
		patchTemplateFields(&properties)
	}
}

func patchTemplateFields(field *apiext.JSONSchemaProps) {
	properties := (*field).Properties
	container := properties["container"]
	container.Required = []string{"image"}
	script := properties["script"]
	script.Required = []string{"image", "source"}
	steps := properties["steps"]
	nestedSteps := steps.Items.Schema.Properties["steps"]
	steps.Items = &apiext.JSONSchemaPropsOrArray{Schema: &nestedSteps}
}

// minimizeCRD generates a stripped-down CRD as a workaround for "Request entity too large: limit is 3145728" errors due to https://github.com/kubernetes/kubernetes/issues/82292.
func minimizeCRD(filename string) {
	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		panic(err)
	}

	shouldMinimize := false
	if len(data) > 1024*1024 {
		fmt.Printf("Minimizing %s due to CRD size (%d) exceeding 1MB\n", filename, len(data))
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

	crd = stripSpecAndStatusFields(crd)

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
func stripSpecAndStatusFields(crd obj) obj {
	spec := crd["spec"].(obj)
	versions := spec["versions"].([]interface{})
	version := versions[0].(obj)
	properties := version["schema"].(obj)["openAPIV3Schema"].(obj)["properties"].(obj)
	for k := range properties {
		if k == "spec" || k == "status" {
			properties[k] = obj{"type": "object", "x-kubernetes-preserve-unknown-fields": true, "x-kubernetes-map-type": "atomic"}
		}
	}
	return crd
}
