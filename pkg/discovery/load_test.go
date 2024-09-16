package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLoadTest(t *testing.T) {
	basePath := "../../testdata/discovery"
	tests := []struct {
		name     string
		fileName string
		path     string
		want     []Test
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
		want: []Test{{
			BasePath: "../../testdata/discovery/test",
			Test: &model.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestStep{{
						Name: "create configmap",
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{
								{
									Apply: &v1alpha1.Apply{
										ActionResourceRef: v1alpha1.ActionResourceRef{
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
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Assert: &v1alpha1.Assert{
									ActionCheckRef: v1alpha1.ActionCheckRef{
										FileRef: v1alpha1.FileRef{
											File: "configmap.yaml",
										},
									},
								},
							}},
						},
					}},
				},
			},
		}},
		wantErr: false,
	}, {
		name:     "test (no extension - yaml)",
		fileName: "chainsaw-test",
		path:     filepath.Join(basePath, "test"),
		want: []Test{{
			BasePath: "../../testdata/discovery/test",
			Test: &model.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestStep{{
						Name: "create configmap",
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{
								{
									Apply: &v1alpha1.Apply{
										ActionResourceRef: v1alpha1.ActionResourceRef{
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
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Assert: &v1alpha1.Assert{
									ActionCheckRef: v1alpha1.ActionCheckRef{
										FileRef: v1alpha1.FileRef{
											File: "configmap.yaml",
										},
									},
								},
							}},
						},
					}},
				},
			},
		}},
		wantErr: false,
	}, {
		name:     "test (no extension - yml)",
		fileName: "chainsaw-test",
		path:     filepath.Join(basePath, "test-yml"),
		want: []Test{{
			BasePath: "../../testdata/discovery/test-yml",
			Test: &model.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestStep{{
						Name: "create configmap",
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{
								{
									Apply: &v1alpha1.Apply{
										ActionResourceRef: v1alpha1.ActionResourceRef{
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
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Assert: &v1alpha1.Assert{
									ActionCheckRef: v1alpha1.ActionCheckRef{
										FileRef: v1alpha1.FileRef{
											File: "configmap.yaml",
										},
									},
								},
							}},
						},
					}},
				},
			},
		}},
		wantErr: false,
	}, {
		name:     "manifests",
		fileName: "chainsaw-test.yaml",
		path:     filepath.Join(basePath, "manifests"),
		want: []Test{{
			BasePath: "../../testdata/discovery/manifests",
			Test: &model.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "manifests",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestStep{{
						Name: "step-01",
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Apply: &v1alpha1.Apply{
									ActionResourceRef: v1alpha1.ActionResourceRef{
										FileRef: v1alpha1.FileRef{
											File: "01-configmap.yaml",
										},
									},
								},
							}, {
								Assert: &v1alpha1.Assert{
									ActionCheckRef: v1alpha1.ActionCheckRef{
										FileRef: v1alpha1.FileRef{
											File: "01-assert.yaml",
										},
									},
								},
							}, {
								Error: &v1alpha1.Error{
									ActionCheckRef: v1alpha1.ActionCheckRef{
										FileRef: v1alpha1.FileRef{
											File: "01-errors.yaml",
										},
									},
								},
							}},
						},
					}},
				},
			},
		}},
		wantErr: false,
	}, {
		name:     "empty test",
		fileName: "",
		path:     filepath.Join(basePath, "empty-test"),
		want:     nil,
		wantErr:  false,
	}, {
		name:     "multiple tests",
		fileName: "chainsaw-test.yaml",
		path:     filepath.Join(basePath, "multiple-tests"),
		want: []Test{{
			BasePath: "../../testdata/discovery/multiple-tests",
			Test: &model.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestStep{{
						Name: "create configmap",
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{
								{
									Apply: &v1alpha1.Apply{
										ActionResourceRef: v1alpha1.ActionResourceRef{
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
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Assert: &v1alpha1.Assert{
									ActionCheckRef: v1alpha1.ActionCheckRef{
										FileRef: v1alpha1.FileRef{
											File: "configmap.yaml",
										},
									},
								},
							}},
						},
					}},
				},
			},
		}, {
			BasePath: "../../testdata/discovery/multiple-tests",
			Test: &model.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-2",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestStep{{
						Name: "create configmap",
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{
								{
									Apply: &v1alpha1.Apply{
										ActionResourceRef: v1alpha1.ActionResourceRef{
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
						TestStepSpec: v1alpha1.TestStepSpec{
							Try: []v1alpha1.Operation{{
								Assert: &v1alpha1.Assert{
									ActionCheckRef: v1alpha1.ActionCheckRef{
										FileRef: v1alpha1.FileRef{
											File: "configmap.yaml",
										},
									},
								},
							}},
						},
					}},
				},
			},
		}},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadTest(tt.fileName, tt.path, false)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_tryLoadTest(t *testing.T) {
	dir := t.TempDir()
	fileName := "chainsaw-test.yaml"
	filePath := filepath.Join(dir, fileName)
	_, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	err = os.Chmod(filePath, 0o000)
	if err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}
	_, err = tryLoadTestFile(filePath, false)
	assert.Error(t, err)
}
