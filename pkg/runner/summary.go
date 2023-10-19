package runner

import (
	"time"
)

type Summary struct {
	Duration     time.Duration
	PassedTests  int32
	FailedTests  int32
	SkippedTests int32
}
