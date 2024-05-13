package clusters

import (
	"github.com/kyverno/chainsaw/pkg/client"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"k8s.io/client-go/rest"
)

func makeClient(config *rest.Config, dryRun bool) (client.Client, error) {
	c, err := client.New(config)
	if err != nil {
		return nil, err
	}
	c = runnerclient.New(c)
	if !dryRun {
		return c, nil
	}
	return client.DryRun(c), nil
}
