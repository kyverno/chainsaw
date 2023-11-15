package report

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"strings"
	"time"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

type ReportSerializer interface {
	Serialize(report *Report) ([]byte, error)
}

type Report struct {
	Tests []TestsReport `json:"tests" xml:"tests"`
}

type TestsReport struct {
	Name      string       `json:"name" xml:"name"`
	StartTime time.Time    `json:"startTime" xml:"startTime"`
	EndTime   time.Time    `json:"endTime" xml:"endTime"`
	Steps     []TestReport `json:"steps" xml:"steps"`
}

type TestReport struct {
	Name      string       `json:"name" xml:"name"`
	StartTime time.Time    `json:"startTime" xml:"startTime"`
	EndTime   time.Time    `json:"endTime" xml:"endTime"`
	Steps     []StepReport `json:"steps" xml:"steps"`
}

type StepReport struct {
	Name      string            `json:"name" xml:"name"`
	StartTime time.Time         `json:"startTime" xml:"startTime"`
	EndTime   time.Time         `json:"endTime" xml:"endTime"`
	Steps     []OperationReport `json:"steps" xml:"steps"`
}

type OperationReport struct {
	Name    string `json:"name" xml:"name"`
	Result  string `json:"result" xml:"result"`
	Message string `json:"message,omitempty" xml:"message,omitempty"`
}

type JSONSerializer struct{}

func (s JSONSerializer) Serialize(report *Report) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

type XMLSerializer struct{}

func (s XMLSerializer) Serialize(report *Report) ([]byte, error) {
	return xml.MarshalIndent(report, "", "  ")
}

func SaveReport(report *Report, serializer ReportSerializer, filePath string) error {
	data, err := serializer.Serialize(report)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
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

func (report *Report) SaveReportBasedOnType(reportFormat v1alpha1.ReportFormatType, reportName string) error {
	serializer, err := GetSerializer(reportFormat)
	if err != nil {
		return err
	}
	filePath := reportName + "." + strings.ToLower(string(reportFormat))
	return SaveReport(report, serializer, filePath)
}
