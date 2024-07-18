package values

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	basePath := "../../../testdata/values"
	tests := []struct {
		name    string
		paths   []string
		want    map[string]any
		wantErr bool
	}{{
		name:    "values-1",
		paths:   []string{filepath.Join(basePath, "values-1.yaml")},
		wantErr: false,
		want: map[string]any{
			"foo": map[string]any{
				"bar": "baz",
			},
			"test": 42.0,
		},
	}, {
		name:    "values-2",
		paths:   []string{filepath.Join(basePath, "values-2.yaml")},
		wantErr: false,
		want: map[string]any{
			"foo": map[string]any{
				"bar": "baz",
			},
			"test": nil,
		},
	}, {
		name:    "values-1 and values-2",
		paths:   []string{filepath.Join(basePath, "values-1.yaml"), filepath.Join(basePath, "values-2.yaml")},
		wantErr: false,
		want: map[string]any{
			"foo": map[string]any{
				"bar": "baz",
			},
			"test": nil,
		},
	}, {
		name:    "values-2 and values-1",
		paths:   []string{filepath.Join(basePath, "values-2.yaml"), filepath.Join(basePath, "values-1.yaml")},
		wantErr: false,
		want: map[string]any{
			"foo": map[string]any{
				"bar": "baz",
			},
			"test": 42.0,
		},
	}, {
		name:    "not found",
		paths:   []string{filepath.Join(basePath, "not-found.yaml")},
		wantErr: true,
	}, {
		name:    "empty",
		paths:   []string{filepath.Join(basePath, "empty.yaml")},
		want:    map[string]any{},
		wantErr: false,
	}, {
		name:    "invalid",
		paths:   []string{filepath.Join(basePath, "invalid.yaml")},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.paths...)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
