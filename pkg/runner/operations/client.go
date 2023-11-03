package operations

import (
	"context"

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
	Apply(ctx context.Context, timeout *metav1.Duration, obj ctrlclient.Object, shouldFail bool, cleanup CleanupFunc) (_err error)
	Assert(ctx context.Context, timeout *metav1.Duration, expected unstructured.Unstructured) (_err error)
	Create(ctx context.Context, timeout *metav1.Duration, obj ctrlclient.Object, shouldFail bool, cleanup CleanupFunc) (_err error)
	Delete(ctx context.Context, timeout *metav1.Duration, obj ctrlclient.Object) error
	Error(ctx context.Context, timeout *metav1.Duration, expected unstructured.Unstructured) (_err error)
	Exec(ctx context.Context, exec v1alpha1.Exec, log bool, namespace string) error
}

type opClient struct {
	namespacer namespacer.Namespacer
	client     client.Client
	config     v1alpha1.ConfigurationSpec
	test       v1alpha1.TestSpec
	step       v1alpha1.TestStepSpec
}

func NewClient(
	namespacer namespacer.Namespacer,
	client client.Client,
	config v1alpha1.ConfigurationSpec,
	test v1alpha1.TestSpec,
	step v1alpha1.TestStepSpec,
) Client {
	return &opClient{
		namespacer: namespacer,
		client:     client,
		config:     config,
		test:       test,
		step:       step,
	}
}

func (c *opClient) Apply(ctx context.Context, to *metav1.Duration, obj ctrlclient.Object, shouldFail bool, cleanup func(ctrlclient.Object, client.Client)) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(obj); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultApplyTimeout, c.config.Timeouts.Apply, c.test.Timeouts.Apply, c.step.Timeouts.Apply, to)
	defer cancel()
	return operationApply(ctx, logger, obj, c.client, shouldFail, cleanup)
}

func (c *opClient) Assert(ctx context.Context, to *metav1.Duration, expected unstructured.Unstructured) (_err error) {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(&expected); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultAssertTimeout, c.config.Timeouts.Assert, c.test.Timeouts.Assert, c.step.Timeouts.Assert, to)
	defer cancel()
	return operationAssert(ctx, logger, expected, c.client)
}

func (c *opClient) Create(ctx context.Context, to *metav1.Duration, obj ctrlclient.Object, shouldFail bool, cleanup func(ctrlclient.Object, client.Client)) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(obj); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultApplyTimeout, c.config.Timeouts.Apply, c.test.Timeouts.Apply, c.step.Timeouts.Apply, to)
	defer cancel()
	return operationCreate(ctx, logger, obj, c.client, shouldFail, cleanup)
}

func (c *opClient) Delete(ctx context.Context, to *metav1.Duration, obj ctrlclient.Object) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(obj); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultDeleteTimeout, c.config.Timeouts.Delete, c.test.Timeouts.Delete, c.step.Timeouts.Delete, to)
	defer cancel()
	return operationDelete(ctx, logger, obj, c.client)
}

func (c *opClient) Error(ctx context.Context, to *metav1.Duration, expected unstructured.Unstructured) (_err error) {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(&expected); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		// fail(t, operation.ContinueOnError)
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultErrorTimeout, c.config.Timeouts.Error, c.test.Timeouts.Error, c.step.Timeouts.Error, to)
	defer cancel()
	return operationError(ctx, logger, expected, c.client)
}

func (c *opClient) Exec(ctx context.Context, exec v1alpha1.Exec, log bool, namespace string) error {
	ctx, cancel := timeout.Context(ctx, timeout.DefaultExecTimeout, c.config.Timeouts.Exec, c.test.Timeouts.Exec, c.step.Timeouts.Exec, exec.Timeout)
	defer cancel()
	return operationExec(ctx, logging.FromContext(ctx), exec, log, namespace)
}
