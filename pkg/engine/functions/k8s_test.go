package functions

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/client/simple"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Test_jpKubernetesResourceExists(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := simple.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name    string
		args    []any
		want    any
		wantErr bool
	}{{
		name:    "nil",
		wantErr: true,
	}, {
		name:    "not enough args",
		args:    []any{&tclient.FakeClient{}, nil},
		wantErr: true,
	}, {
		name:    "not enough args",
		args:    []any{&tclient.FakeClient{}, "", nil},
		wantErr: true,
	}, {
		name:    "bad group",
		args:    []any{client, "foo/v1/bar", "Baz"},
		wantErr: true,
	}, {
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
	client, err := simple.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name    string
		args    []any
		want    any
		wantErr bool
	}{{
		name:    "nil",
		wantErr: true,
	}, {
		name:    "not enough args",
		args:    []any{&tclient.FakeClient{}, nil},
		wantErr: true,
	}, {
		name:    "not enough args",
		args:    []any{&tclient.FakeClient{}, "", nil},
		wantErr: true,
	}, {
		name:    "not enough args",
		args:    []any{&tclient.FakeClient{}, "", "", nil},
		wantErr: true,
	}, {
		name:    "not enough args",
		args:    []any{&tclient.FakeClient{}, "", "", "", nil},
		wantErr: true,
	}, {
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

func Test_jpKubernetesGet(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := simple.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name      string
		arguments []any
		want      any
		wantErr   bool
	}{{
		name:    "nil",
		wantErr: true,
	}, {
		name:      "not enough args",
		arguments: []any{&tclient.FakeClient{}, nil},
		wantErr:   true,
	}, {
		name:      "not enough args",
		arguments: []any{&tclient.FakeClient{}, "", nil},
		wantErr:   true,
	}, {
		name:      "not enough args",
		arguments: []any{&tclient.FakeClient{}, "", "", nil},
		wantErr:   true,
	}, {
		name:      "not enough args",
		arguments: []any{&tclient.FakeClient{}, "", "", "", nil},
		wantErr:   true,
	}, {
		name:      "kube-apiserver-kind-control-plane",
		arguments: []any{client, "v1", "Pod", "kube-system", "kube-apiserver-kind-control-plane"},
		wantErr:   false,
	}, {
		name:      "foo",
		arguments: []any{client, "v1", "Foo", "kube-system", "kube-apiserver-kind-control-plane"},
		wantErr:   true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpKubernetesGet(tt.arguments)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_jpKubernetesList(t *testing.T) {
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	client, err := simple.New(config)
	assert.NoError(t, err)
	tests := []struct {
		name      string
		arguments []any
		want      any
		wantErr   bool
	}{{
		name:    "nil",
		wantErr: true,
	}, {
		name:      "not enough args",
		arguments: []any{&tclient.FakeClient{}, nil},
		wantErr:   true,
	}, {
		name:      "not enough args",
		arguments: []any{&tclient.FakeClient{}, "", nil},
		wantErr:   true,
	}, {
		name:      "not enough args",
		arguments: []any{&tclient.FakeClient{}, "", "", nil},
		wantErr:   true,
	}, {
		name:      "nodes",
		arguments: []any{client, "v1", "Node"},
		wantErr:   false,
	}, {
		name:      "pods",
		arguments: []any{client, "v1", "Pod", "kube-system"},
		wantErr:   false,
	}, {
		name:      "foos",
		arguments: []any{client, "v1", "Foo", "kube-system"},
		wantErr:   true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpKubernetesList(tt.arguments)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_jpKubernetesServerVersion(t *testing.T) {
	var nilConfig *rest.Config
	config, err := restutils.DefaultConfig(clientcmd.ConfigOverrides{})
	assert.NoError(t, err)
	tests := []struct {
		name      string
		arguments []any
		wantErr   bool
	}{{
		name:    "nil",
		wantErr: true,
	}, {
		name:      "nil config",
		arguments: []any{nilConfig},
		wantErr:   true,
	}, {
		name:      "config",
		arguments: []any{config},
		wantErr:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpKubernetesServerVersion(tt.arguments)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
