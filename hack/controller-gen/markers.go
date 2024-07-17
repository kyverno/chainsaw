package main

import (
	"encoding/json"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type OneOf []any

func (m OneOf) ApplyToSchema(schema *apiext.JSONSchemaProps) error {
	for _, prop := range m {
		if sprop, ok := prop.(string); ok {
			schema.OneOf = append(schema.OneOf, apiext.JSONSchemaProps{
				Required: []string{sprop},
			})
		} else {
			var props apiext.JSONSchemaProps
			if data, err := json.Marshal(prop); err != nil {
				return err
			} else if err := json.Unmarshal(data, &props); err != nil {
				return err
			}
			schema.OneOf = append(schema.OneOf, props)
		}
	}
	return nil
}

type Not struct {
	Value interface{}
}

func (m Not) ApplyToSchema(schema *apiext.JSONSchemaProps) error {
	var props apiext.JSONSchemaProps
	if data, err := json.Marshal(m.Value); err != nil {
		return err
	} else if err := json.Unmarshal(data, &props); err != nil {
		return err
	}
	schema.Not = &props
	return nil
}
