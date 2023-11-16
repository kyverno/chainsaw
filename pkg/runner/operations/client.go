package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations/apply"
	"github.com/kyverno/chainsaw/pkg/runner/operations/assert"
	"github.com/kyverno/chainsaw/pkg/runner/operations/command"
	"github.com/kyverno/chainsaw/pkg/runner/operations/create"
	"github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	cmderror "github.com/kyverno/chainsaw/pkg/runner/operations/error"
	"github.com/kyverno/chainsaw/pkg/runner/operations/script"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/kyverno/ext/output/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type OperationClient interface {
	Apply(context.Context, *metav1.Duration, ctrlclient.Object, bool, interface{}, cleanup.Cleaner) error
	Assert(context.Context, *metav1.Duration, unstructured.Unstructured) error
	Create(context.Context, *metav1.Duration, ctrlclient.Object, bool, interface{}, cleanup.Cleaner) error
	Delete(context.Context, *metav1.Duration, ctrlclient.Object) error
	Error(context.Context, *metav1.Duration, unstructured.Unstructured) error
	Command(context.Context, *metav1.Duration, v1alpha1.Command) error
	Script(context.Context, *metav1.Duration, v1alpha1.Script) error
}

type opClient struct {
	namespacer   namespacer.Namespacer
	client       client.Client
	config       v1alpha1.ConfigurationSpec
	test         v1alpha1.TestSpec
	stepTimeouts v1alpha1.Timeouts
}

func NewOperationClient(
	namespacer namespacer.Namespacer,
	client client.Client,
	config v1alpha1.ConfigurationSpec,
	test v1alpha1.TestSpec,
	stepTimeouts v1alpha1.Timeouts,
) OperationClient {
	return &opClient{
		namespacer:   namespacer,
		client:       client,
		config:       config,
		test:         test,
		stepTimeouts: stepTimeouts,
	}
}

func (c *opClient) Apply(ctx context.Context, to *metav1.Duration, obj ctrlclient.Object, dryRun bool, check interface{}, cleanup func(ctrlclient.Object, client.Client)) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(obj); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultApplyTimeout, c.config.Timeouts.Apply, c.test.Timeouts.Apply, c.stepTimeouts.Apply, to)
	defer cancel()
	operation := apply.New(c.client, obj, dryRun, cleanup, check)
	return operation.Exec(ctx)
}

func (c *opClient) Assert(ctx context.Context, to *metav1.Duration, expected unstructured.Unstructured) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(&expected); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultAssertTimeout, c.config.Timeouts.Assert, c.test.Timeouts.Assert, c.stepTimeouts.Assert, to)
	defer cancel()
	operation := assert.New(c.client, expected)
	return operation.Exec(ctx)
}

func (c *opClient) Create(ctx context.Context, to *metav1.Duration, obj ctrlclient.Object, dryRun bool, check interface{}, cleanup func(ctrlclient.Object, client.Client)) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(obj); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultApplyTimeout, c.config.Timeouts.Apply, c.test.Timeouts.Apply, c.stepTimeouts.Apply, to)
	defer cancel()
	operation := create.New(c.client, obj, dryRun, cleanup, check)
	return operation.Exec(ctx)
}

func (c *opClient) Delete(ctx context.Context, to *metav1.Duration, obj ctrlclient.Object) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(obj); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultDeleteTimeout, c.config.Timeouts.Delete, c.test.Timeouts.Delete, c.stepTimeouts.Delete, to)
	defer cancel()
	operation := delete.New(c.client, obj)
	return operation.Exec(ctx)
}

func (c *opClient) Error(ctx context.Context, to *metav1.Duration, expected unstructured.Unstructured) error {
	logger := logging.FromContext(ctx)
	if err := c.namespacer.Apply(&expected); err != nil {
		logger.Log("LOAD  ", color.BoldRed, err)
		return err
	}
	ctx, cancel := timeout.Context(ctx, timeout.DefaultErrorTimeout, c.config.Timeouts.Error, c.test.Timeouts.Error, c.stepTimeouts.Error, to)
	defer cancel()
	operation := cmderror.New(c.client, expected)
	return operation.Exec(ctx)
}

func (c *opClient) Command(ctx context.Context, to *metav1.Duration, exec v1alpha1.Command) error {
	ctx, cancel := timeout.Context(ctx, timeout.DefaultExecTimeout, c.config.Timeouts.Exec, c.test.Timeouts.Exec, c.stepTimeouts.Exec, to)
	defer cancel()
	operation := command.New(exec, c.namespacer.GetNamespace())
	return operation.Exec(ctx)
}

func (c *opClient) Script(ctx context.Context, to *metav1.Duration, exec v1alpha1.Script) error {
	ctx, cancel := timeout.Context(ctx, timeout.DefaultExecTimeout, c.config.Timeouts.Exec, c.test.Timeouts.Exec, c.stepTimeouts.Exec, to)
	defer cancel()
	operation := script.New(exec, c.namespacer.GetNamespace())
	return operation.Exec(ctx)
}
