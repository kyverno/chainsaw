package functions

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func Test_jpKubernetesResourceExists(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := client.New(config)
	assert.NoError(t, err)
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

func Test_jpKubernetesExists(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := client.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name    string
		args    []any
		want    any
		wantErr bool
	}{{
		name: "exist",
		args: []any{
			client,
			"v1",
			"Pod",
			"kube-system",
			"kube-apiserver-kind-control-plane",
		},
		want:    true,
		wantErr: false,
	}, {
		name: "namespace not exists",
		args: []any{
			client,
			"v1",
			"Pod",
			"foo",
			"kube-apiserver-kind-control-plane",
		},
		want:    false,
		wantErr: false,
	}, {
		name: "name not exists",
		args: []any{
			client,
			"v1",
			"Pod",
			"kube-system",
			"foo",
		},
		want:    false,
		wantErr: false,
	}, {
		name: "kind not exists",
		args: []any{
			client,
			"v1",
			"Foo",
			"kube-system",
			"foo",
		},
		want:    nil,
		wantErr: true,
	}, {
		name: "apiVersion not exists",
		args: []any{
			client,
			"v1alpha1",
			"Pod",
			"kube-system",
			"kube-apiserver-kind-control-plane",
		},
		want:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpKubernetesExists(tt.args)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
