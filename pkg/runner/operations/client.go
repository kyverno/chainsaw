package operations

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	Apply(obj ctrlclient.Object, shouldFail bool, cleanup CleanupFunc) (_err error)
	Assert(expected unstructured.Unstructured) (_err error)
	Delete(expected ctrlclient.Object) error
	Error(expected unstructured.Unstructured) (_err error)
	Exec(exec v1alpha1.Exec, log bool, namespace string) error
}

type opClient struct {
	logger     logging.Logger
	namespacer namespacer.Namespacer
	client     client.Client
}

func NewClient(logger logging.Logger, namespacer namespacer.Namespacer, client client.Client) Client {
	return &opClient{
		logger:     logger,
		namespacer: namespacer,
		client:     client,
	}
}

// Apply implements Client.
func (c *opClient) Apply(obj ctrlclient.Object, shouldFail bool, cleanup func(ctrlclient.Object, client.Client)) error {
	if err := c.namespacer.Apply(obj); err != nil {
		c.logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	operationCtx, cancel := timeoutCtx(DefaultApplyTimeout, config.Timeouts.Apply, test.Timeouts.Apply, step.Spec.Timeouts.Apply, operation.Timeout)
	defer cancel()
	return Apply(operationCtx, c.logger, obj, c.client, shouldFail, cleanup)
}

// Assert implements Client.
func (*opClient) Assert(expected unstructured.Unstructured) (_err error) {
	panic("unimplemented")
}

// Delete implements Client.
func (*opClient) Delete(expected ctrlclient.Object) error {
	panic("unimplemented")
}

// Error implements Client.
func (*opClient) Error(expected unstructured.Unstructured) (_err error) {
	panic("unimplemented")
}

// Exec implements Client.
func (*opClient) Exec(exec v1alpha1.Exec, log bool, namespace string) error {
	panic("unimplemented")
}
