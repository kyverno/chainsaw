package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	_ "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
)

type (
	Any             = v1alpha1.Any
	Clusters        = v1alpha1.Clusters
	DefaultTimeouts = v1alpha1.DefaultTimeouts
)
