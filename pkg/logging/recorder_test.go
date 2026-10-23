package logging

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWithRecorder(t *testing.T) {
	pod := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "ns", Name: "pod"},
	}
	tests := []struct {
		name      string
		operation Operation
		status    Status
		obj       client.Object
		args      []fmt.Stringer
		want      string
	}{{
		name:      "without object and arguments",
		operation: Try,
		status:    BeginStatus,
		want:      "TRY | BEGIN\n",
	}, {
		name:      "with object",
		operation: Apply,
		status:    DoneStatus,
		obj:       pod,
		want:      "APPLY | DONE | v1/Pod @ ns/pod\n",
	}, {
		name:      "with sections",
		operation: Command,
		status:    LogStatus,
		args:      []fmt.Stringer{Section("STDOUT", "out"), Section("STDERR", "err")},
		want:      "CMD | LOG\n=== STDOUT\nout\n=== STDERR\nerr\n",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var forwarded int
			next := LoggerFunc(func(_ context.Context, operation Operation, status Status, obj client.Object, _ *color.Color, args ...fmt.Stringer) {
				assert.Equal(t, tt.operation, operation)
				assert.Equal(t, tt.status, status)
				assert.Equal(t, tt.obj, obj)
				assert.Equal(t, tt.args, args)
				forwarded++
			})
			var out bytes.Buffer
			ctx := WithRecorder(WithLogger(context.Background(), next), &out)
			Log(ctx, tt.operation, tt.status, tt.obj, color.New(color.FgBlue), tt.args...)
			assert.Equal(t, tt.want, out.String())
			assert.Equal(t, 1, forwarded)
		})
	}
}

func TestWithRecorder_withoutLogger(t *testing.T) {
	var out bytes.Buffer
	ctx := WithRecorder(context.Background(), &out)
	Log(ctx, Sleep, DoneStatus, nil, nil)
	assert.Equal(t, "SLEEP | DONE\n", out.String())
}
