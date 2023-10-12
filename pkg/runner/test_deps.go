package runner

import (
	"io"
	"reflect"
	"regexp"
	"runtime/pprof"
	"time"
)

// testDeps implements the testDeps interface for MainStart.
type testDeps struct {
	matchPat string
	matchRe  *regexp.Regexp
}

func (d *testDeps) MatchString(pat, str string) (bool, error) {
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

func (*testDeps) SetPanicOnExit0(bool) {}

func (*testDeps) StartCPUProfile(w io.Writer) error {
	return pprof.StartCPUProfile(w)
}

func (*testDeps) StopCPUProfile() {
	pprof.StopCPUProfile()
}

func (*testDeps) WriteProfileTo(name string, w io.Writer, debug int) error {
	return pprof.Lookup(name).WriteTo(w, debug)
}

func (*testDeps) ImportPath() string {
	return ""
}

func (*testDeps) StartTestLog(w io.Writer) {}

func (*testDeps) StopTestLog() error {
	return nil
}

func (*testDeps) CoordinateFuzzing(time.Duration, int64, time.Duration, int64, int, []corpusEntry, []reflect.Type, string, string) error {
	return nil
}

func (*testDeps) RunFuzzWorker(func(corpusEntry) error) error {
	return nil
}

func (*testDeps) ReadCorpus(string, []reflect.Type) ([]corpusEntry, error) {
	return nil, nil
}

func (*testDeps) CheckCorpus([]interface{}, []reflect.Type) error {
	return nil
}

func (*testDeps) ResetCoverage() {}

func (*testDeps) SnapshotCoverage() {}
