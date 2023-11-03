package config

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
)

// MockLoader is a mock implementation of the loader.Loader interface for testing.
type mockLoader struct {
	LoadFunc func(document []byte) (schema.GroupVersionKind, unstructured.Unstructured, error)
}

func (m *mockLoader) Load(document []byte) (schema.GroupVersionKind, unstructured.Unstructured, error) {
	return m.LoadFunc(document)
}

// mockDocumentParser is our test implementation that satisfies the DocumentParser interface.
type mockDocumentParser struct {
	SplitDocumentsFunc func(content []byte) ([][]byte, error)
}

func (m *mockDocumentParser) SplitDocuments(content []byte) ([][]byte, error) {
	return m.SplitDocumentsFunc(content)
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *v1alpha1.Configuration
		wantErr bool
	}{
		{
			name:    "confimap",
			path:    "../../testdata/config/configmap.yaml",
			wantErr: true,
		},
		{
			name:    "not found",
			path:    "../../testdata/config/not-found.yaml",
			wantErr: true,
		},
		{
			name:    "empty",
			path:    "../../testdata/config/empty.yaml",
			wantErr: true,
		},
		{
			name: "default",
			path: "../../testdata/config/default.yaml",
			want: &v1alpha1.Configuration{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Configuration",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: v1alpha1.ConfigurationSpec{
					SkipDelete:       false,
					FailFast:         false,
					ReportFormat:     "",
					ReportName:       "chainsaw-report",
					FullName:         false,
					IncludeTestRegex: "",
					ExcludeTestRegex: "",
				},
			},
		},
		{
			name: "custom-config",
			path: "../../testdata/config/custom-config.yaml",
			want: &v1alpha1.Configuration{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Configuration",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "custom-config",
				},
				Spec: v1alpha1.ConfigurationSpec{
					Timeouts: v1alpha1.Timeouts{
						Apply:   &metav1.Duration{Duration: 5 * time.Second},
						Assert:  &metav1.Duration{Duration: 10 * time.Second},
						Error:   &metav1.Duration{Duration: 10 * time.Second},
						Delete:  &metav1.Duration{Duration: 5 * time.Second},
						Cleanup: &metav1.Duration{Duration: 5 * time.Second},
						Exec:    &metav1.Duration{Duration: 10 * time.Second},
					},
					SkipDelete:       true,
					FailFast:         true,
					Parallel:         ptr.To(4),
					ReportFormat:     "JSON",
					ReportName:       "custom-report",
					FullName:         true,
					IncludeTestRegex: "include-*",
					ExcludeTestRegex: "exclude-*",
				},
			},
		},
		{
			name:    "multiple",
			path:    "../../testdata/config/multiple.yaml",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		content    []byte
		mockReturn [][]byte
		mockErr    error
		wantErr    bool
	}{
		{
			name:       "error from yaml.SplitDocuments",
			content:    []byte("invalid content"),
			mockReturn: nil,
			mockErr:    errors.New("failed to split documents"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := &mockDocumentParser{
				SplitDocumentsFunc: func(content []byte) ([][]byte, error) {
					return tt.mockReturn, tt.mockErr
				},
			}
			_, err := Parse(tt.content, mockParser)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewLoader(t *testing.T) {
	_, err := newLoader(data.Crds(), data.CrdsFolder)
	if err != nil {
		t.Errorf("newLoader with valid input returned an error: %v", err)
	}

	_, err = newLoader(nil, "non-existent-folder")
	if err == nil {
		t.Errorf("newLoader with invalid input did not return an error")
	}
}

func TestParseDocument(t *testing.T) {
	testCases := []struct {
		name                  string
		document              []byte
		mockLoad              func(document []byte) (schema.GroupVersionKind, unstructured.Unstructured, error)
		expectedConfiguration *v1alpha1.Configuration
		expectError           bool
	}{
		{
			name: "ValidConfiguration",
			document: []byte(`apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Configuration
metadata:
  name: test-configuration`),
			mockLoad: func(document []byte) (schema.GroupVersionKind, unstructured.Unstructured, error) {
				u := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "chainsaw.kyverno.io/v1alpha1",
						"kind":       "Configuration",
						"metadata": map[string]interface{}{
							"name": "test-configuration",
						},
						"spec": map[string]interface{}{},
					},
				}
				return configuration_v1alpha1, *u, nil
			},
			expectedConfiguration: &v1alpha1.Configuration{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Configuration",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-configuration",
				},
				Spec: v1alpha1.ConfigurationSpec{},
			},
			expectError: false,
		},
		{
			name: "InvalidGVK",
			document: []byte(`apiVersion: chainsaw.kyverno.io/v1alpha1
	kind: Configuration
	metadata:
	  name: wrong-configuration`),
			mockLoad: func(document []byte) (schema.GroupVersionKind, unstructured.Unstructured, error) {
				wrongGVK := schema.GroupVersionKind{Group: "wrongGroup", Version: "v1alpha1", Kind: "WrongKind"}
				u := unstructured.Unstructured{}
				return wrongGVK, u, nil
			},
			expectedConfiguration: nil,
			expectError:           true,
		},
		{
			name: "LoaderError",
			document: []byte(`apiVersion: chainsaw.kyverno.io/v1alpha1
	kind: Configuration
	metadata:
	  name: test-configuration`),
			mockLoad: func(document []byte) (schema.GroupVersionKind, unstructured.Unstructured, error) {
				// Simulate an error from the loader
				return schema.GroupVersionKind{}, unstructured.Unstructured{}, errors.New("loader error")
			},
			expectedConfiguration: nil,
			expectError:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLoader := &mockLoader{LoadFunc: tc.mockLoad}

			configuration, err := parseDocument(mockLoader, tc.document)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfiguration, configuration)
			}
		})
	}
}

func TestParseDocuments(t *testing.T) {
	badCRDFolder := func(path string) string {
		return "/invalid/path/" + path
	}

	documents := [][]byte{
		[]byte("apiVersion: v1alpha1\nkind: Configuration\n..."),
	}

	_, err := parseDocuments(documents, badCRDFolder)
	assert.Error(t, err, "expected an error due to invalid CRD folder path")
}
