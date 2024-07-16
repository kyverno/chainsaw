package main

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-tools/pkg/crd/markers"
)

type MergePatch string

func (m MergePatch) ApplyToSchema(schema *apiext.JSONSchemaProps) error {
	in, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	patch := []byte(m)
	patched, err := jsonpatch.MergeMergePatches(in, patch)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(patched, schema); err != nil {
		return err
	}
	return nil
}

func (m MergePatch) ApplyPriority() markers.ApplyPriority {
	return markers.ApplyPriorityDefault * 2
}

type OneOf []string

func (m OneOf) ApplyToSchema(schema *apiext.JSONSchemaProps) error {
	for _, prop := range m {
		schema.OneOf = append(schema.OneOf, apiext.JSONSchemaProps{
			Required: []string{prop},
		})
	}
	return nil
}
