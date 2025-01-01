package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type obj map[string]interface{}

func (o *obj) RemoveNestedField(fields ...string) {
	unstructured.RemoveNestedField(*o, fields...)
}

func (o *obj) SetNestedField(value interface{}, fields ...string) {
	parentField := o.NestedFieldNoCopy(fields[:len(fields)-1]...)
	parentField.(map[string]interface{})[fields[len(fields)-1]] = value
}

func (o *obj) NestedFieldNoCopy(fields ...string) interface{} {
	value, found, err := unstructured.NestedFieldNoCopy(*o, fields...)
	if !found {
		panic(fmt.Sprintf("failed to find field %v", fields))
	}
	if err != nil {
		panic(err.Error())
	}
	return value
}

func (o *obj) Name() string {
	return o.NestedFieldNoCopy("metadata", "name").(string)
}

func (o *obj) OpenAPIV3Schema() obj {
	versions := o.NestedFieldNoCopy("spec", "versions")
	version := obj(versions.([]interface{})[0].(map[string]interface{}))
	return version.NestedFieldNoCopy("schema", "openAPIV3Schema").(map[string]interface{})
}
