package operations

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/kyverno/ext/output/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	Apply(timeout *metav1.Duration, obj ctrlclient.Object, shouldFail bool, cleanup CleanupFunc) (_err error)
	Assert(timeout *metav1.Duration, expected unstructured.Unstructured) (_err error)
	Delete(timeout *metav1.Duration, obj ctrlclient.Object) error
	Error(timeout *metav1.Duration, expected unstructured.Unstructured) (_err error)
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
func (c *opClient) Apply(to *metav1.Duration, obj ctrlclient.Object, shouldFail bool, cleanup func(ctrlclient.Object, client.Client)) error {
	if err := c.namespacer.Apply(obj); err != nil {
		c.logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(timeout.DefaultApplyTimeout, c.config.Timeouts.Apply, c.test.Timeouts.Apply, c.step.Timeouts.Apply, to)
	defer cancel()
	return operationApply(ctx, c.logger, obj, c.client, shouldFail, cleanup)
}

// Assert implements Client.
func (c *opClient) Assert(to *metav1.Duration, expected unstructured.Unstructured) (_err error) {
	if err := c.namespacer.Apply(&expected); err != nil {
		c.logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(timeout.DefaultAssertTimeout, c.config.Timeouts.Assert, c.test.Timeouts.Assert, c.step.Timeouts.Assert, to)
	defer cancel()
	return operationAssert(ctx, c.logger, expected, c.client)
}

// Delete implements Client.
func (c *opClient) Delete(to *metav1.Duration, obj ctrlclient.Object) error {
	if err := c.namespacer.Apply(obj); err != nil {
		c.logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(timeout.DefaultDeleteTimeout, c.config.Timeouts.Delete, c.test.Timeouts.Delete, c.step.Timeouts.Delete, to)
	defer cancel()
	return operationDelete(ctx, c.logger, obj, c.client)
}

// Error implements Client.
func (c *opClient) Error(to *metav1.Duration, expected unstructured.Unstructured) (_err error) {
	if err := c.namespacer.Apply(&expected); err != nil {
		c.logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(timeout.DefaultErrorTimeout, c.config.Timeouts.Error, c.test.Timeouts.Error, c.step.Timeouts.Error, to)
	defer cancel()
	return operationError(ctx, c.logger, expected, c.client)
}

// Exec implements Client.
func (c *opClient) Exec(exec v1alpha1.Exec, log bool, namespace string) error {
	ctx, cancel := timeout.Context(timeout.DefaultExecTimeout, c.config.Timeouts.Exec, c.test.Timeouts.Exec, c.step.Timeouts.Exec, exec.Timeout)
	defer cancel()
	return operationExec(ctx, c.logger, exec, log, namespace)
}
