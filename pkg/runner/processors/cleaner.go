package processors

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
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

func (c *cleaner) register(obj unstructured.Unstructured, cluster clusters.Cluster, timeout *time.Duration) {
	client := cluster.Client()
	c.operations = append(c.operations, newLazyOperation(
		nil,
		OperationInfo{},
		true,
		timeout,
		func(_ context.Context, _ binding.Bindings) (operations.Operation, error) {
			return opdelete.New(client, obj, c.namespacer, false), nil
		},
		nil,
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
