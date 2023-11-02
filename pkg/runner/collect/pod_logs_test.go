package collect

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func Test_podLogs(t *testing.T) {
	tests := []struct {
		name      string
		collector *v1alpha1.PodLogs
		want      *v1alpha1.Command
		wantErr   bool
	}{{
		name:      "nil",
		collector: nil,
		want:      nil,
		wantErr:   false,
	}, {
		name:      "empty",
		collector: &v1alpha1.PodLogs{},
		wantErr:   true,
	}, {
		name: "with name",
		collector: &v1alpha1.PodLogs{
			Name: "foo",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"logs", "--prefix", "foo", "-n", "$NAMESPACE", "--all-containers"},
		},
		wantErr: false,
	}, {
		name: "with namespace",
		collector: &v1alpha1.PodLogs{
			Namespace: "foo",
		},
		wantErr: true,
	}, {
		name: "with name and namespace",
		collector: &v1alpha1.PodLogs{
			Name:      "foo",
			Namespace: "bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"logs", "--prefix", "foo", "-n", "bar", "--all-containers"},
		},
		wantErr: false,
	}, {
		name: "with name and container",
		collector: &v1alpha1.PodLogs{
			Name:      "foo",
			Container: "bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"logs", "--prefix", "foo", "-n", "$NAMESPACE", "-c", "bar"},
		},
		wantErr: false,
	}, {
		name: "with name, namespace and container",
		collector: &v1alpha1.PodLogs{
			Name:      "foo",
			Namespace: "lorem",
			Container: "bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"logs", "--prefix", "foo", "-n", "lorem", "-c", "bar"},
		},
		wantErr: false,
	}, {
		name: "with tail",
		collector: &v1alpha1.PodLogs{
			Name:      "foo",
			Namespace: "lorem",
			Container: "bar",
			Tail:      ptr.To(100),
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"logs", "--prefix", "foo", "-n", "lorem", "-c", "bar", "--tail", "100"},
		},
		wantErr: false,
	}, {
		name: "with selector",
		collector: &v1alpha1.PodLogs{
			Selector: "foo=bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"logs", "--prefix", "-l", "foo=bar", "-n", "$NAMESPACE", "--all-containers"},
		},
		wantErr: false,
	}, {
		name: "with name and selector",
		collector: &v1alpha1.PodLogs{
			Name:     "foo",
			Selector: "foo=bar",
		},
		want:    nil,
		wantErr: true,
	}, {
		name: "with namespace and selector",
		collector: &v1alpha1.PodLogs{
			Namespace: "foo",
			Selector:  "foo=bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"logs", "--prefix", "-l", "foo=bar", "-n", "foo", "--all-containers"},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := podLogs(tt.collector)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
