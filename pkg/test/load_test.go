package test

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    []*v1alpha1.Test
		wantErr bool
	}{{
		name:    "confimap",
		path:    "../../testdata/test/configmap.yaml",
		wantErr: true,
	}, {
		name:    "not found",
		path:    "../../testdata/test/not-found.yaml",
		wantErr: true,
	}, {
		name:    "empty",
		path:    "../../testdata/test/empty.yaml",
		wantErr: true,
	}, {
		name:    "no spec",
		path:    "../../testdata/test/no-spec.yaml",
		wantErr: true,
	}, {
		name:    "no steps",
		path:    "../../testdata/test/no-steps.yaml",
		wantErr: true,
	}, {
		name:    "invalid step",
		path:    "../../testdata/test/bad-step.yaml",
		wantErr: true,
	}, {
		name: "ok",
		path: "../../testdata/test/ok.yaml",
		want: []*v1alpha1.Test{{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Test",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Spec: v1alpha1.TestSpec{
				Steps: []v1alpha1.TestStepSpec{{
					Apply: []v1alpha1.Apply{{
						File: "foo.yaml",
					}},
				}, {
					Assert: []v1alpha1.Assert{{
						File: "bar.yaml",
					}},
				}},
			},
		}},
	}, {
		name: "multiple",
		path: "../../testdata/test/multiple.yaml",
		want: []*v1alpha1.Test{{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Test",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
			Spec: v1alpha1.TestSpec{
				Steps: []v1alpha1.TestStepSpec{},
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
				Steps: []v1alpha1.TestStepSpec{},
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
