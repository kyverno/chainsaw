package kubectl

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestDescribe(t *testing.T) {
	tests := []struct {
		name      string
		collector *v1alpha1.Describe
		want      *v1alpha1.Command
		wantErr   bool
	}{{
		name:      "nil",
		collector: nil,
		want:      nil,
		wantErr:   false,
	}, {
		name:      "empty",
		collector: &v1alpha1.Describe{},
		wantErr:   true,
	}, {
		name: "without resource",
		collector: &v1alpha1.Describe{
			Name: "foo",
		},
		wantErr: true,
	}, {
		name: "with resource",
		collector: &v1alpha1.Describe{
			Resource: "foos",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "foos", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with name",
		collector: &v1alpha1.Describe{
			Resource: "foos",
			Name:     "foo",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "foos", "foo", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with namespace",
		collector: &v1alpha1.Describe{
			Resource:  "foos",
			Namespace: "bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "foos", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with name and namespace",
		collector: &v1alpha1.Describe{
			Resource:  "foos",
			Name:      "foo",
			Namespace: "bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "foos", "foo", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with selector",
		collector: &v1alpha1.Describe{
			Resource: "foos",
			Selector: "foo=bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "foos", "-l", "foo=bar", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with name and selector",
		collector: &v1alpha1.Describe{
			Resource: "foos",
			Name:     "foo",
			Selector: "foo=bar",
		},
		wantErr: true,
	}, {
		name: "with namespace and selector",
		collector: &v1alpha1.Describe{
			Resource:  "foos",
			Namespace: "bar",
			Selector:  "foo=bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "foos", "-l", "foo=bar", "-n", "bar"},
		},
		wantErr: false,
	}, {
		name: "with show-events marked as true",
		collector: &v1alpha1.Describe{
			Resource:   "foos",
			ShowEvents: ptr.To[bool](true),
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "foos", "-n", "$NAMESPACE", "--show-events=true"},
		},
		wantErr: false,
	}, {
		name: "with show-events marked as false",
		collector: &v1alpha1.Describe{
			Resource:   "foos",
			ShowEvents: ptr.To[bool](false),
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"describe", "foos", "-n", "$NAMESPACE", "--show-events=false"},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Describe(tt.collector)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
