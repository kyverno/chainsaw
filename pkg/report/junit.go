package report

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type failureNode struct {
	XMLName struct{} `xml:"failure"`
}

type skippedNode struct {
	XMLName struct{} `xml:"skipped"`
}

type propertyNode struct {
	XMLName struct{} `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr,omitempty"`
}

type propertiesNode struct {
	XMLName struct{} `xml:"properties"`
	Inner   []any
}

type testcaseNode struct {
	XMLName   struct{} `xml:"testcase"`
	Name      string   `xml:"name,attr,omitempty"`
	Timestamp string   `xml:"timestamp,attr,omitempty"`
	Time      float64  `xml:"time,attr,omitempty"`
	File      string   `xml:"file,attr,omitempty"`
	Inner     []any
}

type testsuiteNode struct {
	XMLName   struct{} `xml:"testsuite"`
	Name      string   `xml:"name,attr,omitempty"`
	Timestamp string   `xml:"timestamp,attr,omitempty"`
	File      string   `xml:"file,attr,omitempty"`
	Inner     []any
}

type testsuitesNode struct {
	XMLName   struct{} `xml:"testsuites"`
	Name      string   `xml:"name,attr,omitempty"`
	Timestamp string   `xml:"timestamp,attr,omitempty"`
	Time      float64  `xml:"time,attr,omitempty"`
	Inner     []any
}

func saveJUnit(report *Report, file string) error {
	testsuites := testsuitesNode{
		Name:      report.name,
		Timestamp: report.startTime.UTC().Format(time.RFC3339),
		Time:      report.endTime.Sub(report.startTime).Seconds(),
	}
	perFolder := map[string][]*TestReport{}
	for _, test := range report.tests {
		perFolder[test.test.BasePath] = append(perFolder[test.test.BasePath], test)
	}
	for folder, tests := range perFolder {
		testsuite := testsuiteNode{
			Name: folder,
		}
		for _, test := range tests {
			var properties []any
			if test.namespace != "" {
				properties = append(properties, propertyNode{
					Name:  "namespace",
					Value: test.namespace,
				})
			}
			for i, step := range test.steps {
				if step.step.Name != "" {
					properties = append(properties, propertyNode{
						Name:  fmt.Sprintf("step%d", i),
						Value: step.step.Name,
					})
				}
				for j, op := range step.reports {
					if op.err != nil {
						properties = append(properties, propertyNode{
							Name:  fmt.Sprintf("step%d", i),
							Value: fmt.Sprintf("op %d - %s: %s", j, op.operationType, op.err),
						})
					} else {
						properties = append(properties, propertyNode{
							Name:  fmt.Sprintf("step%d", i),
							Value: fmt.Sprintf("op %d - %s", j, op.operationType),
						})
					}
				}
			}
			testcase := testcaseNode{
				Name:      test.test.Name,
				Timestamp: test.startTime.UTC().Format(time.RFC3339),
				Time:      test.endTime.Sub(test.startTime).Seconds(),
				File:      test.test.BasePath,
			}
			if len(properties) != 0 {
				testcase.Inner = append(testcase.Inner, propertiesNode{Inner: properties})
			}
			if test.skipped {
				testcase.Inner = append(testcase.Inner, skippedNode{})
			}
			if test.failed {
				testcase.Inner = append(testcase.Inner, failureNode{})
			}
			testsuite.Inner = append(testsuite.Inner, testcase)
		}
		testsuites.Inner = append(testsuites.Inner, testsuite)
	}
	data, err := xml.MarshalIndent(testsuites, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o600)
}
