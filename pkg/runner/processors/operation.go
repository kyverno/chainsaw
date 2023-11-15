package processors

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
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
	Run(ctx context.Context, namespace string, test discovery.Test, step v1alpha1.TestStepSpec, operation v1alpha1.Operation) []report.OperationReport
}

func NewOperationProcessor(config v1alpha1.ConfigurationSpec, operationClient operations.OperationClient, clock clock.PassiveClock) OperationProcessor {
	return &operationProcessor{
		config:          config,
		operationClient: operationClient,
		clock:           clock,
	}
}

type operationProcessor struct {
	config          v1alpha1.ConfigurationSpec
	operationClient operations.OperationClient
	clock           clock.PassiveClock
}

func (p *operationProcessor) Run(ctx context.Context, namespace string, test discovery.Test, step v1alpha1.TestStepSpec, operation v1alpha1.Operation) []report.OperationReport {
	fail := func(t *testing.T, continueOnError *bool) {
		t.Helper()
		if continueOnError != nil && *continueOnError {
			t.Fail()
		} else {
			t.FailNow()
		}
	}
	t := testing.FromContext(ctx)

	// Initialize an OperationReport
	var operationReports []report.OperationReport

	// Handle Delete
	if operation.Delete != nil {
		deleteOpReport := report.OperationReport{
			Name:      "Delete",
			StartTime: p.clock.Now(),
		}
		var resource unstructured.Unstructured
		resource.SetAPIVersion(operation.Delete.APIVersion)
		resource.SetKind(operation.Delete.Kind)
		resource.SetName(operation.Delete.Name)
		resource.SetNamespace(operation.Delete.Namespace)
		resource.SetLabels(operation.Delete.Labels)
		if err := p.operationClient.Delete(ctx, operation.Timeout, &resource); err != nil {
			deleteOpReport.Result = "Failed"
			deleteOpReport.Message = err.Error()
			fail(t, operation.ContinueOnError)
		} else {
			deleteOpReport.Result = "Success"
		}
		deleteOpReport.EndTime = p.clock.Now()
		operationReports = append(operationReports, deleteOpReport)
	}
	// Handle Exec
	if operation.Command != nil {
		commandOpReport := report.OperationReport{
			Name:      "Command",
			StartTime: p.clock.Now(),
		}
		if err := p.operationClient.Command(ctx, operation.Timeout, *operation.Command); err != nil {
			commandOpReport.Result = "Failed"
			commandOpReport.Message = err.Error()
			fail(t, operation.ContinueOnError)
		} else {
			commandOpReport.Result = "Success"
		}
	}
	if operation.Script != nil {
		scriptOpReport := report.OperationReport{
			Name:      "Script",
			StartTime: p.clock.Now(),
		}
		if err := p.operationClient.Script(ctx, operation.Timeout, *operation.Script); err != nil {
			scriptOpReport.Result = "Failed"
			scriptOpReport.Message = err.Error()
			fail(t, operation.ContinueOnError)
		} else {
			scriptOpReport.Result = "Success"
		}
	}
	var cleaner cleanup.Cleaner
	if !cleanup.Skip(p.config.SkipDelete, test.Spec.SkipDelete, step.SkipDelete) {
		cleaner = func(obj ctrlclient.Object, c client.Client) {
			t.Cleanup(func() {
				if err := p.operationClient.Delete(ctx, nil, obj); err != nil {
					t.Fail()
				}
			})
		}
	}
	// Handle Apply
	if operation.Apply != nil {
		applyOpReport := report.OperationReport{
			Name:      "Apply",
			StartTime: p.clock.Now(),
		}
		var resources []ctrlclient.Object
		var applyErrors []string
		if operation.Apply.Resource != nil {
			resources = append(resources, operation.Apply.Resource)
		} else {
			loaded, err := resource.Load(filepath.Join(test.BasePath, operation.Apply.File))
			if err != nil {
				logging.FromContext(ctx).Log("LOAD  ", color.BoldRed, err)
				fail(t, operation.ContinueOnError)
				applyOpReport.Message = fmt.Sprintf("Resource loading error: %s", err.Error())
				applyErrors = append(applyErrors, applyOpReport.Message)
			}
			for i := range loaded {
				resources = append(resources, &loaded[i])
			}
		}

		shouldFail := operation.Apply.ShouldFail != nil && *operation.Apply.ShouldFail
		dryRun := operation.Apply.DryRun != nil && *operation.Apply.DryRun

		for _, res := range resources {
			if err := p.operationClient.Apply(ctx, operation.Timeout, res, shouldFail, dryRun, cleaner); err != nil {
				errMsg := fmt.Sprintf("Apply error for resource %v: %s", res, err.Error())
				logging.FromContext(ctx).Log("APPLY ", color.BoldRed, errMsg)
				fail(t, operation.ContinueOnError)
				applyErrors = append(applyErrors, errMsg)
			}
		}

		applyOpReport.EndTime = p.clock.Now()
		if len(applyErrors) > 0 {
			applyOpReport.Result = "Failed"
			applyOpReport.Message = strings.Join(applyErrors, "; ")
		} else {
			applyOpReport.Result = "Success"
		}

		operationReports = append(operationReports, applyOpReport)
	}
	// Handle Create
	if operation.Create != nil {
		createOpReport := report.OperationReport{
			Name:      "Create",
			StartTime: p.clock.Now(),
		}
		var resources []ctrlclient.Object
		var createErrors []string
		if operation.Create.Resource != nil {
			resources = append(resources, operation.Create.Resource)
		} else {
			loaded, err := resource.Load(filepath.Join(test.BasePath, operation.Create.File))
			if err != nil {
				logging.FromContext(ctx).Log("LOAD  ", color.BoldRed, err)
				createErrors = append(createErrors, err.Error())
				fail(t, operation.ContinueOnError)
			}
			for i := range loaded {
				resources = append(resources, &loaded[i])
			}
		}
		shouldFail := operation.Create.ShouldFail != nil && *operation.Create.ShouldFail
		dryRun := operation.Create.DryRun != nil && *operation.Create.DryRun

		for _, res := range resources {
			if err := p.operationClient.Create(ctx, operation.Timeout, res, shouldFail, dryRun, cleaner); err != nil {
				errMsg := fmt.Sprintf("Create error for resource %v: %s", res, err.Error())
				createErrors = append(createErrors, errMsg)
				fail(t, operation.ContinueOnError)
			}
		}

		createOpReport.EndTime = p.clock.Now()
		if len(createErrors) > 0 {
			createOpReport.Result = "Failed"
			createOpReport.Message = strings.Join(createErrors, "; ")
		} else {
			createOpReport.Result = "Success"
		}

		operationReports = append(operationReports, createOpReport)
	}
	// Handle Assert
	if operation.Assert != nil {
		assertOpReport := report.OperationReport{
			Name:      "Assert",
			StartTime: p.clock.Now(),
		}
		var assertErrors []string
		resources, err := resource.Load(filepath.Join(test.BasePath, operation.Assert.File))
		if err != nil {
			errMsg := fmt.Sprintf("Resource loading error: %s", err.Error())
			logging.FromContext(ctx).Log("LOAD  ", color.BoldRed, errMsg)
			fail(t, operation.ContinueOnError)
			assertErrors = append(assertErrors, errMsg)
		}

		for _, res := range resources {
			if err := p.operationClient.Assert(ctx, operation.Timeout, res); err != nil {
				errMsg := fmt.Sprintf("Assert error for resource %v: %s", res, err.Error())
				assertErrors = append(assertErrors, errMsg)
				fail(t, operation.ContinueOnError)
			}
		}

		assertOpReport.EndTime = p.clock.Now()
		if len(assertErrors) > 0 {
			assertOpReport.Result = "Failed"
			assertOpReport.Message = strings.Join(assertErrors, "; ")
		} else {
			assertOpReport.Result = "Success"
		}

		operationReports = append(operationReports, assertOpReport)
	}
	// Handle Error
	if operation.Error != nil {
		errorOpReport := report.OperationReport{
			Name:      "Error",
			StartTime: p.clock.Now(),
		}
		var errorErrors []string
		resources, err := resource.Load(filepath.Join(test.BasePath, operation.Error.File))
		if err != nil {
			errMsg := fmt.Sprintf("Resource loading error: %s", err.Error())
			logging.FromContext(ctx).Log("LOAD  ", color.BoldRed, errMsg)
			fail(t, operation.ContinueOnError)
			errorErrors = append(errorErrors, errMsg)
		}

		for _, res := range resources {
			if err := p.operationClient.Error(ctx, operation.Timeout, res); err != nil {
				errMsg := fmt.Sprintf("Error operation failed for resource %v: %s", res, err.Error())
				errorErrors = append(errorErrors, errMsg)
				fail(t, operation.ContinueOnError)
			}
		}
		errorOpReport.EndTime = p.clock.Now()
		if len(errorErrors) > 0 {
			errorOpReport.Result = "Failed"
			errorOpReport.Message = strings.Join(errorErrors, "; ")
		} else {
			errorOpReport.Result = "Success"
		}
		operationReports = append(operationReports, errorOpReport)
	}

	return operationReports
}
