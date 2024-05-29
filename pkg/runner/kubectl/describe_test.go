package kubectl

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
)

func TestDescribe(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := client.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name      string
		collector *v1alpha1.Describe
		want      *v1alpha1.Command
		wantErr   bool
	}{{
		name:      "nil",
		collector: nil,
		want:      nil,
		wantErr:   true,
	}, {
		name:      "empty",
		collector: &v1alpha1.Describe{},
		wantErr:   true,
	}, {
		name: "without resource",
		collector: &v1alpha1.Describe{
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name: "foo",
			},
		},
		wantErr: true,
	}, {
		name: "with resource",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "v1",
				Kind:       "Pod",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "pods", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with clustered resource",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "rbac.authorization.k8s.io/v1",
				Kind:       "ClusterRole",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "clusterroles.v1.rbac.authorization.k8s.io"},
		},
		wantErr: false,
	}, {
		name: "with name",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name: "foo",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "pods", "foo", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with namespace",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Namespace: "bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "pods", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with name and namespace",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "pods", "foo", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with selector",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Selector: "foo=bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "pods", "-l", "foo=bar", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with name and selector",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name:     "foo",
				Selector: "foo=bar",
			},
		},
		wantErr: true,
	}, {
		name: "with namespace and selector",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Namespace: "bar",
				Selector:  "foo=bar",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "pods", "-l", "foo=bar", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with show-events marked as true",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ShowEvents: ptr.To(true),
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "pods", "-n", "$NAMESPACE", "--show-events=true"},
		},
		wantErr: false,
	}, {
		name: "with show-events marked as false",
		collector: &v1alpha1.Describe{
			ResourceReference: v1alpha1.ResourceReference{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ShowEvents: ptr.To(false),
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "pods", "-n", "$NAMESPACE", "--show-events=false"},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Describe(client, nil, tt.collector)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
