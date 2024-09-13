package model

import (
	"fmt"
	"sync"
	"time"
)

type OperationType string

const (
	OperationTypeApply   OperationType = "apply"
	OperationTypeAssert  OperationType = "assert"
	OperationTypeCommand OperationType = "command"
	OperationTypeCreate  OperationType = "create"
	OperationTypeDelete  OperationType = "delete"
	OperationTypeError   OperationType = "error"
	OperationTypePatch   OperationType = "patch"
	OperationTypeScript  OperationType = "script"
	OperationTypeSleep   OperationType = "sleep"
	OperationTypeUpdate  OperationType = "update"
)

type Report struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Tests     []*TestReport
	lock      sync.Mutex
}

func (r *Report) Add(report *TestReport) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Tests = append(r.Tests, report)
}

type TestReport struct {
	BasePath   string
	Name       string
	Concurrent *bool
	StartTime  time.Time
	EndTime    time.Time
	Namespace  string
	Skipped    bool
	Steps      []*StepReport
}

func (r *TestReport) Add(report *StepReport) {
	if report.Name == "" {
		report.Name = fmt.Sprintf("step %d", len(r.Steps)+1)
	}
	r.Steps = append(r.Steps, report)
}

type StepReport struct {
	Name       string
	StartTime  time.Time
	EndTime    time.Time
	Operations []*OperationReport
}

func (r *StepReport) Add(report *OperationReport) {
	if report.Name == "" {
		report.Name = fmt.Sprintf("operation %d", len(r.Operations)+1)
	}
	r.Operations = append(r.Operations, report)
}

func (r *StepReport) Failed() bool {
	for _, operation := range r.Operations {
		if operation.Err != nil {
			return true
		}
	}
	return false
}

type OperationReport struct {
	Name      string
	Type      OperationType
	StartTime time.Time
	EndTime   time.Time
	Err       error
}
