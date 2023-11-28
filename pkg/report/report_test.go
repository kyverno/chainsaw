package report

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

type FakeSerializer struct{}

func (s FakeSerializer) Serialize(report *TestsReport) ([]byte, error) {
	return nil, errors.New("FakeSerializer error")
}

func TestSaveReport(t *testing.T) {
	testCases := []struct {
		name        string
		reportName  string
		serializer  ReportSerializer
		expectError bool
	}{
		{
			name:        "SuccessfulSaveJSON",
			reportName:  "test_report.json",
			serializer:  JSONSerializer{},
			expectError: false,
		},
		{
			name:        "SuccessfulSaveXML",
			reportName:  "test_report.xml",
			serializer:  XMLSerializer{},
			expectError: false,
		},
		{
			name:        "FailedSaveInvalidPathJSON",
			reportName:  "/invalid_path/test_report.json",
			serializer:  JSONSerializer{},
			expectError: true,
		},
		{
			name:        "FailedSaveInvalidPathXML",
			reportName:  "/invalid_path/test_report.xml",
			serializer:  XMLSerializer{},
			expectError: true,
		},
		{
			name:        "FakeSerializerError",
			reportName:  "test_report.json",
			serializer:  FakeSerializer{},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			report := NewTests("SampleTestSuite")
			report.AddTest(NewTest("Test1"))

			filePath := filepath.Join(t.TempDir(), tc.reportName)
			err := SaveReport(report, tc.serializer, filePath)

			if tc.expectError {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Expected no error")

				content, err := os.ReadFile(filePath)
				assert.NoError(t, err, "Failed to read saved file")
				assert.NotEmpty(t, content, "File should not be empty")
			}
		})
	}
}

func TestSaveReportBasedOnType(t *testing.T) {
	testCases := []struct {
		name        string
		format      v1alpha1.ReportFormatType
		fileSuffix  string
		expectError bool
	}{
		{
			name:        "SuccessfulSaveJSON",
			format:      v1alpha1.JSONFormat,
			fileSuffix:  "json",
			expectError: false,
		},
		{
			name:        "SuccessfulSaveXML",
			format:      v1alpha1.XMLFormat,
			fileSuffix:  "xml",
			expectError: false,
		},
		{
			name:        "UnsupportedFormat",
			format:      "Unsupported",
			fileSuffix:  "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			report := NewTests("SampleTestSuite")
			report.AddTest(NewTest("Test1"))

			reportName := filepath.Join(t.TempDir(), "test_report")
			err := report.SaveReportBasedOnType(tc.format, reportName)

			if tc.expectError {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Expected no error")

				filePath := reportName + "." + tc.fileSuffix
				content, err := os.ReadFile(filePath)
				assert.NoError(t, err, "Failed to read saved file")
				assert.NotEmpty(t, content, "File should not be empty")
			}
		})
	}
}

func TestAddTest(t *testing.T) {
	report := NewTests("SampleTestSuite")
	testReport := NewTest("Test1")
	report.AddTest(testReport)

	assert.Equal(t, 1, len(report.Reports), "Expected 1 test in report")
}

func TestAddTestStep(t *testing.T) {
	testReport := NewTest("Test1")
	step := NewTestSpecStep("Step1")

	testReport.AddTestStep(step)

	assert.Equal(t, 1, len(testReport.Steps), "Expected 1 step in test report")
	assert.Equal(t, step, testReport.Steps[0], "The added step does not match the expected step")
}

func TestAddOperation(t *testing.T) {
	testSpecStep := NewTestSpecStep("Step1")
	operation := NewOperation("Operation1", OperationTypeCreate)

	testSpecStep.AddOperation(operation)

	assert.Equal(t, 1, len(testSpecStep.Results), "Expected 1 operation in test spec step")
	assert.Equal(t, operation, testSpecStep.Results[0], "The added operation does not match the expected operation")
}

func TestNewFailure(t *testing.T) {
	testReport := NewTest("Test1")
	testReport.NewFailure("Sample failure message")

	assert.NotNil(t, testReport.Failure, "Failure object should not be nil")
	assert.Equal(t, "Sample failure message", testReport.Failure.Message, "Failure message does not match")
}

func TestMarkTestEnd(t *testing.T) {
	startTime := time.Now().Add(-10 * time.Second)
	testReport := &TestReport{
		TimeStamp: startTime,
		Steps:     []*TestSpecStepReport{{Results: make([]*OperationReport, 2)}},
	}

	testReport.MarkTestEnd()

	assert.Regexp(t, `\d+\.\d{3}`, testReport.Time, "Duration format is incorrect")
	assert.Equal(t, 2, testReport.Test, "Total tests count should be 2")
}

func TestMarkOperationEnd(t *testing.T) {
	testCases := []struct {
		name           string
		success        bool
		expectedResult string
		message        string
	}{
		{
			name:           "OperationSuccessful",
			success:        true,
			expectedResult: "Success",
			message:        "Completed successfully",
		},
		{
			name:           "OperationFailed",
			success:        false,
			expectedResult: "Failure",
			message:        "An error occurred",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startTime := time.Now().Add(-5 * time.Second) // mock start time 5 seconds ago
			operation := &OperationReport{TimeStamp: startTime}

			operation.MarkOperationEnd(tc.success, tc.message)

			assert.Regexp(t, `\d+\.\d{3}`, operation.Time, "Duration format is incorrect")
			assert.Equal(t, tc.expectedResult, operation.Result, "Result does not match expected value")
			assert.Equal(t, tc.message, operation.Message, "Message does not match")
		})
	}
}

func TestCalculateDuration(t *testing.T) {
	startTime := time.Now().Add(-30 * time.Second) // mock start time 30 seconds ago
	durationStr := calculateDuration(startTime, time.Now())

	assert.Regexp(t, `\d+\.\d{3}`, durationStr, "Duration format is incorrect")
}

func TestClose(t *testing.T) {
	startTime := time.Now().Add(-60 * time.Second)
	testsReport := &TestsReport{
		TimeStamp: startTime,
		Reports: []*TestReport{
			{Test: 1},
			{Test: 1, Failure: &Failure{}},
		},
	}

	testsReport.Close()

	assert.Regexp(t, `\d+\.\d{3}`, testsReport.Time, "Duration format is incorrect")
	assert.Equal(t, 1, testsReport.Failures, "Failures count should be 1")
	assert.Equal(t, 2, testsReport.Test, "Total tests count should be 2")
}
