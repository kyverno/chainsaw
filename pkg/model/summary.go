package model

import (
	"sync/atomic"
)

type SummaryResult interface {
	Passed() int32
	Failed() int32
	Skipped() int32
}

type Summary interface {
	SummaryResult
	IncPassed()
	IncFailed()
	IncSkipped()
}

type summary struct {
	passed  atomic.Int32
	failed  atomic.Int32
	skipped atomic.Int32
}

func (s *summary) IncPassed() {
	s.passed.Add(1)
}

func (s *summary) IncFailed() {
	s.failed.Add(1)
}

func (s *summary) IncSkipped() {
	s.skipped.Add(1)
}

func (s *summary) Passed() int32 {
	return s.passed.Load()
}

func (s *summary) Failed() int32 {
	return s.failed.Load()
}

func (s *summary) Skipped() int32 {
	return s.skipped.Load()
}
