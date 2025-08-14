package apis

import (
	"sync"

	gocel "github.com/google/cel-go/cel"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/engine/functions"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/cel"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/jp"
	"k8s.io/apiserver/pkg/cel/library"
)

var (
	env = sync.OnceValues(func() (*gocel.Env, error) {
		env, err := cel.DefaultEnv()
		if err != nil {
			return nil, err
		}
		return env.Extend(
			library.URLs(),
			library.Regex(),
			library.Lists(),
			library.Authz(),
			library.Quantity(),
			library.IP(),
			library.CIDR(),
			library.Format(),
			library.AuthzSelectors(),
		)
	})
	defaultCompilers = compilers.Compilers{
		Jp:  jp.NewCompiler(jp.WithFunctionCaller(functions.Caller())),
		Cel: cel.NewCompiler(env),
	}
	DefaultCompilers = defaultCompilers.WithDefaultCompiler(compilers.CompilerJP)
)

type Bindings = binding.Bindings

var (
	NewBinding  = binding.NewBinding
	NewBindings = binding.NewBindings
)
