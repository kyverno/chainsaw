package processors

import (
	"context"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/clock"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type OperationProcessor interface {
	Run(ctx context.Context, namespace string, test discovery.Test, step v1alpha1.TestStepSpec, operation v1alpha1.Operation)
}

func NewOperationProcessor(config v1alpha1.ConfigurationSpec, client operations.Client, clock clock.PassiveClock) OperationProcessor {
	return &operationProcessor{
		config: config,
		client: client,
		clock:  clock,
	}
}

type operationProcessor struct {
	config v1alpha1.ConfigurationSpec
	client operations.Client
	clock  clock.PassiveClock
}

func (p *operationProcessor) Run(ctx context.Context, namespace string, test discovery.Test, step v1alpha1.TestStepSpec, operation v1alpha1.Operation) {
	fail := func(t *testing.T, continueOnError *bool) {
		t.Helper()
		if continueOnError != nil && *continueOnError {
			t.Fail()
		} else {
			t.FailNow()
		}
	}
	t := testing.FromContext(ctx)
	// Handle Delete
	if operation.Delete != nil {
		var resource unstructured.Unstructured
		resource.SetAPIVersion(operation.Delete.APIVersion)
		resource.SetKind(operation.Delete.Kind)
		resource.SetName(operation.Delete.Name)
		resource.SetNamespace(operation.Delete.Namespace)
		resource.SetLabels(operation.Delete.Labels)
		if err := p.client.Delete(ctx, operation.Timeout, &resource); err != nil {
			fail(t, operation.ContinueOnError)
		}
	}
	// Handle Exec
	if operation.Command != nil {
		if err := p.client.Command(ctx, operation.Timeout, *operation.Command); err != nil {
			fail(t, operation.ContinueOnError)
		}
	}
	if operation.Script != nil {
		if err := p.client.Script(ctx, operation.Timeout, *operation.Script); err != nil {
			fail(t, operation.ContinueOnError)
		}
	}
	var cleaner cleanup.Cleaner
	if !cleanup.Skip(p.config.SkipDelete, test.Spec.SkipDelete, step.SkipDelete) {
		cleaner = func(obj ctrlclient.Object, c client.Client) {
			t.Cleanup(func() {
				if err := p.client.Delete(ctx, nil, obj); err != nil {
					t.Fail()
				}
			})
		}
	}
	// Handle Apply
	if operation.Apply != nil {
		var resources []ctrlclient.Object
		if operation.Apply.Resource != nil {
			resources = append(resources, operation.Apply.Resource)
		} else {
			loaded, err := resource.Load(filepath.Join(test.BasePath, operation.Apply.File))
			if err != nil {
				logging.FromContext(ctx).Log("LOAD  ", color.BoldRed, err)
				fail(t, operation.ContinueOnError)
			}
			for i := range loaded {
				resources = append(resources, &loaded[i])
			}
		}
		shouldFail := operation.Apply.ShouldFail != nil && *operation.Apply.ShouldFail
		for _, resource := range resources {
			if err := p.client.Apply(ctx, operation.Timeout, resource, shouldFail, operation.Apply.DryRun, cleaner); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}
	// Handle Create
	if operation.Create != nil {
		var resources []ctrlclient.Object
		if operation.Create.Resource != nil {
			resources = append(resources, operation.Create.Resource)
		} else {
			loaded, err := resource.Load(filepath.Join(test.BasePath, operation.Create.File))
			if err != nil {
				logging.FromContext(ctx).Log("LOAD  ", color.BoldRed, err)
				fail(t, operation.ContinueOnError)
			}
			for i := range loaded {
				resources = append(resources, &loaded[i])
			}
		}
		shouldFail := operation.Create.ShouldFail != nil && *operation.Create.ShouldFail
		for _, resource := range resources {
			if err := p.client.Create(ctx, operation.Timeout, resource, shouldFail, operation.Create.DryRun, cleaner); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}
	// Handle Assert
	if operation.Assert != nil {
		resources, err := resource.Load(filepath.Join(test.BasePath, operation.Assert.File))
		if err != nil {
			logging.FromContext(ctx).Log("LOAD  ", color.BoldRed, err)
			fail(t, operation.ContinueOnError)
		}
		for _, resource := range resources {
			if err := p.client.Assert(ctx, operation.Timeout, resource); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}
	// Handle Error
	if operation.Error != nil {
		resources, err := resource.Load(filepath.Join(test.BasePath, operation.Error.File))
		if err != nil {
			logging.FromContext(ctx).Log("LOAD  ", color.BoldRed, err)
			fail(t, operation.ContinueOnError)
		}
		for _, resource := range resources {
			if err := p.client.Error(ctx, operation.Timeout, resource); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}
}
