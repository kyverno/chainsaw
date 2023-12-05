package processors

import (
	"context"
	"slices"
	"time"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type cleaner struct {
	namespacer namespacer.Namespacer
	operations []operation
}

func newCleaner(namespacer namespacer.Namespacer) *cleaner {
	return &cleaner{
		namespacer: namespacer,
	}
}

func (c *cleaner) register(obj unstructured.Unstructured, client client.Client, timeout *time.Duration) {
	c.operations = append(c.operations, operation{
		continueOnError: true,
		timeout:         timeout,
		operation:       opdelete.New(client, obj, c.namespacer),
	})
}

func (c *cleaner) run(ctx context.Context) {
	slices.Reverse(c.operations)
	for _, operation := range c.operations {
		operation.execute(ctx)
	}
}
