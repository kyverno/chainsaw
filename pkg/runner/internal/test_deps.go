package internal

import (
	"io"
	"reflect"
	"regexp"
	"time"
)

// TestDeps implements the TestDeps interface for MainStart.
type TestDeps struct {
	Test     bool
	matchPat string
	matchRe  *regexp.Regexp
}

func (d *TestDeps) MatchString(pat, str string) (bool, error) {
	if d.Test {
		return true, nil
	}
	if d.matchRe == nil || d.matchPat != pat {
		d.matchPat = pat
		matchRe, err := regexp.Compile(d.matchPat)
		if err != nil {
			return false, err
		}
		d.matchRe = matchRe
	}
	return d.matchRe.MatchString(str), nil
}

func (*TestDeps) SetPanicOnExit0(bool) {}

func (*TestDeps) StartCPUProfile(w io.Writer) error {
	return nil
}

func (*TestDeps) StopCPUProfile() {
}

func (*TestDeps) WriteProfileTo(name string, w io.Writer, debug int) error {
	return nil
}

func (*TestDeps) ImportPath() string {
	return ""
}

func (*TestDeps) StartTestLog(w io.Writer) {}

func (*TestDeps) StopTestLog() error {
	return nil
}

func (*TestDeps) CoordinateFuzzing(time.Duration, int64, time.Duration, int64, int, []corpusEntry, []reflect.Type, string, string) error {
	return nil
}

func (*TestDeps) RunFuzzWorker(func(corpusEntry) error) error {
	return nil
}

func (*TestDeps) ReadCorpus(string, []reflect.Type) ([]corpusEntry, error) {
	return nil, nil
}

func (*TestDeps) CheckCorpus([]any, []reflect.Type) error {
	return nil
}

func (*TestDeps) ResetCoverage() {}

func (*TestDeps) SnapshotCoverage() {}

func (*TestDeps) InitRuntimeCoverage() (mode string, tearDown func(string, string) (string, error), snapcov func() float64) {
	return
}
