package config

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateConfiguration(t *testing.T) {
	tests := []struct {
		name string
		obj  *v1alpha1.Configuration
		want field.ErrorList
	}{{
		name: "null",
	}, {
		name: "empty",
		obj:  &v1alpha1.Configuration{},
	}, {
		name: "with cluster",
		obj: &v1alpha1.Configuration{
			Spec: v1alpha1.ConfigurationSpec{
				Clusters: map[string]v1alpha1.Cluster{
					"foo": {
						Kubeconfig: "foo",
					},
				},
			},
		},
	}, {
		name: "with catch",
		obj: &v1alpha1.Configuration{
			Spec: v1alpha1.ConfigurationSpec{
				Catch: []v1alpha1.Catch{{
					Events: &v1alpha1.Events{},
				}},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateConfiguration(tt.obj)
			assert.Equal(t, tt.want, got)
		})
	}
}
