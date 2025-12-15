package names

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestStep(t *testing.T) {
	tests := []struct {
		name string
		step v1alpha1.TestStep
		i    int
		want string
	}{{
		name: "no name",
		step: v1alpha1.TestStep{
			Name: "",
		},
		i:    10,
		want: "step #11",
	}, {
		name: "with name",
		step: v1alpha1.TestStep{
			Name: "foo",
		},
		i:    10,
		want: "foo",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Step(tt.step, tt.i)
			assert.Equal(t, tt.want, got)
		})
	}
}
