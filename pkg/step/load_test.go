package step

import (
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLoad(t *testing.T) {
	basePath := "../../testdata/step"
	tests := []struct {
		name    string
		path    string
		want    []*v1alpha1.TestStep
		wantErr bool
	}{
		{
			name:    "confimap",
			path:    filepath.Join(basePath, "configmap.yaml"),
			wantErr: true,
		},
		{
			name:    "not found",
			path:    filepath.Join(basePath, "not-found.yaml"),
			wantErr: true,
		},
		{
			name:    "empty",
			path:    filepath.Join(basePath, "empty.yaml"),
			wantErr: true,
		},
		{
			name:    "no spec",
			path:    filepath.Join(basePath, "no-spec.yaml"),
			wantErr: true,
		},
		{
			name:    "invalid testStep",
			path:    filepath.Join(basePath, "invalid-testStep.yaml"),
			wantErr: true,
		},
		{
			name: "ok",
			path: filepath.Join(basePath, "ok.yaml"),
			want: []*v1alpha1.TestStep{
				{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "chainsaw.kyverno.io/v1alpha1",
						Kind:       "TestStep",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-1",
					},
					Spec: v1alpha1.TestStepSpec{
						Apply: []v1alpha1.Apply{
							{
								File: "foo.yaml",
							},
						},
						Assert: []v1alpha1.Assert{
							{
								File: "bar.yaml",
							},
						},
					},
				},
			},
		},
		{
			name: "mlutiple testStep",
			path: filepath.Join(basePath, "multiple-testStep.yaml"),
			want: []*v1alpha1.TestStep{
				{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "chainsaw.kyverno.io/v1alpha1",
						Kind:       "TestStep",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-1",
					},
					Spec: v1alpha1.TestStepSpec{
						Apply: []v1alpha1.Apply{
							{
								File: "foo.yaml",
							},
						},
						Assert: []v1alpha1.Assert{
							{
								File: "bar.yaml",
							},
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "chainsaw.kyverno.io/v1alpha1",
						Kind:       "TestStep",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-2",
					},
					Spec: v1alpha1.TestStepSpec{
						Apply: []v1alpha1.Apply{
							{
								File: "bar.yaml",
							},
						},
						Assert: []v1alpha1.Assert{
							{
								File: "foo.yaml",
							},
						},
					},
				},
			},
		},
	}
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
