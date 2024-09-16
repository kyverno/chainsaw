package loaders

import (
	"io/fs"
	"sync"

	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/pkg/ext/resource/loader"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

func defaultLoader(_fs func() (fs.FS, error)) (loader.Loader, error) {
	if _fs == nil {
		_fs = data.Crds
	}
	crdsFs, err := _fs()
	if err != nil {
		return nil, err
	}
	return loader.New(openapiclient.NewLocalCRDFiles(crdsFs))
}

var DefaultLoader = sync.OnceValues(func() (loader.Loader, error) { return defaultLoader(nil) })
