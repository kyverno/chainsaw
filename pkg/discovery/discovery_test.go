package discovery

import (
	"os"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
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
							Try: []v1alpha1.Operation{{
								Apply: &v1alpha1.Apply{
									FileRefOrResource: v1alpha1.FileRefOrResource{
										FileRef: v1alpha1.FileRef{
											File: "configmap.yaml",
										},
									},
								},
							}},
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
		}},
		wantErr: false,
	}, {
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
										File: "01-errors.yaml",
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

func TestHelpDiscoverTests(t *testing.T) {
	type testCase struct {
		name              string
		folders           []string
		expectedTestCount int
		expectError       bool
	}

	testCases := []testCase{
		{
			name:              "Successful Discovery",
			folders:           []string{"../../testdata/discovery/test"},
			expectedTestCount: 1,
			expectError:       false,
		},
		{
			name:              "LoadTest Returns Error",
			folders:           []string{"folder1"},
			expectedTestCount: 0,
			expectError:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tests, err := helpDiscoverTests("chainsaw-test.yaml", func() []string {
				return tc.folders
			})
			if tc.expectError {
				require.Error(t, err, "Expected an error but got none")
			} else {
				require.NoError(t, err, "Expected no error but got one")
			}
			assert.Equal(t, tc.expectedTestCount, len(tests), "Unexpected number of tests returned")
		})
	}

}

func TestDiscoverTests_UnreadableFolder(t *testing.T) {
	tempDir := t.TempDir()

	err := os.Chmod(tempDir, 0000)
	if err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}

	_, err = DiscoverTests("chainsaw-test.yaml", tempDir)
	assert.Error(t, err, "Expected an error for unreadable folder")
}
