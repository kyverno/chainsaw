package functions

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func Test_jpKubernetesResourceExists(t *testing.T) {
	config, _ := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	client, _ := client.New(config)
	tests := []struct {
		name    string
		args    []any
		want    any
		wantErr bool
	}{{
		name: "pods",
		args: []any{
			client,
			"v1",
			"Pod",
		},
		want:    true,
		wantErr: false,
	}, {
		name: "deployments",
		args: []any{
			client,
			"apps/v1",
			"Deployment",
		},
		want:    true,
		wantErr: false,
	}, {
		name: "foos",
		args: []any{
			client,
			"v1",
			"Foo",
		},
		want:    false,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpKubernetesResourceExists(tt.args)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
