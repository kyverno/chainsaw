package report

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

type OperationType string

const (
	OperationTypeCreate  OperationType = "create"
	OperationTypeDelete  OperationType = "delete"
	OperationTypeApply   OperationType = "apply"
	OperationTypeAssert  OperationType = "assert"
	OperationTypeError   OperationType = "error"
	OperationTypeScript  OperationType = "script"
	OperationTypeCommand OperationType = "command"
)

type ReportSerializer interface {
	Serialize(report *TestsReport) ([]byte, error)
}

// Failure represents details of a test failure.
type Failure struct {
	// Message provides a summary of the failure.
	Message string `json:"message" xml:"message,attr"`
	// Type indicates the type of failure.
	Type string `json:"type" xml:"type,attr"`
}

// TestsReport encapsulates the entire report for a test suite.
type TestsReport struct {
	// Name of the test suite.
	Name string `json:"name" xml:"name,attr"`
	// TimeStamp marks when the test suite began execution.
	TimeStamp time.Time `json:"timestamp" xml:"timestamp,attr"`
	// Time indicates the total duration of the test suite.
	Time string `json:"time" xml:"time,attr"`
	// Test count the number of tests in the files/TestReports.
	Test int `json:"tests" xml:"tests,attr"`
	// Reports is an array of individual test reports within this suite.
	Reports []*TestReport `json:"testsuite" xml:"testsuite"`
	// Failures count the number of failed tests in the suite.
	Failures int `json:"failures" xml:"failures,attr"`
}

// TestReport represents a report for a single test.
type TestReport struct {
	// Name of the test.
	Name string `json:"name" xml:"name,attr"`
	// TimeStamp marks when the test began execution.
	TimeStamp time.Time `json:"timestamp" xml:"timestamp,attr"`
	// Time indicates the total duration of the test.
	Time string `json:"time" xml:"time,attr"`
	// Failure captures details if the test failed it should be nil otherwise.
	Failure *Failure `json:"failure,omitempty" xml:"failure,omitempty"`
	// Test count the number of tests in the suite/TestReport.
	Test int `json:"tests" xml:"tests,attr"`
	// Spec represents the specifications of the test.
	Steps []*TestSpecStepReport `json:"testcase,omitempty" xml:"testcase,omitempty"`
	// Concurrent indicates if the test runs concurrently with other tests.
	Concurrent bool `json:"concurrent,omitempty" xml:"concurrent,attr,omitempty"`
	// Namespace in which the test runs.
	Namespace string `json:"namespace,omitempty" xml:"namespace,attr,omitempty"`
	// Skip indicates if the test is skipped.
	Skip bool `json:"skip,omitempty" xml:"skip,attr,omitempty"`
	// SkipDelete indicates if resources are not deleted after test execution.
	SkipDelete bool `json:"skipDelete,omitempty" xml:"skipDelete,attr,omitempty"`
}

// TestSpecStepReport represents a report of a single step in a test.
type TestSpecStepReport struct {
	// Name of the test step.
	Name string `json:"name,omitempty" xml:"name,attr,omitempty"`
	// Results are the outcomes of operations performed in this step.
	Results []*OperationReport `json:"results,omitempty" xml:"results,omitempty"`
}

// OperationReport details the outcome of a single operation within a test step.
type OperationReport struct {
	// Name of the operation.
	Name string `json:"name" xml:"name,attr"`
	// TimeStamp marks when the operation began execution.
	TimeStamp time.Time `json:"timestamp" xml:"timestamp,attr"`
	// Time indicates the total duration of the operation.
	Time string `json:"time" xml:"time,attr"`
	// Result of the operation.
	Result string `json:"result" xml:"result,attr"`
	// Message provides additional information about the operation's outcome.
	Message string `json:"message,omitempty" xml:"message,omitempty"`
	// Type indicates the type of operation.
	OperationType OperationType `json:"operationType,omitempty" xml:"operationType,attr"`
}

type JSONSerializer struct{}

