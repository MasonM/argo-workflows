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
	parentField := NestedFieldNoCopy[map[string]interface{}](o, fields[:len(fields)-1]...)
	parentField[fields[len(fields)-1]] = value
}

func (o *obj) CopyNestedField(sourceFields []string, targetFields []string) {
	value := NestedFieldNoCopy[any](o, sourceFields...)
	parentField := NestedFieldNoCopy[map[string]interface{}](o, targetFields[:len(targetFields)-1]...)
	parentField[targetFields[len(targetFields)-1]] = value
}

func (o *obj) Name() string {
	return NestedFieldNoCopy[string](o, "metadata", "name")
}

func (o *obj) OpenAPIV3Schema() obj {
	versions := NestedFieldNoCopy[[]interface{}](o, "spec", "versions")
	version := obj(versions[0].(map[string]interface{}))
	return NestedFieldNoCopy[map[string]interface{}](&version, "schema", "openAPIV3Schema", "properties")
}

func NestedFieldNoCopy[T any](o *obj, fields ...string) T {
	value, found, err := unstructured.NestedFieldNoCopy(*o, fields...)
	if !found {
		panic(fmt.Sprintf("failed to find field %v", fields))
	}
	if err != nil {
		panic(err.Error())
	}
	return value.(T)
}
