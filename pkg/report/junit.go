package report

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/jstemmer/go-junit-report/v2/junit"
	"github.com/kyverno/chainsaw/pkg/model"
)

const (
	FailureTypeAssertionError = "AssertionError"
)

// durationInSecondsString is a helper function to convert a start and end time into a string
// representing the duration in seconds. This is needed by the junit package for generating
// the JUnit XML report.
func durationInSecondsString(start, end time.Time) string {
	duration := end.Sub(start)
	return fmt.Sprintf("%.6f", duration.Seconds())
}

// addTestSuite loops through all the operations reports for each directory that is
// being tested and adds them to the JUnit XML report. This is done by looping through
// all tests and then looping through all of their steps and all the reports for all steps.
//
// The end goal is to have each Test file represented as a TestSuite and each step as a TestCase.
func addTestSuite(testSuites *junit.Testsuites, report *model.TestReport) {
	// Loop through all the Tests in the report
	suite := junit.Testsuite{
		Name:    report.Name,
		Package: report.BasePath,
		Time:    durationInSecondsString(report.StartTime, report.EndTime),
	}
	suite.SetTimestamp(report.StartTime)
	suite.AddProperty("namespace", report.Namespace)
	if report.Skipped {
		suite.Skipped = suite.Skipped + 1
	}
	for _, report := range report.Steps {
		testCase := junit.Testcase{
			Name: report.Name,
			Time: durationInSecondsString(report.StartTime, report.EndTime),
		}
		if report.Failed {
			testCase.Failure = &junit.Result{}
			for _, report := range report.Operations {
				if report.Err != nil {
					testCase.Failure.Message = report.Err.Error()
					break
				}
			}
		}
		suite.AddTestcase(testCase)
	}
	testSuites.AddSuite(suite)
}

// saveJUnit writes the JUnit XML report to disk. The spec is defined here:
// https://github.com/testmoapp/junitxml
//
// This method makes use of https://github.com/jstemmer/go-junit-report to generate the XML.
func saveJUnit(report *model.Report, file string) error {
	// Initialize the top-level TestSuites object
	testSuites := &junit.Testsuites{
		Name: report.Name,
		Time: durationInSecondsString(report.StartTime, report.EndTime),
	}
	// Append the individual test suites to the parent
	for _, test := range report.Tests {
		addTestSuite(testSuites, test)
	}
	data, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o600)
}
