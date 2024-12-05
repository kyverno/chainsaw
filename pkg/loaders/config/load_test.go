package config

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	tloader "github.com/kyverno/chainsaw/pkg/loaders/testing"
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *v1alpha2.Configuration
		wantErr bool
	}{{
		name:    "confimap",
		path:    "../../../testdata/config/configmap.yaml",
		wantErr: true,
	}, {
		name:    "not found",
		path:    "../../../testdata/config/not-found.yaml",
		wantErr: true,
	}, {
		name:    "empty",
		path:    "../../../testdata/config/empty.yaml",
		wantErr: true,
	}, {
		name:    "multiple",
		path:    "../../../testdata/config/multiple.yaml",
		wantErr: true,
	}, {
		name:    "bad catch",
		path:    "../../../testdata/config/v1alpha1/bad-catch.yaml",
		wantErr: true,
	}, {
		name: "default",
		path: "../../../testdata/config/v1alpha1/default.yaml",
		want: &v1alpha2.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
			Spec: v1alpha2.ConfigurationSpec{
				Timeouts: v1alpha1.DefaultTimeouts{
					Apply:   metav1.Duration{Duration: 5 * time.Second},
					Assert:  metav1.Duration{Duration: 30 * time.Second},
					Cleanup: metav1.Duration{Duration: 30 * time.Second},
					Delete:  metav1.Duration{Duration: 15 * time.Second},
					Error:   metav1.Duration{Duration: 30 * time.Second},
					Exec:    metav1.Duration{Duration: 5 * time.Second},
				},
				Discovery: v1alpha2.DiscoveryOptions{
					TestFile:         "chainsaw-test",
					FullName:         false,
					IncludeTestRegex: "",
					ExcludeTestRegex: "",
				},
				Cleanup: v1alpha2.CleanupOptions{
					SkipDelete: false,
				},
				Deletion: v1alpha2.DeletionOptions{
					Propagation: metav1.DeletePropagationBackground,
				},
				Execution: v1alpha2.ExecutionOptions{
					FailFast: false,
				},
				Report: &v1alpha2.ReportOptions{
					Format: "",
					Name:   "chainsaw-report",
				},
				Templating: v1alpha2.TemplatingOptions{
					Enabled: true,
				},
			},
		},
	}, {
		name: "custom-config",
		path: "../../../testdata/config/v1alpha1/custom-config.yaml",
		want: &v1alpha2.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "custom-config",
			},
			Spec: v1alpha2.ConfigurationSpec{
				Timeouts: v1alpha1.DefaultTimeouts{
					Apply:   metav1.Duration{Duration: 5 * time.Second},
					Assert:  metav1.Duration{Duration: 10 * time.Second},
					Cleanup: metav1.Duration{Duration: 5 * time.Second},
					Delete:  metav1.Duration{Duration: 5 * time.Second},
					Error:   metav1.Duration{Duration: 10 * time.Second},
					Exec:    metav1.Duration{Duration: 10 * time.Second},
				},
				Discovery: v1alpha2.DiscoveryOptions{
					TestFile:         "custom-test.yaml",
					FullName:         true,
					IncludeTestRegex: "include-*",
					ExcludeTestRegex: "exclude-*",
				},
				Cleanup: v1alpha2.CleanupOptions{
					SkipDelete: true,
				},
				Deletion: v1alpha2.DeletionOptions{
					Propagation: metav1.DeletePropagationBackground,
				},
				Execution: v1alpha2.ExecutionOptions{
					FailFast:                    true,
					Parallel:                    ptr.To(4),
					ForceTerminationGracePeriod: &metav1.Duration{Duration: 10 * time.Second},
				},
				Report: &v1alpha2.ReportOptions{
					Format: "JSON",
					Name:   "custom-report",
				},
				Templating: v1alpha2.TemplatingOptions{
					Enabled: true,
				},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(fsutils.NewLocal(), tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parse(t *testing.T) {
	content, err := os.ReadFile("../../../testdata/config/v1alpha1/custom-config.yaml")
	assert.NoError(t, err)
	tests := []struct {
		name          string
		splitter      splitter
		loaderFactory loaderFactory
		converter     converter
		wantErr       bool
	}{{
		name:          "default",
		splitter:      nil,
		loaderFactory: nil,
		converter:     nil,
		wantErr:       false,
	}, {
		name: "splitter error",
		splitter: func([]byte) ([][]byte, error) {
			return nil, errors.New("splitter")
		},
		loaderFactory: nil,
		converter:     nil,
		wantErr:       true,
	}, {
		name:     "loader factory error",
		splitter: nil,
		loaderFactory: func() (loader.Loader, error) {
			return nil, errors.New("loader factory")
		},
		converter: nil,
		wantErr:   true,
	}, {
		name:     "loader error",
		splitter: nil,
		loaderFactory: func() (loader.Loader, error) {
			return &tloader.FakeLoader{
				LoadFn: func(_ int, _ []byte) (schema.GroupVersionKind, unstructured.Unstructured, error) {
					return schema.GroupVersionKind{Group: "v1", Kind: "Something"}, unstructured.Unstructured{}, nil
				},
			}, nil
		},
		converter: nil,
		wantErr:   true,
	}, {
		name:          "converter error",
		splitter:      nil,
		loaderFactory: nil,
		converter: func(schema.GroupVersionKind, unstructured.Unstructured) (*v1alpha2.Configuration, error) {
			return nil, errors.New("converter")
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse(content, tt.splitter, tt.loaderFactory, tt.converter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
