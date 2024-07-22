package processors

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/model"
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

func (c *cleaner) register(ops ...operation) {
	c.operations = append(c.operations, ops...)
}

func (c *cleaner) addObject(obj unstructured.Unstructured, client client.Client, timeout *time.Duration) {
	c.register(newOperation(
		model.OperationInfo{},
		true,
		timeout,
		func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
			return opdelete.New(client, obj, c.namespacer, false, metav1.DeletePropagationBackground), bindings, nil
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

func (c *cleaner) isEmpty() bool {
	return len(c.operations) == 0
}
