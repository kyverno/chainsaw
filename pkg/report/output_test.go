package report

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jstemmer/go-junit-report/v2/junit"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func reportWithOutput() *model.Report {
	start := time.Date(2009, 11, 17, 20, 34, 58, 0, time.UTC)
	return &model.Report{
		Name:      "report",
		StartTime: start,
		EndTime:   start.Add(time.Minute),
		Tests: []*model.TestReport{{
			BasePath:  "base-path",
			Name:      "test",
			StartTime: start,
			EndTime:   start.Add(time.Minute),
			Steps: []*model.StepReport{{
				Name:      "step",
				StartTime: start,
				EndTime:   start.Add(time.Minute),
				Operations: []*model.OperationReport{{
					Name:      "apply the pod",
					Type:      model.OperationTypeApply,
					StartTime: start,
					EndTime:   start.Add(time.Second),
					Output:    "APPLY | DONE | v1/Pod @ ns/pod\n",
				}, {
					Name:      "operation 2",
					Type:      model.OperationTypeSleep,
					StartTime: start,
					EndTime:   start.Add(time.Second),
				}},
			}},
		}},
	}
}

func TestSave_operationNameAndOutput(t *testing.T) {
	path := t.TempDir()
	report := reportWithOutput()

	require.NoError(t, Save(report, v1alpha2.JUnitOperationFormat, path, "operation"))
	data, err := os.ReadFile(filepath.Join(path, "operation.xml"))
	require.NoError(t, err)
	var operation junit.Testsuites
	require.NoError(t, xml.Unmarshal(data, &operation))
	require.Len(t, operation.Suites, 1)
	cases := operation.Suites[0].Testcases
	require.Len(t, cases, 2)
	assert.Equal(t, "step / apply the pod", cases[0].Name)
	require.NotNil(t, cases[0].SystemOut)
	assert.Equal(t, "APPLY | DONE | v1/Pod @ ns/pod\n", cases[0].SystemOut.Data)
	assert.Equal(t, "step / operation 2", cases[1].Name)
	assert.Nil(t, cases[1].SystemOut)

	require.NoError(t, Save(report, v1alpha2.JUnitStepFormat, path, "step"))
	data, err = os.ReadFile(filepath.Join(path, "step.xml"))
	require.NoError(t, err)
	var step junit.Testsuites
	require.NoError(t, xml.Unmarshal(data, &step))
	require.Len(t, step.Suites, 1)
	require.Len(t, step.Suites[0].Testcases, 1)
	require.NotNil(t, step.Suites[0].Testcases[0].SystemOut)
	assert.Equal(t, "APPLY | DONE | v1/Pod @ ns/pod\n", step.Suites[0].Testcases[0].SystemOut.Data)

	require.NoError(t, Save(report, v1alpha2.JSONFormat, path, "json"))
	data, err = os.ReadFile(filepath.Join(path, "json.json"))
	require.NoError(t, err)
	var parsed struct {
		Tests []struct {
			Steps []struct {
				Operations []map[string]any `json:"operations"`
			} `json:"steps"`
		} `json:"tests"`
	}
	require.NoError(t, json.Unmarshal(data, &parsed))
	operations := parsed.Tests[0].Steps[0].Operations
	require.Len(t, operations, 2)
	assert.Equal(t, "apply the pod", operations[0]["name"])
	assert.Equal(t, "APPLY | DONE | v1/Pod @ ns/pod\n", operations[0]["output"])
	assert.NotContains(t, operations[1], "output")
}
