package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	_ "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
)

type (
	Any             = v1alpha1.Any
	Binding         = v1alpha1.Binding
	Catch           = v1alpha1.CatchFinally
	Check           = v1alpha1.Check
	Cluster         = v1alpha1.Cluster
	Clusters        = v1alpha1.Clusters
	Expectation     = v1alpha1.Expectation
	Format          = v1alpha1.Format
	ObjectName      = v1alpha1.ObjectName
	ObjectReference = v1alpha1.ObjectReference
	ObjectType      = v1alpha1.ObjectType
	TestStepSpec    = v1alpha1.TestStepSpec
	Timeouts        = v1alpha1.Timeouts
)
