package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport_Add(t *testing.T) {
	tests := []struct {
		name    string
		tests   []*TestReport
		report  *TestReport
		wantLen int
	}{{
		name:    "add nil",
		wantLen: 0,
	}, {
		name:    "add not nil",
		report:  &TestReport{},
		wantLen: 1,
	}, {
		name:    "add not nil",
		tests:   []*TestReport{{}, {}, {}},
		report:  &TestReport{},
		wantLen: 4,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Report{
				Tests: tt.tests,
			}
			r.Add(tt.report)
			assert.Equal(t, tt.wantLen, len(r.Tests))
			if tt.report != nil {
				assert.Equal(t, tt.report, r.Tests[len(r.Tests)-1])
			}
		})
	}
}

func TestTestReport_Add(t *testing.T) {
	tests := []struct {
		name    string
		steps   []*StepReport
		report  *StepReport
		wantLen int
	}{{
		name:    "add nil",
		wantLen: 0,
	}, {
		name:    "add not nil",
		report:  &StepReport{},
		wantLen: 1,
	}, {
		name:    "add not nil",
		steps:   []*StepReport{{}, {}, {}},
		report:  &StepReport{},
		wantLen: 4,
	}, {
		name:    "add not nil",
		steps:   []*StepReport{{}, {}, {}},
		report:  &StepReport{Name: "foo"},
		wantLen: 4,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TestReport{
				Steps: tt.steps,
			}
			r.Add(tt.report)
			assert.Equal(t, tt.wantLen, len(r.Steps))
			if tt.report != nil {
				assert.Equal(t, tt.report, r.Steps[len(r.Steps)-1])
				assert.NotEmpty(t, r.Steps[len(r.Steps)-1].Name)
			}
		})
	}
}

func TestStepReport_Add(t *testing.T) {
	tests := []struct {
		name       string
		operations []*OperationReport
		report     *OperationReport
		wantLen    int
	}{{
		name:    "add nil",
		wantLen: 0,
	}, {
		name:    "add not nil",
		report:  &OperationReport{},
		wantLen: 1,
	}, {
		name:       "add not nil",
		operations: []*OperationReport{{}, {}, {}},
		report:     &OperationReport{},
		wantLen:    4,
	}, {
		name:       "add not nil",
		operations: []*OperationReport{{}, {}, {}},
		report:     &OperationReport{Name: "foo"},
		wantLen:    4,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &StepReport{
				Operations: tt.operations,
			}
			r.Add(tt.report)
			assert.Equal(t, tt.wantLen, len(r.Operations))
			if tt.report != nil {
				assert.Equal(t, tt.report, r.Operations[len(r.Operations)-1])
				assert.NotEmpty(t, r.Operations[len(r.Operations)-1].Name)
			}
		})
	}
}

func TestStepReport_Failed(t *testing.T) {
	tests := []struct {
		name   string
		report StepReport
		want   bool
	}{{
		name:   "nil",
		report: StepReport{},
	}, {
		name: "empty",
		report: StepReport{
			Operations: []*OperationReport{},
		},
	}, {
		name: "failed",
		report: StepReport{
			Operations: []*OperationReport{{}, {Err: errors.New("dummy")}, {}},
		},
		want: true,
	}, {
		name: "not failed",
		report: StepReport{
			Operations: []*OperationReport{{}, {}, {}},
		},
		want: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.report.Failed()
			assert.Equal(t, tt.want, got)
		})
	}
}
