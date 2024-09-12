package model

import (
	"fmt"
	"sync"
	"time"
)

type Report struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Tests     []TestReport
	lock      sync.Mutex
}

func (r *Report) Add(report TestReport) {
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
	Failed     bool
	Steps      []StepReport
}

func (r *TestReport) Add(report StepReport) {
	if report.Name == "" {
		report.Name = fmt.Sprintf("step %d", len(r.Steps)+1)
	}
	if report.Failed {
		r.Failed = true
	}
	r.Steps = append(r.Steps, report)
}

type StepReport struct {
	Name       string
	StartTime  time.Time
	EndTime    time.Time
	Failed     bool
	Operations []OperationReport
}

func (r *StepReport) Add(report OperationReport) {
	if report.Name == "" {
		report.Name = fmt.Sprintf("operation %d", len(r.Operations)+1)
	}
	if report.Err != nil {
		r.Failed = true
	}
	r.Operations = append(r.Operations, report)
}

type OperationReport struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Err       error
}
