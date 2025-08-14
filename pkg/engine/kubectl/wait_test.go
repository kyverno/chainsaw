package kubectl

import (
	"context"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client/simple"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
)

func TestWait(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := simple.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name           string
		collector      *v1alpha1.Wait
		wantEntrypoint string
		wantArgs       []string
		wantErr        bool
	}{{
		name:      "nil",
		collector: nil,
		wantErr:   true,
	}, {
		name:      "empty",
		collector: &v1alpha1.Wait{},
		wantErr:   true,
	}, {
		name: "valid resource and condition",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
			WaitFor: v1alpha1.WaitFor{
				Condition: &v1alpha1.WaitForCondition{
					Name: "Ready",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", "--for=condition=Ready", "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		wantErr:        false,
	}, {
		name: "valid clustered resource and condition",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "rbac.authorization.k8s.io/v1",
					Kind:       "ClusterRole",
				},
			},
			WaitFor: v1alpha1.WaitFor{
				Condition: &v1alpha1.WaitForCondition{
					Name: "Ready",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "clusterroles.v1.rbac.authorization.k8s.io", "--for=condition=Ready", "--all", "--timeout=-1s"},
		wantErr:        false,
	}, {
		name: "valid resource and condition with value",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
			WaitFor: v1alpha1.WaitFor{
				Condition: &v1alpha1.WaitForCondition{
					Name:  "Ready",
					Value: ptr.To(v1alpha1.Expression("test")),
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", `--for=condition=Ready=test`, "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		wantErr:        false,
	}, {
		name: "valid resource and condition with empty value",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
			WaitFor: v1alpha1.WaitFor{
				Condition: &v1alpha1.WaitForCondition{
					Name:  "Ready",
					Value: ptr.To(v1alpha1.Expression("")),
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", `--for=condition=Ready=`, "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		wantErr:        false,
	}, {
		name: "valid resource and delete",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
			WaitFor: v1alpha1.WaitFor{
				Deletion: &v1alpha1.WaitForDeletion{},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", "--for=delete", "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		wantErr:        false,
	}, {
		name: "valid resource and jsonpath without value",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
			WaitFor: v1alpha1.WaitFor{
				JsonPath: &v1alpha1.WaitForJsonPath{
					Path: "{.status.phase}",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", "--for=jsonpath={.status.phase}", "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		wantErr:        false,
	}, {
		name: "valid resource and jsonpath with value",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
			WaitFor: v1alpha1.WaitFor{
				JsonPath: &v1alpha1.WaitForJsonPath{
					Path:  "{.status.phase}",
					Value: ptr.To(v1alpha1.Expression("Running")),
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", "--for=jsonpath={.status.phase}=Running", "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		wantErr:        false,
	}, {
		name: "with resource name",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "my-pod",
					},
				},
			},
			WaitFor: v1alpha1.WaitFor{
				Condition: &v1alpha1.WaitForCondition{
					Name: "Ready",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", "--for=condition=Ready", "my-pod", "-n", "$NAMESPACE", "--timeout=-1s"},
		wantErr:        false,
	}, {
		name: "with selector",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					Selector: "app=my-app",
				},
			},
			WaitFor: v1alpha1.WaitFor{
				Condition: &v1alpha1.WaitForCondition{
					Name: "Ready",
				},
			},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", "--for=condition=Ready", "-l", "app=my-app", "-n", "$NAMESPACE", "--timeout=-1s"},
		wantErr:        false,
	}, {
		name: "with timeout",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
			WaitFor: v1alpha1.WaitFor{
				Condition: &v1alpha1.WaitForCondition{
					Name: "Ready",
				},
			},
			ActionTimeout: v1alpha1.ActionTimeout{Timeout: &metav1.Duration{Duration: 120 * time.Second}},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", "--for=condition=Ready", "--all", "-n", "$NAMESPACE", "--timeout", "2m0s"},
		wantErr:        false,
	}, {
		name: "name and selector error",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					Selector: "app=my-app",
					ObjectName: v1alpha1.ObjectName{
						Name: "my-pod",
					},
				},
			},
			WaitFor: v1alpha1.WaitFor{
				Condition: &v1alpha1.WaitForCondition{
					Name: "Ready",
				},
			},
		},
		wantErr: true,
	}, {
		name: "missing condition",
		collector: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "my-pod",
					},
				},
			},
		},
		wantErr: true,
	}, {
		name: "with all namespaces",
		collector: &v1alpha1.Wait{
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
			WaitFor: v1alpha1.WaitFor{
				Condition: &v1alpha1.WaitForCondition{
					Name: "Ready",
				},
			},
			ActionTimeout: v1alpha1.ActionTimeout{Timeout: &metav1.Duration{Duration: 120 * time.Second}},
		},
		wantEntrypoint: "kubectl",
		wantArgs:       []string{"wait", "pods", "--for=condition=Ready", "--all", "--all-namespaces", "--timeout", "2m0s"},
		wantErr:        false,
	}, {
		name: "bad name",
		collector: &v1alpha1.Wait{
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
		collector: &v1alpha1.Wait{
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
		collector: &v1alpha1.Wait{
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
		collector: &v1alpha1.Wait{
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
			entrypoint, args, err := Wait(context.TODO(), apis.DefaultCompilers, client, nil, tt.collector)
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
