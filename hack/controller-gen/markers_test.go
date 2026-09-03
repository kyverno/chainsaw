package main

import (
	"reflect"
	"testing"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-tools/pkg/crd"
	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
)

func TestSchemaMarkers_ApplyToSchema(t *testing.T) {
	tests := []struct {
		name    string
		marker  crd.SchemaMarker
		schema  apiext.JSONSchemaProps
		want    apiext.JSONSchemaProps
		wantErr bool
	}{
		{
			name:   "oneOf appends without replacing existing rules",
			marker: OneOf{Value: map[string]any{"required": []string{"script"}}},
			schema: apiext.JSONSchemaProps{OneOf: []apiext.JSONSchemaProps{{Required: []string{"command"}}}},
			want: apiext.JSONSchemaProps{OneOf: []apiext.JSONSchemaProps{
				{Required: []string{"command"}},
				{Required: []string{"script"}},
			}},
		},
		{
			name: "not preserves nested constraints",
			marker: Not{Value: map[string]any{
				"required": []string{"use"},
				"anyOf":    []any{map[string]any{"required": []string{"try"}}},
			}},
			want: apiext.JSONSchemaProps{Not: &apiext.JSONSchemaProps{
				Required: []string{"use"},
				AnyOf:    []apiext.JSONSchemaProps{{Required: []string{"try"}}},
			}},
		},
		{
			name:    "oneOf rejects values that cannot be marshaled",
			marker:  OneOf{Value: make(chan int)},
			wantErr: true,
		},
		{
			name:    "not rejects values that cannot be marshaled",
			marker:  Not{Value: make(chan int)},
			wantErr: true,
		},
		{
			name:    "oneOf rejects invalid schemas",
			marker:  OneOf{Value: map[string]any{"required": 42}},
			wantErr: true,
		},
		{
			name:    "not rejects invalid schemas",
			marker:  Not{Value: map[string]any{"required": 42}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.marker.ApplyToSchema(&crdmarkers.SchemaContext{}, &tt.schema)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ApplyToSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.schema, tt.want) {
				t.Errorf("ApplyToSchema() schema = %#v, want %#v", tt.schema, tt.want)
			}
		})
	}
}
