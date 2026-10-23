package runner

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
)

type loggingAction struct {
	err error
}

func (a loggingAction) Execute(ctx context.Context, _ enginecontext.TestContext) (outputs.Outputs, error) {
	logging.Log(ctx, logging.Command, logging.LogStatus, nil, nil, logging.Section("STDOUT", "hello"))
	return nil, a.err
}

func Test_operationName(t *testing.T) {
	tests := []struct {
		name        string
		opName      string
		description string
		action      int
		actions     int
		want        string
	}{{
		name:    "without name and description",
		action:  0,
		actions: 1,
		want:    "",
	}, {
		name:    "with name",
		opName:  "check-pod",
		action:  0,
		actions: 1,
		want:    "check-pod",
	}, {
		name:        "with description",
		description: "check the pod",
		action:      0,
		actions:     1,
		want:        "check the pod",
	}, {
		name:        "name before description",
		opName:      "check-pod",
		description: "check the pod",
		action:      0,
		actions:     1,
		want:        "check-pod",
	}, {
		name:        "second of three actions",
		description: "apply the resources",
		action:      1,
		actions:     3,
		want:        "apply the resources (2/3)",
	}, {
		name:        "name with three actions",
		opName:      "apply-resources",
		description: "apply the resources",
		action:      2,
		actions:     3,
		want:        "apply-resources (3/3)",
	}, {
		name:    "several actions without name and description",
		action:  1,
		actions: 3,
		want:    "",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, operationName(tt.opName, tt.description, tt.action, tt.actions))
		})
	}
}

func Test_runner_runAction_report(t *testing.T) {
	tests := []struct {
		name     string
		opName   string
		err      error
		wantName string
	}{{
		name:     "with name",
		opName:   "check the pod",
		wantName: "check the pod",
	}, {
		name:     "without name",
		wantName: "operation 1",
	}, {
		name:     "with error",
		opName:   "check the pod",
		err:      errors.New("dummy"),
		wantName: "check the pod",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &runner{clock: clock.RealClock{}}
			var forwarded int
			ctx := logging.WithLogger(context.Background(), logging.LoggerFunc(func(context.Context, logging.Operation, logging.Status, client.Object, *color.Color, ...fmt.Stringer) {
				forwarded++
			}))
			stepReport := &model.StepReport{}
			_, err := r.runAction(ctx, loggingAction{err: tt.err}, model.OperationTypeCommand, tt.opName, 0, 0, enginecontext.EmptyContext(clock.RealClock{}), stepReport)
			assert.Equal(t, tt.err, err)
			assert.Len(t, stepReport.Operations, 1)
			report := stepReport.Operations[0]
			assert.Equal(t, tt.wantName, report.Name)
			assert.Equal(t, model.OperationTypeCommand, report.Type)
			assert.Equal(t, tt.err, report.Err)
			assert.Equal(t, "CMD | LOG\n=== STDOUT\nhello\n", report.Output)
			assert.Equal(t, 1, forwarded)
		})
	}
}

func Test_runner_runOperation_name(t *testing.T) {
	tests := []struct {
		name      string
		operation v1alpha1.OperationBase
		want      string
	}{{
		name:      "name",
		operation: v1alpha1.OperationBase{Name: "wait", Description: "wait a moment"},
		want:      "wait",
	}, {
		name:      "description",
		operation: v1alpha1.OperationBase{Description: "wait a moment"},
		want:      "wait a moment",
	}, {
		name: "default",
		want: "operation 1",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &runner{clock: clock.RealClock{}}
			operation := v1alpha1.Operation{
				OperationBase: tt.operation,
				Sleep:         &v1alpha1.Sleep{Duration: metav1.Duration{}},
			}
			stepReport := &model.StepReport{}
			_, _, err := r.runOperation(context.Background(), enginecontext.EmptyContext(clock.RealClock{}), operation, 0, cleaner.New(0, false, nil, metav1.DeletePropagationBackground), stepReport)
			assert.NoError(t, err)
			assert.Len(t, stepReport.Operations, 1)
			assert.Equal(t, tt.want, stepReport.Operations[0].Name)
			assert.Equal(t, model.OperationTypeSleep, stepReport.Operations[0].Type)
			assert.Contains(t, stepReport.Operations[0].Output, "SLEEP | DONE")
		})
	}
}
