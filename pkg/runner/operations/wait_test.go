package operations

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/stretchr/testify/assert"
)

func Test_waitOperation(t *testing.T) {
	tests := []struct {
		name       string
		namespacer namespacer.Namespacer
		op         v1alpha1.Wait
		want       Operation
	}{{
		namespacer: namespacer.New("bar"),
		op:         v1alpha1.Wait{},
		want: waitAction{
			op: v1alpha1.Wait{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := waitOperation(tt.op)
			assert.Equal(t, tt.want, got)
		})
	}
}
