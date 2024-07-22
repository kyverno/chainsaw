package clusters

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/client/dryrun"
	"github.com/kyverno/chainsaw/pkg/client/simple"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"k8s.io/client-go/rest"
)

func makeClient(config *rest.Config, dryRun bool) (client.Client, error) {
	c, err := simple.New(config)
	if err != nil {
		return nil, err
	}
	c = runnerclient.New(c)
	if !dryRun {
		return c, nil
	}
	return dryrun.New(c), nil
}
