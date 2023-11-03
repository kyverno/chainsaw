package test

import (
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
						Try: []v1alpha1.Operation{
							{
								Apply: []v1alpha1.Apply{{
									FileRef: v1alpha1.FileRef{
										File: "foo.yaml",
									},
								}},
							},
						},
						Catch: []v1alpha1.Catch{{
							Collect: &v1alpha1.Collect{
								PodLogs: &v1alpha1.PodLogs{
									Namespace: "foo",
								},
							},
						}, {
							Collect: &v1alpha1.Collect{
								Events: &v1alpha1.Events{
									Namespace: "foo",
								},
							},
						}, {
							Exec: &v1alpha1.Exec{
								Command: &v1alpha1.Command{
									Entrypoint: "time",
								},
							},
						}, {
							Exec: &v1alpha1.Exec{
								Script: &v1alpha1.Script{
									Content: `echo "hello"`,
								},
							},
						}},
					},
				}, {
					Spec: v1alpha1.TestStepSpec{
						Try: []v1alpha1.Operation{
							{
								Assert: []v1alpha1.Assert{{
									FileRef: v1alpha1.FileRef{
										File: "bar.yaml",
									},
								}},
							},
						},
						Finally: []v1alpha1.Finally{{
							Collect: &v1alpha1.Collect{
								PodLogs: &v1alpha1.PodLogs{
									Namespace: "foo",
								},
							},
						}, {
							Collect: &v1alpha1.Collect{
								Events: &v1alpha1.Events{
									Namespace: "foo",
								},
							},
						}, {
							Exec: &v1alpha1.Exec{
								Command: &v1alpha1.Command{
									Entrypoint: "time",
								},
							},
						}, {
							Exec: &v1alpha1.Exec{
								Script: &v1alpha1.Script{
									Content: `echo "hello"`,
								},
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
