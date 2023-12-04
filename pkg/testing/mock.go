package testing

import "time"

type MockT struct {
	NameVar           string
	FailedVar         bool
	ImmeditateFailVar bool
	SkippedVar        bool
}

func (c *MockT) Cleanup(f func()) {
}

func (t *MockT) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (c *MockT) Error(args ...any) {
}

func (c *MockT) Errorf(format string, args ...any) {
}

func (c *MockT) Fail() {
	c.FailedVar = true
}

func (c *MockT) FailNow() {
	c.ImmeditateFailVar = true
	c.FailedVar = true
}

func (c *MockT) Failed() bool {
	return c.FailedVar
}

func (c *MockT) Fatal(args ...any) {
}

func (c *MockT) Fatalf(format string, args ...any) {
}

func (c *MockT) Helper() {
}

func (c *MockT) Log(args ...any) {
}

func (c *MockT) Logf(format string, args ...any) {
}

func (c *MockT) Name() string {
	return c.NameVar
}

func (t *MockT) Parallel() {
}

func (t *MockT) Run(name string, f func(t *T)) bool {
	return true
}

func (t *MockT) Setenv(key, value string) {
}

func (c *MockT) Skip(args ...any) {
}

func (c *MockT) SkipNow() {
}

func (c *MockT) Skipf(format string, args ...any) {
}

func (c *MockT) Skipped() bool {
	return c.SkippedVar
}

func (c *MockT) TempDir() string {
	return ""
}
