package metrics

import (
	"context"
	"testing"

	"fmt"

	"github.com/prometheus/common/model"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDecode(t *testing.T) {
	clcfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		panic(err.Error())
	}
	restcfg, err := clientcmd.NewNonInteractiveClientConfig(
		*clcfg, "", &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(restcfg)
	res := clientset.CoreV1().RESTClient().Get().
		Namespace("kyverno").
		Resource("pods").
		Name("kyverno-admission-controller-86664c8d5c-kgf2t:8000").
		SubResource("proxy").
		// The server URL path, without leading "/" goes here...
		Suffix("metrics").
		Do(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	rawbody, err := res.Raw()
	if err != nil {
		panic(err.Error())
	}
	fmt.Print(string(rawbody))
	tests := []struct {
		name    string
		in      string
		ts      model.Time
		wantErr bool
	}{
		{
			in: `
	# Only a quite simple scenario with two metric families.
	# More complicated tests of the parser itself can be found in the text package.
	# TYPE mf2 counter
	mf2 3
	mf1{label="value1"} -3.14 123456
	mf1{label="value2"} 42
	mf2 4
	`,
			ts:      model.Now(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Decode(tt.in, tt.ts); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// kubectl get --raw /api/v1/namespaces/kyverno/services/kyverno-svc-metrics:metrics-port/proxy/metrics
