package client

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	ctrlclient.Reader
	ctrlclient.Writer
	IsObjectNamespaced(obj runtime.Object) (bool, error)
}

func New(cfg *rest.Config) (Client, error) {
	var opts ctrlclient.Options
	return ctrlclient.New(cfg, opts)
}
