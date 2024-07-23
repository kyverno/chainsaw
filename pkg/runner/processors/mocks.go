package processors

import (
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"k8s.io/client-go/rest"
)

type testContext struct {
	config   model.Configuration
	bindings binding.Bindings
	clusters clusters.Registry
}

func (tc *testContext) Bindings() binding.Bindings {
	return tc.bindings
}

func (tc *testContext) Clusters() clusters.Registry {
	return tc.clusters
}

func (tc *testContext) Cluster() (*rest.Config, client.Client, error) {
	return tc.clusters.Resolve(false)
}

func (tc *testContext) Configuration() model.Configuration {
	return tc.config
}

func (tc *testContext) WithBindings(bindings binding.Bindings) model.TestContext {
	return &testContext{
		config:   tc.config,
		bindings: bindings,
		clusters: tc.clusters,
	}
}

type registryMock struct {
	client client.Client
}

func (r registryMock) Register(string, clusters.Cluster) clusters.Registry {
	return r
}

func (r registryMock) Resolve(bool, ...string) (*rest.Config, client.Client, error) {
	return nil, r.client, nil
}
