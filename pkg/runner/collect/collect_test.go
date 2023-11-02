package collect

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestCommands(t *testing.T) {
	tests := []struct {
		name      string
		collector *v1alpha1.Collect
		want      []*v1alpha1.Command
		wantErr   bool
	}{{
		name:      "nil",
		collector: nil,
		want:      nil,
		wantErr:   false,
	}, {
		name:      "empty",
		collector: &v1alpha1.Collect{},
		want:      nil,
		wantErr:   false,
	}, {
		name: "with pod logs",
		collector: &v1alpha1.Collect{
			PodLogs: &v1alpha1.PodLogs{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: []*v1alpha1.Command{{
			Entrypoint: "kubectl",
			Args:       []string{"logs", "--prefix", "foo", "-n", "bar", "--all-containers"},
		}},
		wantErr: false,
	}, {
		name: "with events",
		collector: &v1alpha1.Collect{
			Events: &v1alpha1.Events{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: []*v1alpha1.Command{{
			Entrypoint: "kubectl",
			Args:       []string{"get", "events", "foo", "-n", "bar"},
		}},
		wantErr: false,
	}, {
		name: "with pod logs and events",
		collector: &v1alpha1.Collect{
			PodLogs: &v1alpha1.PodLogs{
				Name:      "foo",
				Namespace: "bar",
			},
			Events: &v1alpha1.Events{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: []*v1alpha1.Command{{
			Entrypoint: "kubectl",
			Args:       []string{"logs", "--prefix", "foo", "-n", "bar", "--all-containers"},
		}, {
			Entrypoint: "kubectl",
			Args:       []string{"get", "events", "foo", "-n", "bar"},
		}},
		wantErr: false,
	}, {
		name: "with error",
		collector: &v1alpha1.Collect{
			PodLogs: &v1alpha1.PodLogs{
				Name:      "foo",
				Namespace: "bar",
			},
			Events: &v1alpha1.Events{
				Name:     "foo",
				Selector: "foo=bar",
			},
		},
		wantErr: true,
	}, {
		name: "with error",
		collector: &v1alpha1.Collect{
			PodLogs: &v1alpha1.PodLogs{
				Name:     "foo",
				Selector: "foo=bar",
			},
			Events: &v1alpha1.Events{
				Name:      "foo",
				Selector:  "foo=bar",
				Namespace: "bar",
			},
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Commands(tt.collector)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
