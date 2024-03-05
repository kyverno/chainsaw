package config

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateCluster(t *testing.T) {
	tests := []struct {
		name string
		path *field.Path
		obj  v1alpha1.Cluster
		want field.ErrorList
	}{{
		name: "empty",
		path: field.NewPath("foo"),
		want: field.ErrorList{
			&field.Error{
				Type:     field.ErrorTypeRequired,
				BadValue: "",
				Field:    "foo.kubeconfig",
				Detail:   "a kubeconfig is required",
			},
		},
	}, {
		name: "valid",
		path: field.NewPath("foo"),
		obj: v1alpha1.Cluster{
			Kubeconfig: "foo",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateCluster(tt.path, tt.obj)
			assert.Equal(t, tt.want, got)
		})
	}
}
