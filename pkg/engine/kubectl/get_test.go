package kubectl

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client/simple"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGet(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := simple.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name           string
		collector      *v1alpha1.Get
		wantEntrypoint string
		wantArgs       []string
		wantErr        bool
	}{{
		name:      "nil",
		collector: nil,
		wantErr:   true,
	}, {
		name:      "empty",
		collector: &v1alpha1.Get{},
		wantErr:   true,
	}, {
		name: "without resource",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "foo",
					},
				},
			},
		},
		wantErr: true,
	}, {
		name: "with resource",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"get", "pods", "-n", "$NAMESPACE"},
		wantErr:        false,
	}, {
		name: "with clustered resource",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "rbac.authorization.k8s.io/v1",
					Kind:       "ClusterRole",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"get", "clusterroles.v1.rbac.authorization.k8s.io"},
		wantErr:        false,
	}, {
		name: "with name",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "foo",
					},
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"get", "pods", "foo", "-n", "$NAMESPACE"},
		wantErr:        false,
	}, {
		name: "with namespace",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Namespace: "bar",
					},
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"get", "pods", "-n", "bar"},
		wantErr:        false,
	}, {
		name: "with name and namespace",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name:      "foo",
						Namespace: "bar",
					},
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"get", "pods", "foo", "-n", "bar"},
		wantErr:        false,
	}, {
		name: "with selector",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					Selector: "foo=bar",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"get", "pods", "-l", "foo=bar", "-n", "$NAMESPACE"},
		wantErr:        false,
	}, {
		name: "with name and selector",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "foo",
					},
					Selector: "foo=bar",
				},
			},
		},
		wantErr: true,
	}, {
		name: "with namespace and selector",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Namespace: "bar",
					},
					Selector: "foo=bar",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"get", "pods", "-l", "foo=bar", "-n", "bar"},
		wantErr:        false,
	}, {
		name: "with all namespaces",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Namespace: "*",
					},
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"get", "pods", "--all-namespaces"},
		wantErr:        false,
	}, {
		name: "bad name",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "($bad)",
					},
				},
			},
		},
		wantErr: true,
	}, {
		name: "bad namespace",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Namespace: "($bad)",
					},
				},
			},
		},
		wantErr: true,
	}, {
		name: "bad selector",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					Selector: "($bad)",
				},
			},
		},
		wantErr: true,
	}, {
		name: "bad format",
		collector: &v1alpha1.Get{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
			ActionFormat: v1alpha1.ActionFormat{
				Format: "($bad)",
			},
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entrypoint, args, err := Get(context.TODO(), apis.XDefaultCompilers, client, nil, tt.collector)
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
