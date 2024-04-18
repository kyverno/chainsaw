package loader

import (
	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/pkg/ext/resource/loader"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

var (
	OpenApiClient      = openapiclient.NewLocalCRDFiles(data.Crds(), data.CrdsFolder)
	DefaultLoader, Err = loader.New(openapiclient.NewLocalCRDFiles(data.Crds(), data.CrdsFolder))
)
