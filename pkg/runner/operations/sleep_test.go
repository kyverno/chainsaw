package operations

import (
	"context"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
)

func Test_sleepOperation(t *testing.T) {
	tests := []struct {
		name string
		op   v1alpha1.Sleep
		want Operation
	}{{
		op: v1alpha1.Sleep{
			Duration: metav1.Duration{Duration: time.Duration(1) * time.Second},
		},
		want: sleepAction{
			duration: time.Second,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sleepOperation(tt.op)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sleepAction_Execute(t *testing.T) {
	op := sleepOperation(v1alpha1.Sleep{
		Duration: metav1.Duration{Duration: time.Duration(1) * time.Second},
	})
	o, err := op.Execute(context.Background(), enginecontext.EmptyContext(clock.RealClock{}))
	assert.NoError(t, err)
	assert.Nil(t, o)
}
