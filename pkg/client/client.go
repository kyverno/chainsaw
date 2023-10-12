package client

import (
	"context"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // package needed for auth providers like GCP
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	Apply(context.Context, ctrlclient.Object) (updated bool, err error)
}

type c = ctrlclient.Client

type client struct {
	c
}

func New(cfg *rest.Config) (Client, error) {
	var opts ctrlclient.Options
	c, err := ctrlclient.New(cfg, opts)
	if err != nil {
		return nil, err
	}
	return &client{
		c: c,
	}, nil
}

func (c *client) Apply(ctx context.Context, obj ctrlclient.Object) (bool, error) {
	var actual unstructured.Unstructured
	err := c.Get(ctx, ObjectKey(obj), &actual)
	if err == nil {
		bytes, err := PatchObject(&actual, obj)
		if err != nil {
			return false, err
		}
		return true, c.Patch(ctx, &actual, ctrlclient.RawPatch(types.MergePatchType, bytes))
	} else if k8serrors.IsNotFound(err) {
		return false, c.Create(ctx, obj)
	}
	// TODO: context timeout
	return false, err
}
