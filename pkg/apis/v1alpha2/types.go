package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	_ "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
)

const (
	EngineJP  = v1alpha1.EngineJP
	EngineCEL = v1alpha1.EngineCEL
)

type (
	Clusters        = v1alpha1.Clusters
	Compiler        = v1alpha1.Compiler
	DefaultTimeouts = v1alpha1.DefaultTimeouts
	Projection      = v1alpha1.Projection
)
