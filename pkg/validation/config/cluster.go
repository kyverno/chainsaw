package config

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateCluster(path *field.Path, obj v1alpha1.Cluster) field.ErrorList {
	var errs field.ErrorList
	if obj.Kubeconfig == "" {
		errs = append(errs, field.Required(path.Child("kubeconfig"), "a kubeconfig is required"))
	}
	return errs
}
