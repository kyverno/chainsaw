package loaders

import (
	"sync"

	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/pkg/ext/resource/loader"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

func defaultLoader() (loader.Loader, error) {
	crdsFs, err := data.Crds()
	if err != nil {
		return nil, err
	}
	return loader.New(openapiclient.NewLocalCRDFiles(crdsFs))
}

var DefaultLoader = sync.OnceValues(defaultLoader)
