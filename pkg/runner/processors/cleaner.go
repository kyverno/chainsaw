package processors

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type cleaner struct {
	namespacer namespacer.Namespacer
	delay      *metav1.Duration
	operations []operation
}

func newCleaner(namespacer namespacer.Namespacer, delay *metav1.Duration) *cleaner {
	return &cleaner{
		namespacer: namespacer,
		delay:      delay,
	}
}

func (c *cleaner) register(obj unstructured.Unstructured, client client.Client, timeout *time.Duration) {
	c.operations = append(c.operations, newOperation(
		true,
		timeout,
		opdelete.New(client, obj, c.namespacer, false),
		nil,
		nil,
		client,
	))
}

func (c *cleaner) run(ctx context.Context) {
	if c.delay != nil {
		time.Sleep(c.delay.Duration)
	}
	for i := len(c.operations) - 1; i >= 0; i-- {
		c.operations[i].execute(ctx, nil)
	}
}
