package discovery

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDiscoverTests(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		paths    []string
		want     []Test
		wantErr  bool
	}{{
		name:     "quick start",
		fileName: "chainsaw-test.yaml",
		paths:    []string{"../../testdata/tests/quick-start"},
		want: []Test{{
			BasePath: "../../testdata/tests/quick-start",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "chainsaw-quick-start",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestSpecStep{{
						Name: "create configmap",
						Spec: v1alpha1.TestStepSpec{
							Apply: []v1alpha1.Apply{{
								FileRef: v1alpha1.FileRef{
									File: "configmap.yaml",
								},
							}},
						},
					}, {
						Name: "assert configmap",
						Spec: v1alpha1.TestStepSpec{
							Assert: []v1alpha1.Assert{{
								FileRef: v1alpha1.FileRef{
									File: "configmap.yaml",
								},
							}},
						},
					}},
				},
			},
		}},
		wantErr: false,
	}, {
		name:     "manifests based",
		fileName: "chainsaw-test.yaml",
		paths:    []string{"../../testdata/tests/manifests-based"},
		want: []Test{{
			BasePath: "../../testdata/tests/manifests-based",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "manifests-based",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestSpecStep{{
						Name: "configmap",
						Spec: v1alpha1.TestStepSpec{
							Apply: []v1alpha1.Apply{{
								FileRef: v1alpha1.FileRef{
									File: "01-configmap.yaml",
								},
							}},
							Assert: []v1alpha1.Assert{{
								FileRef: v1alpha1.FileRef{
									File: "01-assert.yaml",
								},
							}},
						},
					}},
				},
			},
		}},
		wantErr: false,
	}, {
		name:     "steps based",
		fileName: "chainsaw-test.yaml",
		paths:    []string{"../../testdata/tests/steps-based"},
		want: []Test{{
			BasePath: "../../testdata/tests/steps-based",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "steps-based",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestSpecStep{{
						Name: "test-1",
						Spec: v1alpha1.TestStepSpec{
							Apply: []v1alpha1.Apply{{
								FileRef: v1alpha1.FileRef{
									File: "foo.yaml",
								},
							}},
							Assert: []v1alpha1.Assert{{
								FileRef: v1alpha1.FileRef{
									File: "bar.yaml",
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
			got, err := DiscoverTests(tt.fileName, tt.paths...)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
