package kubectl

import (
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
)

func TestWait(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := client.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name    string
		waiter  *v1alpha1.Wait
		want    *v1alpha1.Command
		wantErr bool
	}{{
		name:    "nil waiter",
		waiter:  nil,
		wantErr: true,
	}, {
		name:    "empty waiter",
		waiter:  &v1alpha1.Wait{},
		wantErr: true,
	}, {
		name: "valid resource and condition",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			For: v1alpha1.For{
				Condition: &v1alpha1.Condition{
					Name: "Ready",
				},
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"wait", "pods", "--for=condition=Ready", "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		},
		wantErr: false,
	}, {
		name: "valid clustered resource and condition",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "clusterroles.v1.rbac.authorization.k8s.io",
			},
			For: v1alpha1.For{
				Condition: &v1alpha1.Condition{
					Name: "Ready",
				},
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"wait", "clusterroles.v1.rbac.authorization.k8s.io", "--for=condition=Ready", "--all", "--timeout=-1s"},
		},
		wantErr: false,
	}, {
		name: "valid resource and condition with value",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			For: v1alpha1.For{
				Condition: &v1alpha1.Condition{
					Name:  "Ready",
					Value: ptr.To("test"),
				},
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"wait", "pods", `--for=condition=Ready=test`, "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		},
		wantErr: false,
	}, {
		name: "valid resource and condition with empty value",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			For: v1alpha1.For{
				Condition: &v1alpha1.Condition{
					Name:  "Ready",
					Value: ptr.To(""),
				},
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"wait", "pods", `--for=condition=Ready=`, "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		},
		wantErr: false,
	}, {
		name: "valid resource and delete",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			For: v1alpha1.For{
				Deletion: &v1alpha1.Deletion{},
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"wait", "pods", "--for=delete", "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		},
		wantErr: false,
	}, {
		name: "valid resource and jsonpath",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			For: v1alpha1.For{
				JsonPath: &v1alpha1.JsonPath{
					Path:  "{.status.phase}",
					Value: "Running",
				},
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"wait", "pods", "--for=jsonpath='{.status.phase}'=Running", "--all", "-n", "$NAMESPACE", "--timeout=-1s"},
		},
		wantErr: false,
	}, {
		name: "with resource name",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name: "my-pod",
			},
			For: v1alpha1.For{
				Condition: &v1alpha1.Condition{
					Name: "Ready",
				},
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"wait", "pods", "--for=condition=Ready", "my-pod", "-n", "$NAMESPACE", "--timeout=-1s"},
		},
		wantErr: false,
	}, {
		name: "with selector",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			For: v1alpha1.For{
				Condition: &v1alpha1.Condition{
					Name: "Ready",
				},
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Selector: "app=my-app",
			},
		},
		want: &v1alpha1.Command{
			Entrypoint: "kubectl",
			Args:       []string{"wait", "pods", "--for=condition=Ready", "-l", "app=my-app", "-n", "$NAMESPACE", "--timeout=-1s"},
		},
		wantErr: false,
	}, {
		name: "with timeout",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			For: v1alpha1.For{
				Condition: &v1alpha1.Condition{
					Name: "Ready",
				},
			},
			Timeout: &metav1.Duration{Duration: 120 * time.Second},
		},
		want: &v1alpha1.Command{
			Timeout:    &metav1.Duration{Duration: 120 * time.Second},
			Entrypoint: "kubectl",
			Args:       []string{"wait", "pods", "--for=condition=Ready", "--all", "-n", "$NAMESPACE", "--timeout", "2m0s"},
		},
		wantErr: false,
	}, {
		name: "name and selector error",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Selector: "app=my-app",
				Name:     "my-pod",
			},
			For: v1alpha1.For{
				Condition: &v1alpha1.Condition{
					Name: "Ready",
				},
			},
		},
		wantErr: true,
	}, {
		name: "missing condition",
		waiter: &v1alpha1.Wait{
			ResourceReference: v1alpha1.ResourceReference{
				Resource: "pods",
			},
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name: "my-pod",
			},
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Wait(client, nil, tt.waiter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
