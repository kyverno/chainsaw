package report

import (
	"encoding/xml"
	"os"
	"time"

	"github.com/kyverno/chainsaw/pkg/model"
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

func saveJUnit(report *model.Report, file string) error {
	testsuites := testsuitesNode{
		Name:      report.Name,
		Timestamp: report.StartTime.UTC().Format(time.RFC3339),
		Time:      report.EndTime.Sub(report.StartTime).Seconds(),
	}
	perFolder := map[string][]model.TestReport{}
	for _, test := range report.Tests {
		perFolder[test.BasePath] = append(perFolder[test.BasePath], test)
	}
	for folder, tests := range perFolder {
		testsuite := testsuiteNode{
			Name: folder,
		}
		for _, test := range tests {
			var properties []any
			if test.Namespace != "" {
				properties = append(properties, propertyNode{
					Name:  "namespace",
					Value: test.Namespace,
				})
			}
			// for i, step := range test.steps {
			// 	if step.step.Name != "" {
			// 		properties = append(properties, propertyNode{
			// 			Name:  fmt.Sprintf("step%d", i),
			// 			Value: step.step.Name,
			// 		})
			// 	}
			// 	for j, op := range step.reports {
			// 		if op.err != nil {
			// 			properties = append(properties, propertyNode{
			// 				Name:  fmt.Sprintf("step%d", i),
			// 				Value: fmt.Sprintf("op %d - %s: %s", j, op.operationType, op.err),
			// 			})
			// 		} else {
			// 			properties = append(properties, propertyNode{
			// 				Name:  fmt.Sprintf("step%d", i),
			// 				Value: fmt.Sprintf("op %d - %s", j, op.operationType),
			// 			})
			// 		}
			// 	}
			// }
			testcase := testcaseNode{
				Name:      test.Name,
				Timestamp: test.StartTime.UTC().Format(time.RFC3339),
				Time:      test.EndTime.Sub(test.StartTime).Seconds(),
				File:      test.BasePath,
			}
			if len(properties) != 0 {
				testcase.Inner = append(testcase.Inner, propertiesNode{Inner: properties})
			}
			if test.Skipped {
				testcase.Inner = append(testcase.Inner, skippedNode{})
			}
			if test.Failed {
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
