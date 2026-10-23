package report

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jstemmer/go-junit-report/v2/junit"
	"github.com/kyverno/chainsaw/pkg/model"
	"go.uber.org/multierr"
)

// systemOut gives the system-out element of a test case, or nil when there is no output.
func systemOut(output string) *junit.Output {
	if output == "" {
		return nil
	}
	return &junit.Output{Data: output}
}

func durationInSecondsString(start, end time.Time) string {
	duration := end.Sub(start)
	return fmt.Sprintf("%.6f", duration.Seconds())
}

func saveJUnitTest(report *model.Report, file string) error {
	testSuites := &junit.Testsuites{
		Name: report.Name,
		Time: durationInSecondsString(report.StartTime, report.EndTime),
	}
	addTestSuite := func(folder string, tests ...*model.TestReport) {
		testSuite := junit.Testsuite{
			Name: folder,
		}
		for _, test := range tests {
			testCase := junit.Testcase{
				Name: test.Name,
				Time: durationInSecondsString(test.StartTime, test.EndTime),
			}
			if test.Skipped {
				testCase.Skipped = &junit.Result{}
			} else {
				var errs []error
				for _, step := range test.Steps {
					for _, operation := range step.Operations {
						if operation.Err != nil {
							errs = append(errs, operation.Err)
						}
					}
				}
				if err := multierr.Combine(errs...); err != nil {
					testCase.Failure = &junit.Result{
						Message: err.Error(),
					}
				}
			}
			testSuite.AddTestcase(testCase)
		}
		testSuites.AddSuite(testSuite)
	}
	perFolder := map[string][]*model.TestReport{}
	for _, test := range report.Tests {
		perFolder[test.BasePath] = append(perFolder[test.BasePath], test)
	}
	for folder, tests := range perFolder {
		addTestSuite(folder, tests...)
	}
	data, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o600)
}

func saveJUnitStep(report *model.Report, file string) error {
	testSuites := &junit.Testsuites{
		Name: report.Name,
		Time: durationInSecondsString(report.StartTime, report.EndTime),
	}
	addTestSuite := func(test *model.TestReport) {
		testSuite := junit.Testsuite{
			Name:    test.Name,
			Package: test.BasePath,
			Time:    durationInSecondsString(test.StartTime, test.EndTime),
		}
		testSuite.SetTimestamp(report.StartTime)
		testSuite.AddProperty("namespace", test.Namespace)
		if test.Skipped {
			testCase := junit.Testcase{
				Name: test.Name,
				Time: durationInSecondsString(test.StartTime, test.EndTime),
			}
			testCase.Skipped = &junit.Result{}
			testSuite.AddTestcase(testCase)
		} else {
			for _, step := range test.Steps {
				testCase := junit.Testcase{
					Name: step.Name,
					Time: durationInSecondsString(step.StartTime, step.EndTime),
				}
				var errs []error
				var output strings.Builder
				for _, operation := range step.Operations {
					if operation.Err != nil {
						errs = append(errs, operation.Err)
					}
					output.WriteString(operation.Output)
				}
				testCase.SystemOut = systemOut(output.String())
				if err := multierr.Combine(errs...); err != nil {
					testCase.Failure = &junit.Result{
						Message: err.Error(),
					}
				}
				testSuite.AddTestcase(testCase)
			}
		}
		testSuites.AddSuite(testSuite)
	}
	for _, test := range report.Tests {
		addTestSuite(test)
	}
	data, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o600)
}

func saveJUnitOperation(report *model.Report, file string) error {
	testSuites := &junit.Testsuites{
		Name: report.Name,
		Time: durationInSecondsString(report.StartTime, report.EndTime),
	}
	addTestSuite := func(test *model.TestReport) {
		testSuite := junit.Testsuite{
			Name:    test.Name,
			Package: test.BasePath,
			Time:    durationInSecondsString(test.StartTime, test.EndTime),
		}
		testSuite.SetTimestamp(report.StartTime)
		testSuite.AddProperty("namespace", test.Namespace)
		if test.Skipped {
			testCase := junit.Testcase{
				Name: test.Name,
				Time: durationInSecondsString(test.StartTime, test.EndTime),
			}
			testCase.Skipped = &junit.Result{}
			testSuite.AddTestcase(testCase)
		} else {
			for _, step := range test.Steps {
				for _, operation := range step.Operations {
					testCase := junit.Testcase{
						Name:      fmt.Sprintf("%s / %s", step.Name, operation.Name),
						Classname: string(operation.Type),
						Time:      durationInSecondsString(operation.StartTime, operation.EndTime),
						SystemOut: systemOut(operation.Output),
					}
					if err := operation.Err; err != nil {
						testCase.Failure = &junit.Result{
							Message: err.Error(),
						}
					}
					testSuite.AddTestcase(testCase)
				}
			}
		}
		testSuites.AddSuite(testSuite)
	}
	for _, test := range report.Tests {
		addTestSuite(test)
	}
	data, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o600)
}
