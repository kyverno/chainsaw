package apis

import (
	"github.com/kyverno/chainsaw/pkg/engine/functions"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/cel"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/jp"
)

var (
	defaultCompilers = compilers.Compilers{
		Jp:  jp.NewCompiler(jp.WithFunctionCaller(functions.Caller())),
		Cel: cel.NewCompiler(),
	}
	DefaultCompilers = defaultCompilers.WithDefaultCompiler(compilers.CompilerJP)
)
