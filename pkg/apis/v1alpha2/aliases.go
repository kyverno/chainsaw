package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	_ "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
)

type (
	// Any represents any type.
	Any = v1alpha1.Any
	// Catch defines actions to be executed on failure.
	Catch = v1alpha1.Catch
	// Cluster defines cluster config and context.
	Cluster = v1alpha1.Cluster
	// Timeouts contains timeouts per operation.
	Timeouts = v1alpha1.Timeouts
)
