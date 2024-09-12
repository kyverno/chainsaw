package model

import (
	"sync"
	"time"
)

type Report struct {
	Tests []TestReport
	lock  sync.Mutex
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
	Failed     bool
	Skipped    bool
	Steps      []StepReport
}

type StepReport struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
}
