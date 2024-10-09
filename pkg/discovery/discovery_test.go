package discovery

import (
	"os"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/model"
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
	"github.com/stretchr/testify/assert"
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
							Try: []v1alpha1.Operation{{
								Apply: &v1alpha1.Apply{
									ActionResourceRef: v1alpha1.ActionResourceRef{
										FileRef: v1alpha1.FileRef{
											File: "configmap.yaml",
										},
									},
								},
							}},
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
		paths:    []string{"../../testdata/discovery/manifests"},
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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DiscoverTests(fsutils.NewLocal(), tt.fileName, nil, false, tt.paths...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
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
	testCases := []testCase{{
		name:              "Successful Discovery",
		folders:           []string{"../../testdata/discovery/test"},
		expectedTestCount: 1,
		expectError:       false,
	}, {
		name:              "LoadTest Returns Error",
		folders:           []string{"folder1"},
		expectedTestCount: 0,
		expectError:       true,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tests, err := discoverTests("chainsaw-test.yaml", nil, false, tc.folders...)
			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
			assert.Equal(t, tc.expectedTestCount, len(tests), "Unexpected number of tests returned")
		})
	}
}

func TestDiscoverTests_UnreadableFolder(t *testing.T) {
	tempDir := t.TempDir()
	err := os.Chmod(tempDir, 0o000)
	if err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}
	_, err = DiscoverTests(fsutils.NewLocal(), "chainsaw-test.yaml", nil, false, tempDir)
	assert.Error(t, err, "Expected an error for unreadable folder")
}
