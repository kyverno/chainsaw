package kubectl

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestLogs(t *testing.T) {
	tests := []struct {
		name           string
		collector      *v1alpha1.PodLogs
		wantEntrypoint string
		wantArgs       []string
		wantErr        bool
	}{{
		name:      "nil",
		collector: nil,
		wantErr:   true,
	}, {
		name:      "empty",
		collector: &v1alpha1.PodLogs{},
		wantErr:   true,
	}, {
		name: "with name",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Name: "foo",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"logs", "--prefix", "foo", "-n", "$NAMESPACE", "--all-containers"},
		wantErr:        false,
	}, {
		name: "with namespace",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Namespace: "foo",
				},
			},
		},
		wantErr: true,
	}, {
		name: "with name and namespace",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Name:      "foo",
					Namespace: "bar",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"logs", "--prefix", "foo", "-n", "bar", "--all-containers"},
		wantErr:        false,
	}, {
		name: "with name and container",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Name: "foo",
				},
			},
			Container: "bar",
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"logs", "--prefix", "foo", "-n", "$NAMESPACE", "-c", "bar"},
		wantErr:        false,
	}, {
		name: "with name, namespace and container",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Name:      "foo",
					Namespace: "lorem",
				},
			},
			Container: "bar",
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"logs", "--prefix", "foo", "-n", "lorem", "-c", "bar"},
		wantErr:        false,
	}, {
		name: "with tail",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Name:      "foo",
					Namespace: "lorem",
				},
			},
			Container: "bar",
			Tail:      ptr.To(100),
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"logs", "--prefix", "foo", "-n", "lorem", "-c", "bar", "--tail", "100"},
		wantErr:        false,
	}, {
		name: "with selector",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				Selector: "foo=bar",
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"logs", "--prefix", "-l", "foo=bar", "-n", "$NAMESPACE", "--all-containers"},
		wantErr:        false,
	}, {
		name: "with name and selector",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Name: "foo",
				},
				Selector: "foo=bar",
			},
		},
		wantErr: true,
	}, {
		name: "with namespace and selector",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Namespace: "foo",
				},
				Selector: "foo=bar",
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"logs", "--prefix", "-l", "foo=bar", "-n", "foo", "--all-containers"},
		wantErr:        false,
	}, {
		name: "bad name",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Name: "($bad)",
				},
			},
		},
		wantErr: true,
	}, {
		name: "bad namespace",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Namespace: "($bad)",
				},
			},
		},
		wantErr: true,
	}, {
		name: "bad selector",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Name: "foo",
				},
				Selector: "($bad)",
			},
		},
		wantErr: true,
	}, {
		name: "bad container",
		collector: &v1alpha1.PodLogs{
			ActionObjectSelector: v1alpha1.ActionObjectSelector{
				ObjectName: v1alpha1.ObjectName{
					Name: "foo",
				},
			},
			Container: "($bad)",
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entrypoint, args, err := Logs(context.TODO(), apis.DefaultCompilers, nil, tt.collector)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantEntrypoint, entrypoint)
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}
