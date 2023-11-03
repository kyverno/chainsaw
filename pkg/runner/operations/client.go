package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/kyverno/ext/output/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	Apply(context.Context, *metav1.Duration, ctrlclient.Object, bool, cleanup.Cleaner) error
	Assert(context.Context, *metav1.Duration, unstructured.Unstructured) error
	Create(context.Context, *metav1.Duration, ctrlclient.Object, bool, cleanup.Cleaner) error
	Delete(context.Context, *metav1.Duration, ctrlclient.Object) error
	Error(context.Context, *metav1.Duration, unstructured.Unstructured) error
	Exec(context.Context, *metav1.Duration, v1alpha1.Exec, bool, string) error
}

type opClient struct {
	namespacer   namespacer.Namespacer
	client       client.Client
	config       v1alpha1.ConfigurationSpec
	test         v1alpha1.TestSpec
	stepTimeouts v1alpha1.Timeouts
}

func NewClient(
	namespacer namespacer.Namespacer,
	client client.Client,
	config v1alpha1.ConfigurationSpec,
	test v1alpha1.TestSpec,
	stepTimeouts v1alpha1.Timeouts,
) Client {
	return &opClient{
		namespacer:   namespacer,
		client:       client,
		config:       config,
		test:         test,
		stepTimeouts: stepTimeouts,
	}
}

func (c *opClient) Apply(ctx context.Context, to *metav1.Duration, obj ctrlclient.Object, shouldFail bool, cleanup func(ctrlclient.Object, client.Client)) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(obj); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultApplyTimeout, c.config.Timeouts.Apply, c.test.Timeouts.Apply, c.stepTimeouts.Apply, to)
	defer cancel()
	return operationApply(ctx, obj, c.client, shouldFail, cleanup)
}

func (c *opClient) Assert(ctx context.Context, to *metav1.Duration, expected unstructured.Unstructured) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(&expected); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultAssertTimeout, c.config.Timeouts.Assert, c.test.Timeouts.Assert, c.stepTimeouts.Assert, to)
	defer cancel()
	return operationAssert(ctx, expected, c.client)
}

func (c *opClient) Create(ctx context.Context, to *metav1.Duration, obj ctrlclient.Object, shouldFail bool, cleanup func(ctrlclient.Object, client.Client)) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(obj); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultApplyTimeout, c.config.Timeouts.Apply, c.test.Timeouts.Apply, c.stepTimeouts.Apply, to)
	defer cancel()
	return operationCreate(ctx, obj, c.client, shouldFail, cleanup)
}

func (c *opClient) Delete(ctx context.Context, to *metav1.Duration, obj ctrlclient.Object) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(obj); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultDeleteTimeout, c.config.Timeouts.Delete, c.test.Timeouts.Delete, c.stepTimeouts.Delete, to)
	defer cancel()
	return operationDelete(ctx, obj, c.client)
}

func (c *opClient) Error(ctx context.Context, to *metav1.Duration, expected unstructured.Unstructured) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(&expected); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultErrorTimeout, c.config.Timeouts.Error, c.test.Timeouts.Error, c.stepTimeouts.Error, to)
	defer cancel()
	return operationError(ctx, expected, c.client)
}

func (c *opClient) Exec(ctx context.Context, to *metav1.Duration, exec v1alpha1.Exec, log bool, namespace string) error {
	ctx, cancel := timeout.Context(ctx, timeout.DefaultExecTimeout, c.config.Timeouts.Exec, c.test.Timeouts.Exec, c.stepTimeouts.Exec, to)
	defer cancel()
	return operationExec(ctx, exec, log, namespace)
}
