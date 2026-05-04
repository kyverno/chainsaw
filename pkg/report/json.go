package report

import (
	"encoding/json"
	"os"
	"time"

	"github.com/kyverno/chainsaw/pkg/model"
)

func saveJson(report *model.Report, file string) error {
	type Failure struct {
		Err string `json:"error,omitempty"`
	}
	type OperationReport struct {
		Name      string              `json:"name,omitempty"`
		Type      model.OperationType `json:"type,omitempty"`
		Status    string              `json:"status"`
		StartTime time.Time           `json:"startTime"`
		EndTime   time.Time           `json:"endTime"`
		Failure   *Failure            `json:"failure,omitempty"`
	}
	type StepReport struct {
		Name       string            `json:"name,omitempty"`
		Status     string            `json:"status"`
		StartTime  time.Time         `json:"startTime"`
		EndTime    time.Time         `json:"endTime"`
		Operations []OperationReport `json:"operations,omitempty"`
	}
	type TestReport struct {
		BasePath   string       `json:"basePath,omitempty"`
		Name       string       `json:"name,omitempty"`
		Concurrent *bool        `json:"concurrent,omitempty"`
		Status     string       `json:"status"`
		StartTime  time.Time    `json:"startTime"`
		EndTime    time.Time    `json:"endTime"`
		Namespace  string       `json:"namespace,omitempty"`
		Steps      []StepReport `json:"steps,omitempty"`
	}
	type Report struct {
		Name      string       `json:"name,omitempty"`
		StartTime time.Time    `json:"startTime"`
		EndTime   time.Time    `json:"endTime"`
		Tests     []TestReport `json:"tests,omitempty"`
	}
	testsReport := Report{
		Name:      report.Name,
		StartTime: report.StartTime,
		EndTime:   report.EndTime,
	}
	for _, test := range report.Tests {
		testStatus := "passed"
		if test.Skipped {
			testStatus = "skipped"
		} else if test.Failed() {
			testStatus = "failed"
		}
		testReport := TestReport{
			BasePath:   test.BasePath,
			Name:       test.Name,
			Concurrent: test.Concurrent,
			Status:     testStatus,
			StartTime:  test.StartTime,
			EndTime:    test.EndTime,
			Namespace:  test.Namespace,
		}
		for _, step := range test.Steps {
			stepStatus := "passed"
			if step.Failed() {
				stepStatus = "failed"
			}
			stepReport := StepReport{
				Name:      step.Name,
				Status:    stepStatus,
				StartTime: step.StartTime,
				EndTime:   step.EndTime,
			}
			for _, operation := range step.Operations {
				opStatus := "passed"
				if operation.Err != nil {
					opStatus = "failed"
				}
				operationReport := OperationReport{
					Name:      operation.Name,
					Type:      operation.Type,
					Status:    opStatus,
					StartTime: operation.StartTime,
					EndTime:   operation.EndTime,
				}
				if operation.Err != nil {
					operationReport.Failure = &Failure{
						Err: operation.Err.Error(),
					}
				}
				stepReport.Operations = append(stepReport.Operations, operationReport)
			}
			testReport.Steps = append(testReport.Steps, stepReport)
		}
		testsReport.Tests = append(testsReport.Tests, testReport)
	}
	data, err := json.MarshalIndent(testsReport, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o600)
}
