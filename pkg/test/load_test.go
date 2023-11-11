package test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	tloader "github.com/kyverno/chainsaw/pkg/internal/loader/testing"
	"github.com/kyverno/kyverno/ext/resource/loader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/openapi"
)

func TestLoad(t *testing.T) {
	basePath := "../../testdata/test"
	tests := []struct {
		name    string
		path    string
		want    []*v1alpha1.Test
		wantErr bool
	}{{
		name:    "confimap",
		path:    filepath.Join(basePath, "configmap.yaml"),
		wantErr: true,
	}, {
		name:    "not found",
		path:    filepath.Join(basePath, "not-found.yaml"),
		wantErr: true,
	}, {
		name:    "empty",
		path:    filepath.Join(basePath, "empty.yaml"),
		wantErr: true,
	}, {
		name:    "no spec",
		path:    filepath.Join(basePath, "no-spec.yaml"),
		wantErr: true,
	}, {
		name:    "no steps",
		path:    filepath.Join(basePath, "no-steps.yaml"),
		wantErr: true,
	}, {
		name:    "invalid step",
		path:    filepath.Join(basePath, "bad-step.yaml"),
		wantErr: true,
	}, {
		name: "ok",
		path: filepath.Join(basePath, "ok.yaml"),
		want: []*v1alpha1.Test{{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Test",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Spec: v1alpha1.TestSpec{
				Steps: []v1alpha1.TestSpecStep{{
					Spec: v1alpha1.TestStepSpec{
						Try: []v1alpha1.Operation{{
							Apply: &v1alpha1.Apply{
								FileRefOrResource: v1alpha1.FileRefOrResource{
									FileRef: v1alpha1.FileRef{
										File: "foo.yaml",
									},
								},
							},
						}},
						Catch: []v1alpha1.Catch{{
							PodLogs: &v1alpha1.PodLogs{
								Namespace: "foo",
							},
						}, {
							Events: &v1alpha1.Events{
								Namespace: "foo",
							},
						}, {
							Command: &v1alpha1.Command{
								Entrypoint: "time",
							},
						}, {
							Script: &v1alpha1.Script{
								Content: `echo "hello"`,
							},
						}},
					},
				}, {
					Spec: v1alpha1.TestStepSpec{
						Try: []v1alpha1.Operation{{
							Assert: &v1alpha1.Assert{
								FileRef: v1alpha1.FileRef{
									File: "bar.yaml",
								},
							},
						}},
						Finally: []v1alpha1.Finally{{
							PodLogs: &v1alpha1.PodLogs{
								Namespace: "foo",
							},
						}, {
							Events: &v1alpha1.Events{
								Namespace: "foo",
							},
						}, {
							Command: &v1alpha1.Command{
								Entrypoint: "time",
							},
						}, {
							Script: &v1alpha1.Script{
								Content: `echo "hello"`,
							},
						}},
					},
				}},
			},
		}},
	}, {
		name: "multiple",
		path: filepath.Join(basePath, "multiple.yaml"),
		want: []*v1alpha1.Test{{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Test",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Spec: v1alpha1.TestSpec{
				Steps: []v1alpha1.TestSpecStep{},
			},
		}, {
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Test",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-2",
			},
			Spec: v1alpha1.TestSpec{
				Steps: []v1alpha1.TestSpecStep{},
			},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.path)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_parse(t *testing.T) {
	content, err := os.ReadFile("../../testdata/test/custom-test.yaml")
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
		converter: func(unstructured.Unstructured) (*v1alpha1.Test, error) {
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