func (s JSONSerializer) Serialize(report *TestsReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

type XMLSerializer struct{}

func (s XMLSerializer) Serialize(report *TestsReport) ([]byte, error) {
	return xml.MarshalIndent(report, "", "  ")
}

func SaveReport(report *TestsReport, serializer ReportSerializer, filePath string) error {
	data, err := serializer.Serialize(report)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0o600)
}

func GetSerializer(format v1alpha1.ReportFormatType) (ReportSerializer, error) {
	switch format {
	case v1alpha1.JSONFormat:
		return JSONSerializer{}, nil
	case v1alpha1.XMLFormat:
		return XMLSerializer{}, nil
	default:
		return nil, errors.New("unsupported report format")
	}
}

func (report *TestsReport) SaveReportBasedOnType(reportFormat v1alpha1.ReportFormatType, reportName string) error {
	serializer, err := GetSerializer(reportFormat)
	if err != nil {
		return err
	}
	filePath := reportName + "." + strings.ToLower(string(reportFormat))
	return SaveReport(report, serializer, filePath)
}

// NewTests initializes a new TestsReport with the given name.
func NewTests(name string) *TestsReport {
	return &TestsReport{
		Name:      name,
		TimeStamp: time.Now(),
		Reports:   []*TestReport{},
	}
}

// NewTest creates a new TestReport with the given name.
func NewTest(name string, concurrent bool, namespace string, skip bool, skipDelete bool) *TestReport {
	return &TestReport{
		Name:       name,
		TimeStamp:  time.Now(),
		Concurrent: concurrent,
		Namespace:  namespace,
		Skip:       skip,
		SkipDelete: skipDelete,
		Steps:      []*TestSpecStepReport{},
	}
}

// NewTestSpecStep initializes a new TestSpecStepReport with the given name.
func NewTestSpecStep(name string) *TestSpecStepReport {
	return &TestSpecStepReport{
		Name:    name,
		Results: []*OperationReport{},
	}
}

// NewOperation creates a new OperationReport with the given details.
func NewOperation(name string, operationType OperationType) *OperationReport {
	return &OperationReport{
		Name:          name,
		TimeStamp:     time.Now(),
		OperationType: operationType,
	}
}

// AddTest adds a test report to the TestsReport.
func (tr *TestsReport) AddTest(test *TestReport) {
	tr.Reports = append(tr.Reports, test)
}

// AddTestStep adds a test step report to the TestReport.
func (t *TestReport) AddTestStep(step *TestSpecStepReport) {
	t.Steps = append(t.Steps, step)
}

// AddOperation adds an operation report to the TestSpecStepReport.
func (ts *TestSpecStepReport) AddOperation(op *OperationReport) {
	ts.Results = append(ts.Results, op)
}

// NewFailure creates a new Failure instance with the given message and type and assigns it to the TestReport.
func (t *TestReport) NewFailure(message, failureType string) {
	if t.Failure == nil {
		t.Failure = &Failure{
			Message: message,
			Type:    failureType,
		}
	}
}

// MarkTestEnd marks the end time of a TestReport and calculates its duration.
func (t *TestReport) MarkTestEnd() {
	t.Time = calculateDuration(t.TimeStamp, time.Now())

	for _, step := range t.Steps {
		t.Test += len(step.Results)
	}
}

// MarkOperationEnd marks the end time of an OperationReport and calculates its duration.
func (op *OperationReport) MarkOperationEnd(success bool, message string) {
	op.Time = calculateDuration(op.TimeStamp, time.Now())
	if success {
		op.Result = "Success"
	} else {
		op.Result = "Failure"
	}
	op.Message = message
}

// calculateDuration calculates the duration between two time points.
func calculateDuration(start, end time.Time) string {
	return fmt.Sprintf("%.3f", end.Sub(start).Seconds())
}

// Close finalizes the TestsReport, marking its end time and calculating the overall duration.
func (tr *TestsReport) Close() {
	tr.Time = calculateDuration(tr.TimeStamp, time.Now())
	totalTests := 0
	for _, test := range tr.Reports {
		if test.Failure != nil {
			tr.Failures++
		}
		totalTests += test.Test
	}
	tr.Test = totalTests
}
