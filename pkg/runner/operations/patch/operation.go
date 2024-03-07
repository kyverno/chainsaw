package patch

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	mutation "github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type operation struct {
	client     client.Client
	base       unstructured.Unstructured
	namespacer namespacer.Namespacer
	template   bool
	expect     []v1alpha1.Expectation
	outputs    []v1alpha1.Output
}

func New(
	client client.Client,
	obj unstructured.Unstructured,
	namespacer namespacer.Namespacer,
	template bool,
	expect []v1alpha1.Expectation,
	outputs []v1alpha1.Output,
) operations.Operation {
	return &operation{
		client:     client,
		base:       obj,
		namespacer: namespacer,
		template:   template,
		expect:     expect,
		outputs:    outputs,
	}
}

func (o *operation) Exec(ctx context.Context, bindings binding.Bindings) (_ operations.Outputs, _err error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	obj := o.base
	logger := internal.GetLogger(ctx, &obj)
	defer func() {
		internal.LogEnd(logger, logging.Patch, _err)
	}()
	if o.template {
		template := v1alpha1.Any{
			Value: obj.UnstructuredContent(),
		}
		if merged, err := mutate.Merge(ctx, obj, bindings, template); err != nil {
			return nil, err
		} else {
			obj = merged
		}
	}
	if err := internal.ApplyNamespacer(o.namespacer, &obj); err != nil {
		return nil, err
	}
	internal.LogStart(logger, logging.Patch)
	return o.execute(ctx, bindings, obj)
}

func (o *operation) execute(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) (operations.Outputs, error) {
	var lastErr error
	var outputs operations.Outputs
	err := wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		outputs, lastErr = o.tryPatchResource(ctx, bindings, obj)
		// TODO: determine if the error can be retried
		return lastErr == nil, nil
	})
	if err == nil {
		return outputs, nil
	}
	if lastErr != nil {
		return outputs, lastErr
	}
	return outputs, err
}

func (o *operation) tryPatchResource(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) (operations.Outputs, error) {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
	err := o.client.Get(ctx, client.ObjectKey(&obj), &actual)
	if err != nil {
		return nil, err
	}
	return o.updateResource(ctx, bindings, &actual, obj)
}

func (o *operation) updateResource(ctx context.Context, bindings binding.Bindings, actual *unstructured.Unstructured, obj unstructured.Unstructured) (operations.Outputs, error) {
	patched, err := client.PatchObject(actual, &obj)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(patched)
	if err != nil {
		return nil, err
	}
	return o.handleCheck(ctx, bindings, obj, o.client.Patch(ctx, actual, ctrlclient.RawPatch(types.MergePatchType, bytes)))
}

func (o *operation) handleCheck(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured, err error) (_outputs operations.Outputs, _err error) {
	defer func() {
		var outputs operations.Outputs
		if _err == nil {
			for _, output := range o.outputs {
				if err := output.CheckName(); err != nil {
					_err = err
					return
				}
				if output.Match != nil && output.Match.Value != nil {
					if errs, err := check.Check(ctx, obj.UnstructuredContent(), nil, output.Match); err != nil {
						_err = err
						return
					} else if len(errs) != 0 {
						continue
					}
				}
				patched, err := mutation.Mutate(ctx, nil, mutation.Parse(ctx, output.Value.Value), obj.UnstructuredContent(), bindings, template.WithFunctionCaller(functions.Caller))
				if err != nil {
					_err = err
					return
				}
				if outputs == nil {
					outputs = operations.Outputs{}
				}
				outputs[output.Name] = binding.NewBinding(patched)
				bindings = bindings.Register("$"+output.Name, outputs[output.Name])
			}
			_outputs = outputs
		}
	}()
	if err == nil {
		bindings = bindings.Register("$error", binding.NewBinding(nil))
	} else {
		bindings = bindings.Register("$error", binding.NewBinding(err.Error()))
	}
	if matched, err := check.Expectations(ctx, obj, bindings, o.expect...); matched {
		return nil, err
	}
	return nil, err
}
