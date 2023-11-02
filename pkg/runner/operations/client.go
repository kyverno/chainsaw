package operations

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
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
	config     v1alpha1.ConfigurationSpec
	test       v1alpha1.TestSpec
	step       v1alpha1.TestStepSpec
}

func NewClient(
	logger logging.Logger,
	namespacer namespacer.Namespacer,
	client client.Client,
	config v1alpha1.ConfigurationSpec,
	test v1alpha1.TestSpec,
	step v1alpha1.TestStepSpec,
) Client {
	return &opClient{
		logger:     logger,
		namespacer: namespacer,
		client:     client,
		config:     config,
		test:       test,
		step:       step,
	}
}

// Apply implements Client.
func (c *opClient) Apply(obj ctrlclient.Object, shouldFail bool, cleanup func(ctrlclient.Object, client.Client)) error {
	if err := c.namespacer.Apply(obj); err != nil {
		c.logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(timeout.DefaultApplyTimeout, c.config.Timeouts.Apply, c.test.Timeouts.Apply, c.step.Timeouts.Apply, nil /*c.operation.Timeout*/)
	defer cancel()
	return operationApply(ctx, c.logger, obj, c.client, shouldFail, cleanup)
}

// Assert implements Client.
func (c *opClient) Assert(expected unstructured.Unstructured) (_err error) {
	if err := c.namespacer.Apply(&expected); err != nil {
		c.logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(timeout.DefaultAssertTimeout, c.config.Timeouts.Assert, c.test.Timeouts.Assert, c.step.Timeouts.Assert, nil /*c.operation.Timeout*/)
	defer cancel()
	return operationAssert(ctx, c.logger, expected, c.client)
}

// Delete implements Client.
func (*opClient) Delete(expected ctrlclient.Object) error {
	panic("unimplemented")
}

// Error implements Client.
func (c *opClient) Error(expected unstructured.Unstructured) (_err error) {
	if err := c.namespacer.Apply(&expected); err != nil {
		c.logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(timeout.DefaultErrorTimeout, c.config.Timeouts.Error, c.test.Timeouts.Error, c.step.Timeouts.Error, nil /*c.operation.Timeout*/)
	defer cancel()
	return operationError(ctx, c.logger, expected, c.client)
}

// Exec implements Client.
func (*opClient) Exec(exec v1alpha1.Exec, log bool, namespace string) error {
	panic("unimplemented")
}
