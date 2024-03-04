package config

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/validation/test"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateConfigurationSpec(path *field.Path, obj v1alpha1.ConfigurationSpec) field.ErrorList {
	var errs field.ErrorList
	path = path.Child("clusters")
	for name, cluster := range obj.Clusters {
		errs = append(errs, ValidateCluster(path.Key(name), cluster)...)
	}
	for i, catch := range obj.Catch {
		errs = append(errs, test.ValidateCatch(path.Child("catch").Index(i), catch)...)
	}
	return errs
}
