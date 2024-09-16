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
		StartTime time.Time           `json:"startTime"`
		EndTime   time.Time           `json:"endTime"`
		Failure   *Failure            `json:"failure,omitempty"`
	}
	type StepReport struct {
		Name       string            `json:"name,omitempty"`
		StartTime  time.Time         `json:"startTime"`
		EndTime    time.Time         `json:"endTime"`
		Operations []OperationReport `json:"operations,omitempty"`
	}
	type TestReport struct {
		BasePath   string       `json:"basePath,omitempty"`
		Name       string       `json:"name,omitempty"`
		Concurrent *bool        `json:"concurrent,omitempty"`
		StartTime  time.Time    `json:"startTime"`
		EndTime    time.Time    `json:"endTime"`
		Namespace  string       `json:"namespace,omitempty"`
		Skipped    bool         `json:"skipped,omitempty"`
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
		testReport := TestReport{
			BasePath:   test.BasePath,
			Name:       test.Name,
			Concurrent: test.Concurrent,
			StartTime:  test.StartTime,
			EndTime:    test.EndTime,
			Namespace:  test.Namespace,
			Skipped:    test.Skipped,
		}
		for _, step := range test.Steps {
			stepReport := StepReport{
				Name:      step.Name,
				StartTime: step.StartTime,
				EndTime:   step.EndTime,
			}
			for _, operation := range step.Operations {
				operationReport := OperationReport{
					Name:      operation.Name,
					Type:      operation.Type,
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
