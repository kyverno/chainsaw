package report

import (
	"encoding/xml"
	"fmt"
	"github.com/jstemmer/go-junit-report/v2/junit"
	"os"
	"time"
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
func addTestSuite(testSuites *junit.Testsuites, report *Report) {
	// Loop through all the Tests in the report
	for _, test := range report.tests {
		// Initialize the TestSuite for this test
		suite := junit.Testsuite{
			Name:      test.test.Name,                              // Name is pulled args to the command (--report-name)
			Timestamp: report.startTime.UTC().Format(time.RFC3339), // Take the time from the beginning of the test
			Time:      durationInSecondsString(test.startTime, test.endTime),
		}

		// Loop through all the Steps in the Test
		for _, step := range test.steps {
			// Loop through each Report
			for _, report := range step.reports {
				// Each report is now an individual Testcase because it is a single operation executed
				// against a cluster and can fail on its own.
				testCase := junit.Testcase{
					Name:      report.name,
					Classname: suite.Name, // Associate the Testcase with the TestSuite
					Time:      durationInSecondsString(report.startTime, report.endTime),
				}

				// Each report will have an error property set if the step failed.
				if report.err != nil {
					testCase.Failure = &junit.Result{
						Message: report.err.Error(),
					}

					// Type hints can help identify what the error was. This is not required but can be helpful.
					// TODO: add more type hints
					// https://github.com/testmoapp/junitxml?tab=readme-ov-file#conventions-types
					switch report.operationType {
					case OperationTypeAssert:
						testCase.Failure.Type = FailureTypeAssertionError
					default:
						testCase.Failure.Type = "" // Do not set if we don't know what to categorize it as
					}
				}

				// Add each Testcase to the parent suite in order to increment the count of failures/successes/etc.
				suite.AddTestcase(testCase)
			}
		}

		// Add the test suite to the parent
		testSuites.AddSuite(suite)
	}
}

// saveJUnit writes the JUnit XML report to disk. The spec is defined here:
// https://github.com/testmoapp/junitxml
//
// This method makes use of https://github.com/jstemmer/go-junit-report to generate the XML.
func saveJUnit(report *Report, file string) error {
	// Initialize the top-level TestSuites object
	testSuites := &junit.Testsuites{
		Name:   report.name,
		Time:   durationInSecondsString(report.startTime, report.endTime),
		Suites: []junit.Testsuite{},
	}

	// Append the individual test suites to the parent
	addTestSuite(testSuites, report)

	data, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0o600)
}
