package config

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	tloader "github.com/kyverno/chainsaw/pkg/internal/loader/testing"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/openapi"
	"k8s.io/utils/ptr"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *v1alpha1.Configuration
		wantErr bool
	}{{
		name:    "confimap",
		path:    "../../testdata/config/configmap.yaml",
		wantErr: true,
	}, {
		name:    "not found",
		path:    "../../testdata/config/not-found.yaml",
		wantErr: true,
	}, {
		name:    "empty",
		path:    "../../testdata/config/empty.yaml",
		wantErr: true,
	}, {
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
				TestFile:         "chainsaw-test",
				SkipDelete:       false,
				FailFast:         false,
				ReportFormat:     "",
				ReportName:       "chainsaw-report",
				FullName:         false,
				IncludeTestRegex: "",
				ExcludeTestRegex: "",
			},
		},
	}, {
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
				TestFile:                    "custom-test.yaml",
				SkipDelete:                  true,
				FailFast:                    true,
				Parallel:                    ptr.To(4),
				ReportFormat:                "JSON",
				ReportName:                  "custom-report",
				FullName:                    true,
				IncludeTestRegex:            "include-*",
				ExcludeTestRegex:            "exclude-*",
				ForceTerminationGracePeriod: &metav1.Duration{Duration: 10 * time.Second},
			},
		},
	}, {
		name:    "multiple",
		path:    "../../testdata/config/multiple.yaml",
		wantErr: true,
	}}
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

func Test_parse(t *testing.T) {
	content, err := os.ReadFile("../../testdata/config/custom-config.yaml")
	assert.NoError(t, err)
	tests := []struct {
		name          string
		splitter      splitter
		loaderFactory loaderFactory
		converter     converter
		validator     validator
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
		loaderFactory: func(openapi.Client) (loader.Loader, error) {
			return nil, errors.New("loader factory")
		},
		converter: nil,
		wantErr:   true,
	}, {
		name:     "loader error",
		splitter: nil,
		loaderFactory: func(openapi.Client) (loader.Loader, error) {
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
		converter: func(schema.GroupVersionKind, unstructured.Unstructured) (*v1alpha1.Configuration, error) {
			return nil, errors.New("converter")
		},
		wantErr: true,
	}, {
		name:          "validator error",
		splitter:      nil,
		loaderFactory: nil,
		converter:     nil,
		validator: func(obj *v1alpha1.Configuration) field.ErrorList {
			return field.ErrorList{
				field.Invalid(nil, nil, ""),
			}
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse(content, tt.splitter, tt.loaderFactory, tt.converter, tt.validator)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
