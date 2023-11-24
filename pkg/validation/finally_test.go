package validation

import (
	"fmt"
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateFinally(t *testing.T) {
	examplePodLogs := &v1alpha1.PodLogs{
		Selector: "app=example",
	}
	exampleEvents := &v1alpha1.Events{}
	exampleCommand := &v1alpha1.Command{
		Entrypoint: "echo",
		Args:       []string{"Hello, World!"},
	}
	exampleScript := &v1alpha1.Script{
		Content: "echo Hello, World!",
	}

	tests := []struct {
		name      string
		input     v1alpha1.Finally
		expectErr bool
		errMsg    string
	}{
		{
			name:      "No Finally statements provided",
			input:     v1alpha1.Finally{},
			expectErr: true,
			errMsg:    "no statement found in operation",
		},
		{
			name: "Multiple Finally statements provided",
			input: v1alpha1.Finally{
				PodLogs: examplePodLogs,
				Events:  exampleEvents,
				Command: exampleCommand,
			},
			expectErr: true,
			errMsg:    fmt.Sprintf("only one statement is allowed per operation (found %d)", 3),
		},
		{
			name: "Only PodLogs statement provided",
			input: v1alpha1.Finally{
				PodLogs: examplePodLogs,
			},
			expectErr: false,
		},
		{
			name: "Only Events statement provided",
			input: v1alpha1.Finally{
				Events: exampleEvents,
			},
			expectErr: false,
		},
		{
			name: "Only Command statement provided",
			input: v1alpha1.Finally{
				Command: exampleCommand,
			},
			expectErr: false,
		},
		{
			name: "Only Script statement provided",
			input: v1alpha1.Finally{
				Script: exampleScript,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateFinally(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
