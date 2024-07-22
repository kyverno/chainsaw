package simple

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func New(cfg *rest.Config) (client.Client, error) {
	var opts ctrlclient.Options
	return ctrlclient.New(cfg, opts)
}
