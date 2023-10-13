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
		name:     "ok",
		fileName: "chainsaw-test.yaml",
		paths:    []string{"../../testdata/tests"},
		want: []Test{{
			BasePath: "../../testdata/tests",
			Test: &v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "chainsaw-quick-start",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestStepSpec{{
						Apply: []v1alpha1.Apply{{
							File: "configmap.yaml",
						}},
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
