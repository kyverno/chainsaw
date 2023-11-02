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
		name:     "test",
		fileName: "chainsaw-test.yaml",
		paths:    []string{"../../testdata/discovery/test"},
		want: []Test{{
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
							Operations: []v1alpha1.Operations{{
								Apply: []v1alpha1.Apply{{
									FileRef: v1alpha1.FileRef{
										File: "configmap.yaml",
									},
								}},
							},
							},
						},
					}, {
						Name: "assert configmap",
						Spec: v1alpha1.TestStepSpec{
							Operations: []v1alpha1.Operations{{
								Assert: []v1alpha1.Assert{{
									FileRef: v1alpha1.FileRef{
										File: "configmap.yaml",
									},
								}},
							},
							},
						},
					}},
				},
			},
		}},
		wantErr: false,
	},
		{
			name:     "manifests",
			fileName: "chainsaw-test.yaml",
			paths:    []string{"../../testdata/discovery/manifests"},
			want: []Test{{
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
							Name: "assert",
							Spec: v1alpha1.TestStepSpec{
								Operations: []v1alpha1.Operations{
									{

										Assert: []v1alpha1.Assert{{
											FileRef: v1alpha1.FileRef{
												File: "01-assert.yaml",
											},
										}},
									},
									{
										Apply: []v1alpha1.Apply{{
											FileRef: v1alpha1.FileRef{
												File: "01-configmap.yaml",
											},
										}},
									},
									{
										Error: []v1alpha1.Error{{
											FileRef: v1alpha1.FileRef{
												File: "01-error.yaml",
											},
										}},
									},
								},
							},
						}},
					},
				},
			}},
			wantErr: false,
		},
		{
			name:     "steps",
			fileName: "chainsaw-test.yaml",
			paths:    []string{"../../testdata/discovery/steps"},
			want: []Test{{
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
								Operations: []v1alpha1.Operations{
									{
										Assert: []v1alpha1.Assert{{
											FileRef: v1alpha1.FileRef{
												File: "bar.yaml",
											},
										}},
										Apply: []v1alpha1.Apply{{
											FileRef: v1alpha1.FileRef{
												File: "foo.yaml",
											},
										}},
									},
								},
							},
						}},
					},
				},
			}},
			wantErr: false,
		},
	}
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
