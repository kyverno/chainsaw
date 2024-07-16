package resource

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestLoad(t *testing.T) {
	baseDir := filepath.Join("..", "..", "..", "testdata", "resource")
	tests := []struct {
		fileName    string
		expectError bool
		expectedLen int
	}{{
		fileName:    filepath.Join(baseDir, "valid.yaml"),
		expectError: false,
		expectedLen: 2,
	}, {
		fileName:    filepath.Join(baseDir, "list.yaml"),
		expectError: false,
		expectedLen: 2,
	}, {
		fileName:    filepath.Join(baseDir, "empty.yaml"),
		expectError: true,
		expectedLen: 0,
	}, {
		fileName:    filepath.Join(baseDir, "invalid.yaml"),
		expectError: true,
		expectedLen: 0,
	}, {
		fileName:    filepath.Join(baseDir, "nonexistent.yaml"),
		expectError: true,
		expectedLen: 0,
	}, {
		fileName:    filepath.Join(baseDir, "folder-valid/*.yaml"),
		expectError: false,
		expectedLen: 3,
	}, {
		fileName:    filepath.Join(baseDir, "folder-invalid/*.yaml"),
		expectError: true,
		expectedLen: 0,
	}, {
		fileName:    filepath.Join(baseDir, "folder-nonexistent/*.yaml"),
		expectError: true,
		expectedLen: 0,
	}}
	for _, tt := range tests {
		resources, err := Load(tt.fileName, true)
		if !tt.expectError {
			assert.NoError(t, err)
			assert.Len(t, resources, tt.expectedLen)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestLoadFromURI(t *testing.T) {
	tests := []struct {
		fileName    string
		expectError bool
		expectedLen int
	}{{
		fileName:    "https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/resource/valid.yaml",
		expectError: false,
		expectedLen: 2,
	}, {
		fileName:    "https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/resource/empty.yaml",
		expectError: true,
		expectedLen: 0,
	}, {
		fileName:    "https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/resource/invalid.yaml",
		expectError: true,
		expectedLen: 0,
	}, {
		fileName:    "https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/resource/nonexistent.yaml",
		expectError: true,
		expectedLen: 0,
	}}
	for _, tt := range tests {
		url, err := url.ParseRequestURI(tt.fileName)
		assert.NoError(t, err)
		resources, err := LoadFromURI(url, true)
		if !tt.expectError {
			assert.NoError(t, err)
			assert.Len(t, resources, tt.expectedLen)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParse(t *testing.T) {
	baseDir := filepath.Join("..", "..", "..", "testdata", "resource")
	tests := []struct {
		fileName          string
		expectError       bool
		expectedLen       int
		expectedResources []unstructured.Unstructured
	}{{
		fileName:    filepath.Join(baseDir, "valid.yaml"),
		expectError: false,
		expectedLen: 2,
		expectedResources: []unstructured.Unstructured{{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name": "test-pod",
				},
			},
		}, {
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Service",
				"metadata": map[string]any{
					"name": "test-service",
				},
			},
		}},
	}, {
		fileName:          filepath.Join(baseDir, "empty.yaml"),
		expectError:       false,
		expectedLen:       0,
		expectedResources: []unstructured.Unstructured{},
	}, {
		fileName:    filepath.Join(baseDir, "invalid.yaml"),
		expectError: true,
		expectedLen: 0,
	}}
	for _, tt := range tests {
		content, readErr := os.ReadFile(tt.fileName)
		assert.NoError(t, readErr)
		resources, err := Parse(content, true)
		if !tt.expectError {
			assert.NoError(t, err)
			assert.Len(t, resources, tt.expectedLen)
			assert.ElementsMatch(t, tt.expectedResources, resources)
		} else {
			assert.Error(t, err)
		}
	}
}

func Test_parse(t *testing.T) {
	content, err := os.ReadFile("../../../testdata/resource/custom-resource.yaml")
	assert.NoError(t, err)
	tests := []struct {
		name      string
		splitter  splitter
		converter converter
		wantErr   bool
	}{{
		name:      "default behavior",
		splitter:  nil,
		converter: nil,
		wantErr:   false,
	}, {
		name: "splitter error",
		splitter: func([]byte) ([][]byte, error) {
			return nil, errors.New("splitter error")
		},
		converter: nil,
		wantErr:   true,
	}, {
		name:     "converter error",
		splitter: nil,
		converter: func([]byte) ([]byte, error) {
			return nil, errors.New("converter error")
		},
		wantErr: true,
	}, {
		name: "splitter and converter error",
		splitter: func([]byte) ([][]byte, error) {
			return nil, errors.New("splitter error")
		},
		converter: func([]byte) ([]byte, error) {
			return nil, errors.New("converter error")
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse(content, tt.splitter, tt.converter, true)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
