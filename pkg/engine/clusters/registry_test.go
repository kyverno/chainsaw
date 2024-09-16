package clusters

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
)

func TestNewRegistry(t *testing.T) {
	tests := []struct {
		name string
		f    clientFactory
		want Registry
	}{{
		name: "nil",
		want: registry{
			clientFactory: nil,
			clusters:      map[string]Cluster{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRegistry(tt.f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_registry_Lookup(t *testing.T) {
	foo := NewClusterFromKubeconfig("foo", "bar")
	tests := []struct {
		name          string
		clientFactory clientFactory
		clusters      map[string]Cluster
		clusterName   string
		want          Cluster
	}{{
		name: "nil",
		want: nil,
	}, {
		name: "not found",
		clusters: map[string]Cluster{
			"foo": foo,
		},
		clusterName: "bar",
		want:        nil,
	}, {
		name: "found",
		clusters: map[string]Cluster{
			"foo": foo,
		},
		clusterName: "foo",
		want:        foo,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := registry{
				clientFactory: tt.clientFactory,
				clusters:      tt.clusters,
			}
			got := c.Lookup(tt.clusterName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_registry_Register(t *testing.T) {
	foo := NewClusterFromKubeconfig("foo", "bar")
	bar := NewClusterFromKubeconfig("bar", "baz")
	tests := []struct {
		name          string
		clientFactory clientFactory
		clusters      map[string]Cluster
		clusterName   string
		cluster       Cluster
		want          Registry
	}{{
		name:        "default",
		clusterName: DefaultClient,
		cluster:     foo,
		want: registry{
			clusters: map[string]Cluster{
				DefaultClient: foo,
			},
		},
	}, {
		name: "override",
		clusters: map[string]Cluster{
			"foo": foo,
		},
		clusterName: "foo",
		cluster:     bar,
		want: registry{
			clusters: map[string]Cluster{
				"foo": bar,
			},
		},
	}, {
		name: "foo and bar",
		clusters: map[string]Cluster{
			"foo": foo,
		},
		clusterName: "bar",
		cluster:     bar,
		want: registry{
			clusters: map[string]Cluster{
				"foo": foo,
				"bar": bar,
			},
		},
		// }, {
		// 	name:          "factory",
		// 	clientFactory: defaultClientFactory,
		// 	clusterName:   "bar",
		// 	cluster:       bar,
		// 	want: registry{
		// 		clientFactory: defaultClientFactory,
		// 		clusters: map[string]Cluster{
		// 			"bar": bar,
		// 		},
		// 	},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := registry{
				clientFactory: tt.clientFactory,
				clusters:      tt.clusters,
			}
			got := c.Register(tt.clusterName, tt.cluster)
			assert.Equal(t, tt.want, got)
			assert.NotEqual(t, c, got)
		})
	}
}

func Test_registry_Build(t *testing.T) {
	factory := func(Cluster) (*rest.Config, client.Client, error) {
		return &rest.Config{}, &tclient.FakeClient{}, nil
	}
	tests := []struct {
		name          string
		clientFactory clientFactory
		clusters      map[string]Cluster
		cluster       Cluster
		config        *rest.Config
		client        client.Client
		wantErr       bool
	}{{
		name:    "nil",
		wantErr: false,
	}, {
		name:          "not nil",
		clientFactory: factory,
		config:        &rest.Config{},
		client:        &tclient.FakeClient{},
		wantErr:       false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := registry{
				clientFactory: tt.clientFactory,
				clusters:      tt.clusters,
			}
			config, client, err := c.Build(tt.cluster)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, config)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.config, config)
				assert.Equal(t, tt.client, client)
			}
		})
	}
}

func Test_defaultClientFactory(t *testing.T) {
	tests := []struct {
		name    string
		cluster Cluster
		config  *rest.Config
		client  client.Client
		wantErr bool
	}{{
		name:    "nil",
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, client, err := defaultClientFactory(tt.cluster)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, config)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.config, config)
				assert.Equal(t, tt.client, client)
			}
		})
	}
}
