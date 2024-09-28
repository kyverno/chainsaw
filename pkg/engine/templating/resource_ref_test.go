package templating

import (
	"context"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestResourceRef(t *testing.T) {
	var noMeta unstructured.Unstructured
	noMeta.SetAPIVersion("v1")
	noMeta.SetKind("Namespace")
	var meta unstructured.Unstructured
	meta.SetAPIVersion("v1")
	meta.SetKind("Namespace")
	meta.SetName("foo")
	meta.SetNamespace("bar")
	meta.SetLabels(
		map[string]string{
			"($foo)": "bar",
		},
	)
	var binds unstructured.Unstructured
	binds.SetAPIVersion("v1")
	binds.SetKind("Namespace")
	binds.SetName("($foo)")
	binds.SetNamespace("($bar)")
	binds.SetLabels(
		map[string]string{
			"($foo)": "($bar)",
		},
	)
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		bindings apis.Bindings
		wantErr  bool
		want     *unstructured.Unstructured
	}{{
		name:     "nil",
		obj:      nil,
		bindings: binding.NewBindings(),
		wantErr:  false,
		want:     nil,
	}, {
		name:     "empty",
		obj:      &unstructured.Unstructured{},
		bindings: binding.NewBindings(),
		wantErr:  false,
		want:     &unstructured.Unstructured{},
	}, {
		name:     "meta",
		obj:      noMeta.DeepCopy(),
		bindings: binding.NewBindings(),
		wantErr:  false,
		want:     &noMeta,
	}, {
		name:     "no meta",
		obj:      meta.DeepCopy(),
		bindings: binding.NewBindings(),
		wantErr:  false,
		want:     &meta,
	}, {
		name:     "bindings",
		obj:      binds.DeepCopy(),
		bindings: binding.NewBindings().Register("$foo", binding.NewBinding("foo")).Register("$bar", binding.NewBinding("bar")),
		wantErr:  false,
		want:     &meta,
	}, {
		name:     "error",
		obj:      binds.DeepCopy(),
		bindings: binding.NewBindings(),
		wantErr:  true,
		want:     &binds,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ResourceRef(context.TODO(), tt.obj, tt.bindings)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, tt.obj)
		})
	}
}
