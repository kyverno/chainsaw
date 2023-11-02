package collect

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func Test_events(t *testing.T) {
	tests := []struct {
		name      string
		collector *v1alpha1.Events
		want      *v1alpha1.Command
		wantErr   bool
	}{{
		name:      "nil",
		collector: nil,
		want:      nil,
		wantErr:   false,
	}, {
		name:      "empty",
		collector: &v1alpha1.Events{},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "events", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with name",
		collector: &v1alpha1.Events{
			Name: "foo",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "events", "foo", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with namespace",
		collector: &v1alpha1.Events{
			Namespace: "foo",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "events", "-n", "foo"},
		},
		wantErr: false,
	}, {
		name: "with name and namespace",
		collector: &v1alpha1.Events{
			Name:      "bar",
			Namespace: "foo",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "events", "bar", "-n", "foo"},
		},
		wantErr: false,
	}, {
		name: "with selector",
		collector: &v1alpha1.Events{
			Selector: "foo=bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "events", "-l", "foo=bar", "-n", "$NAMESPACE"},
		},
		wantErr: false,
	}, {
		name: "with name and selector",
		collector: &v1alpha1.Events{
			Name:     "foo",
			Selector: "foo=bar",
		},
		want:    nil,
		wantErr: true,
	}, {
		name: "with namespace and selector",
		collector: &v1alpha1.Events{
			Namespace: "foo",
			Selector:  "foo=bar",
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"get", "events", "-l", "foo=bar", "-n", "foo"},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := events(tt.collector)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
