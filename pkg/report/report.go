package report

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"time"
)

type Report struct {
	Tests []TestReport `json:"tests" xml:"tests"`
}

type TestReport struct {
	Name      string       `json:"name" xml:"name"`
	StartTime time.Time    `json:"startTime" xml:"startTime"`
	EndTime   time.Time    `json:"endTime" xml:"endTime"`
	Steps     []StepReport `json:"steps" xml:"steps"`
}

type StepReport struct {
	Name    string `json:"name" xml:"name"`
	Result  string `json:"result" xml:"result"`
	Message string `json:"message,omitempty" xml:"message,omitempty"`
}

// Function to serialize to JSON
func (r *Report) ToJSON() ([]byte, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return []byte(""), err
	}
	return data, nil
}

// Function to serialize to XML
func (r *Report) ToXML() ([]byte, error) {
	data, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return []byte(""), err
	}
	return data, nil
}

func (report *Report) SaveReportAsJSON(filePath string) error {
	jsonData, err := report.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0644)
}

func (report *Report) SaveReportAsXML(filePath string) error {
	xmlData, err := report.ToXML()
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, xmlData, 0644)
}
