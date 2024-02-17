package collect

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		collector *v1alpha1.Get
		want      *v1alpha1.Command
		wantErr   bool
	}{{
		name:      "nil",
		collector: nil,
		want:      nil,
		wantErr:   false,
	}, {
		name:      "empty",
		collector: &v1alpha1.Get{},
		wantErr:   true,
	}, {
		name: "without resource",
		collector: &v1alpha1.Get{
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name: "foo",
			},
		},
		wantErr: true,
	}, {
		name: "with resource",
		collector: &v1alpha1.Get{
			Resource: "foos",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "foos", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with name",
		collector: &v1alpha1.Get{
			Resource: "foos",
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name: "foo",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "foos", "foo", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with namespace",
		collector: &v1alpha1.Get{
			Resource: "foos",
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Namespace: "bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "foos", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with name and namespace",
		collector: &v1alpha1.Get{
			Resource: "foos",
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "foos", "foo", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with selector",
		collector: &v1alpha1.Get{
			Resource: "foos",
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Selector: "foo=bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "foos", "-l", "foo=bar", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with name and selector",
		collector: &v1alpha1.Get{
			Resource: "foos",
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name:     "foo",
				Selector: "foo=bar",
			},
		},
		wantErr: true,
	}, {
		name: "with namespace and selector",
		collector: &v1alpha1.Get{
			Resource: "foos",
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Namespace: "bar",
				Selector:  "foo=bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "foos", "-l", "foo=bar", "-n", "bar"},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.collector)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
