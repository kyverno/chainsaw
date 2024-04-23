package loader

import (
	"io/fs"

	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/pkg/ext/resource/loader"
	"k8s.io/client-go/openapi"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

var (
	OpenApiClient = func() openapi.Client {
		fs, err := fs.Sub(data.Crds(), data.CrdsFolder)
		if err != nil {
			panic(err)
		}
		return openapiclient.NewLocalCRDFiles(fs)
	}()
	DefaultLoader, Err = loader.New(OpenApiClient)
)
