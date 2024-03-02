package kubectl

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGet(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := client.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name      string
		collector *v1alpha1.Get
		want      *v1alpha1.Command
		wantErr   bool
	}{{
		name:      "nil",
		collector: nil,
		want:      nil,
		wantErr:   true,
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
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "pods", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with clustered resource",
		collector: &v1alpha1.Get{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "clusterroles.v1.rbac.authorization.k8s.io",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "clusterroles.v1.rbac.authorization.k8s.io"},
		},
		wantErr: false,
	}, {
		name: "with name",
		collector: &v1alpha1.Get{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name: "foo",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "pods", "foo", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with namespace",
		collector: &v1alpha1.Get{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Namespace: "bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "pods", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with name and namespace",
		collector: &v1alpha1.Get{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "pods", "foo", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with selector",
		collector: &v1alpha1.Get{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Selector: "foo=bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "pods", "-l", "foo=bar", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with name and selector",
		collector: &v1alpha1.Get{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name:     "foo",
				Selector: "foo=bar",
			},
		},
		wantErr: true,
	}, {
		name: "with namespace and selector",
		collector: &v1alpha1.Get{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Namespace: "bar",
				Selector:  "foo=bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "pods", "-l", "foo=bar", "-n", "bar"},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(client, tt.collector)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
