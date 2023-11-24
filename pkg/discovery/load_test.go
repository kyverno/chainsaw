package discovery

import (
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestLoadTest(t *testing.T) {
	basePath := "../../testdata/discovery"
	tests := []struct {
		name     string
		fileName string
		path     string
		want     *Test
		wantErr  bool
	}{{
		name:     "invalid path",
		fileName: "chainsaw-test.yaml",
		path:     "/invalid",
		want:     nil,
		wantErr:  true,
	}, {
		name:     "no path",
		fileName: "chainsaw-test.yaml",
		path:     "",
		want:     nil,
		wantErr:  true,
	}, {
		name:     "test",
		fileName: "chainsaw-test.yaml",
		path:     filepath.Join(basePath, "test"),
		want: &Test{
			BasePath: "../../testdata/discovery/test",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestSpecStep{{
						Name: "create configmap",
						Spec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{
								{
									Apply: &v1alpha1.Apply{
										FileRefOrResource: v1alpha1.FileRefOrResource{
											FileRef: v1alpha1.FileRef{
												File: "configmap.yaml",
											},
										},
									},
								},
							},
						},
					}, {
						Name: "assert configmap",
						Spec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Assert: &v1alpha1.Assert{
									FileRef: v1alpha1.FileRef{
										File: "configmap.yaml",
									},
								},
							}},
						},
					}},
				},
			},
		},
		wantErr: false,
	}, {
		name:     "test",
		fileName: "chainsaw-test.yaml",
		path:     filepath.Join(basePath, "invalid-test-and-steps"),
		want: &Test{
			BasePath: "../../testdata/discovery/invalid-test-and-steps",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestSpecStep{{
						Name: "create configmap",
						Spec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{
								{
									Apply: &v1alpha1.Apply{
										FileRefOrResource: v1alpha1.FileRefOrResource{
											FileRef: v1alpha1.FileRef{
												File: "configmap.yaml",
											},
										},
									},
								},
							},
						},
					}, {
						Name: "assert configmap",
						Spec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{
								{
									Assert: &v1alpha1.Assert{
										FileRef: v1alpha1.FileRef{
											File: "configmap.yaml",
										},
									},
								},
							},
						},
					}},
				},
			},
		},
		wantErr: false,
	}, {
		name:     "test",
		fileName: "chainsaw-test.yaml",
		path:     filepath.Join(basePath, "test-and-steps"),
		want: &Test{
			BasePath: "../../testdata/discovery/test-and-steps",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.TestSpec{
					Skip:       ptr.To(false),
					Concurrent: ptr.To(false),
					Steps: []v1alpha1.TestSpecStep{{
						Name: "test-1",
						Spec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Apply: &v1alpha1.Apply{
									FileRefOrResource: v1alpha1.FileRefOrResource{
										FileRef: v1alpha1.FileRef{
											File: "foo.yaml",
										},
									},
								},
							}, {
								Assert: &v1alpha1.Assert{
									FileRef: v1alpha1.FileRef{
										File: "bar.yaml",
									},
								},
							}},
						},
					}},
				},
			},
		},
		wantErr: false,
	}, {
		name:     "steps",
		fileName: "chainsaw-test.yaml",
		path:     filepath.Join(basePath, "steps"),
		want: &Test{
			BasePath: "../../testdata/discovery/steps",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "steps",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestSpecStep{{
						Name: "test-1",
						Spec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Apply: &v1alpha1.Apply{
									FileRefOrResource: v1alpha1.FileRefOrResource{
										FileRef: v1alpha1.FileRef{
											File: "foo.yaml",
										},
									},
								},
							}, {
								Assert: &v1alpha1.Assert{
									FileRef: v1alpha1.FileRef{
										File: "bar.yaml",
									},
								},
							}},
						},
					}},
				},
			},
		},
		wantErr: false,
	}, {
		name:     "steps",
		fileName: "chainsaw-test.yaml",
		path:     filepath.Join(basePath, "steps-and-manifests"),
		want: &Test{
			BasePath: "../../testdata/discovery/steps-and-manifests",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "steps-and-manifests",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestSpecStep{{
						Name: "test-1",
						Spec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Apply: &v1alpha1.Apply{
									FileRefOrResource: v1alpha1.FileRefOrResource{
										FileRef: v1alpha1.FileRef{
											File: "foo.yaml",
										},
									},
								},
							}, {
								Assert: &v1alpha1.Assert{
									FileRef: v1alpha1.FileRef{
										File: "bar.yaml",
									},
								},
							}, {
								Apply: &v1alpha1.Apply{
									FileRefOrResource: v1alpha1.FileRefOrResource{
										FileRef: v1alpha1.FileRef{
											File: "01-configmap.yaml",
										},
									},
								},
							}, {
								Assert: &v1alpha1.Assert{
									FileRef: v1alpha1.FileRef{
										File: "01-assert.yaml",
									},
								},
							}, {
								Error: &v1alpha1.Error{
									FileRef: v1alpha1.FileRef{
										File: "01-error.yaml",
									},
								},
							}},
						},
					}},
				},
			},
		},
		wantErr: false,
	}, {
		name:     "manifests",
		fileName: "chainsaw-test.yaml",
		path:     filepath.Join(basePath, "manifests"),
		want: &Test{
			BasePath: "../../testdata/discovery/manifests",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "manifests",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestSpecStep{{
						Name: "configmap",
						Spec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Apply: &v1alpha1.Apply{
									FileRefOrResource: v1alpha1.FileRefOrResource{
										FileRef: v1alpha1.FileRef{
											File: "01-configmap.yaml",
										},
									},
								},
							}, {
								Assert: &v1alpha1.Assert{
									FileRef: v1alpha1.FileRef{
										File: "01-assert.yaml",
									},
								},
							}, {
								Error: &v1alpha1.Error{
									FileRef: v1alpha1.FileRef{
										File: "01-error.yaml",
									},
								},
							}},
						},
					}},
				},
			},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadTest(tt.fileName, tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.want != nil {
					tt.want.Err = nil
				}
				if got != nil {
					got.Err = nil
				}
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
