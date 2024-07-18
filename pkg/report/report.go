package report

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/discovery"
)

type OperationType string

const (
	OperationTypeCreate  OperationType = "create"
	OperationTypeDelete  OperationType = "delete"
	OperationTypeApply   OperationType = "apply"
	OperationTypeAssert  OperationType = "assert"
	OperationTypeError   OperationType = "error"
	OperationTypeScript  OperationType = "script"
	OperationTypeSleep   OperationType = "sleep"
	OperationTypeCommand OperationType = "command"
)

type Report struct {
	name      string
	startTime time.Time
	endTime   time.Time
	tests     []*TestReport
	lock      sync.Mutex
}

func New(name string) *Report {
	return &Report{
		name: name,
	}
}

func (r *Report) SetStartTime(t time.Time) {
	r.startTime = t
}

func (r *Report) SetEndTime(t time.Time) {
	r.endTime = t
}

func (r *Report) ForTest(test *discovery.Test) *TestReport {
	out := &TestReport{test: test}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.tests = append(r.tests, out)
	return out
}

func (r *Report) Save(format v1alpha2.ReportFormatType, path, name string) error {
	if filepath.Ext(name) == "" {
		name += "." + strings.ToLower(string(format))
	}
	filePath := name
	if path != "" {
		filePath = filepath.Join(path, name)
	}
	switch format {
	case v1alpha2.XMLFormat:
		return saveJUnit(r, filePath)
	case v1alpha2.JSONFormat:
	default:
		return fmt.Errorf("unknown report format: %s", format)
	}
	return nil
}

type TestReport struct {
	test      *discovery.Test
	startTime time.Time
	endTime   time.Time
	namespace string
	skipped   bool
	failed    bool
	steps     []*StepReport
	lock      sync.Mutex
}

func (r *TestReport) SetStartTime(t time.Time) {
	r.startTime = t
}

func (r *TestReport) SetEndTime(t time.Time) {
	r.endTime = t
}

func (r *TestReport) Skip() {
	r.skipped = true
}

func (r *TestReport) Fail() {
	r.failed = true
}

func (r *TestReport) SetNamespace(namespace string) {
	r.namespace = namespace
}

func (r *TestReport) ForStep(step *v1alpha1.TestStep) *StepReport {
	out := &StepReport{step: step}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.steps = append(r.steps, out)
	return out
}

type StepReport struct {
	step      *v1alpha1.TestStep
	startTime time.Time
	endTime   time.Time
	reports   []*OperationReport
	lock      sync.Mutex
}

func (r *StepReport) SetStartTime(t time.Time) {
	r.startTime = t
}

func (r *StepReport) SetEndTime(t time.Time) {
	r.endTime = t
}

func (r *StepReport) ForOperation(name string, operationType OperationType) *OperationReport {
	step := &OperationReport{name: name, operationType: operationType}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.reports = append(r.reports, step)
	return step
}

type OperationReport struct {
	name          string
	operationType OperationType
	startTime     time.Time
	endTime       time.Time
	err           error
}

func (r *OperationReport) SetStartTime(t time.Time) {
	r.startTime = t
}

func (r *OperationReport) SetEndTime(t time.Time) {
	r.endTime = t
}
