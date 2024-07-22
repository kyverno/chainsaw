package model

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
	"github.com/kyverno/pkg/ext/resource/convert"
	corev1 "k8s.io/api/core/v1"
)

func buildNamespace(ctx context.Context, name string, template *v1alpha1.Any, bindings binding.Bindings) (*corev1.Namespace, error) {
	namespace := client.Namespace(name)
	if template == nil {
		return &namespace, nil
	}
	if template.Value == nil {
		return &namespace, nil
	}
	object := client.ToUnstructured(&namespace)
	bindings = apibindings.RegisterNamedBinding(ctx, bindings, "namespace", object.GetName())
	merged, err := mutate.Merge(ctx, object, bindings, *template)
	if err != nil {
		return nil, err
	}
	return convert.To[corev1.Namespace](merged)
}
