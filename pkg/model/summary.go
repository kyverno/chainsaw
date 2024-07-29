package model

import (
	"sync/atomic"
)

type SummaryResult interface {
	Passed() int32
	Failed() int32
	Skipped() int32
}

type Summary struct {
	passed  atomic.Int32
	failed  atomic.Int32
	skipped atomic.Int32
}

func (s *Summary) IncPassed() {
	s.passed.Add(1)
}

func (s *Summary) IncFailed() {
	s.failed.Add(1)
}

func (s *Summary) IncSkipped() {
	s.skipped.Add(1)
}

func (s *Summary) Passed() int32 {
	return s.passed.Load()
}

func (s *Summary) Failed() int32 {
	return s.failed.Load()
}

func (s *Summary) Skipped() int32 {
	return s.skipped.Load()
}
