package report

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func Test_getFile(t *testing.T) {
	tests := []struct {
		path      string
		name      string
		extension string
		want      string
	}{{
		path:      "foo",
		name:      "bar",
		extension: "xml",
		want:      "foo/bar.xml",
	}, {
		path:      "foo",
		name:      "bar.json",
		extension: "xml",
		want:      "foo/bar.json",
	}, {
		path:      "foo",
		name:      "bar.xml",
		extension: "xml",
		want:      "foo/bar.xml",
	}, {
		path:      "",
		name:      "bar",
		extension: "xml",
		want:      "bar.xml",
	}, {
		path:      "",
		name:      "bar.json",
		extension: "xml",
		want:      "bar.json",
	}, {
		path:      "",
		name:      "bar.xml",
		extension: "xml",
		want:      "bar.xml",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFile(tt.path, tt.name, tt.extension)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSave(t *testing.T) {
	report := &model.Report{
		Name:      "report",
		StartTime: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		EndTime:   time.Date(2009, 11, 17, 20, 44, 58, 651387237, time.UTC),
		Tests: []*model.TestReport{{
			BasePath:   "base-path",
			Name:       "test-report",
			Concurrent: nil,
			StartTime:  time.Date(2009, 11, 17, 20, 35, 58, 651387237, time.UTC),
			EndTime:    time.Date(2009, 11, 17, 20, 38, 58, 651387237, time.UTC),
			Namespace:  "demo",
			Skipped:    false,
			Steps: []*model.StepReport{{
				Name:      "step-report",
				StartTime: time.Date(2009, 11, 17, 20, 36, 58, 651387237, time.UTC),
				EndTime:   time.Date(2009, 11, 17, 20, 37, 58, 651387237, time.UTC),
				Operations: []*model.OperationReport{{
					Name:      "operation-report",
					Type:      model.OperationTypeApply,
					StartTime: time.Date(2009, 11, 17, 20, 36, 58, 651387237, time.UTC),
					EndTime:   time.Date(2009, 11, 17, 20, 37, 58, 651387237, time.UTC),
					Err:       nil,
				}, {
					Name:      "operation-report-err",
					Type:      model.OperationTypeAssert,
					StartTime: time.Date(2009, 11, 17, 20, 36, 58, 651387237, time.UTC),
					EndTime:   time.Date(2009, 11, 17, 20, 37, 58, 651387237, time.UTC),
					Err:       errors.New("dummy"),
				}},
			}},
		}, {
			BasePath:   "base-path",
			Name:       "test-report-skipped",
			Concurrent: ptr.To(true),
			StartTime:  time.Date(2009, 11, 17, 20, 35, 58, 651387237, time.UTC),
			EndTime:    time.Date(2009, 11, 17, 20, 38, 58, 651387237, time.UTC),
			Namespace:  "skipped",
			Skipped:    true,
		}},
	}
	tests := []struct {
		name    string
		report  *model.Report
		format  v1alpha2.ReportFormatType
		wantErr bool
		out     string
	}{{
		report: report,
		format: v1alpha2.JSONFormat,
		out:    "JSON.json",
	}, {
		report: report,
		format: v1alpha2.XMLFormat,
		out:    "XML.xml",
	}, {
		report: report,
		format: v1alpha2.JUnitTestFormat,
		out:    "JUNIT-TEST.xml",
	}, {
		report: report,
		format: v1alpha2.JUnitStepFormat,
		out:    "JUNIT-STEP.xml",
	}, {
		report: report,
		format: v1alpha2.JUnitOperationFormat,
		out:    "JUNIT-OPERATION.xml",
	}, {
		report:  report,
		format:  v1alpha2.ReportFormatType("xyz"),
		wantErr: true,
	}}
	path, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Save(tt.report, tt.format, path, string(tt.format))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.out != "" {
					expected, err := os.ReadFile(filepath.Join(path, tt.out))
					assert.NoError(t, err)
					actual, err := os.ReadFile(filepath.Join(path, tt.out))
					assert.NoError(t, err)
					assert.Equal(t, string(expected), string(actual))
				}
			}
		})
	}
}
