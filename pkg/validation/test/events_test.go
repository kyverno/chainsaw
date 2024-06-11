package test

import (
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateEvents(t *testing.T) {
	tests := []struct {
		name      string
		input     *v1alpha1.Events
		expectErr bool
		errMsg    string
	}{
		{
			name: "Both Name and Selector provided",
			input: &v1alpha1.Events{
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "example-name",
					},
					Selector: "example-selector",
				},
			},
			expectErr: true,
			errMsg:    "a name or label selector must be specified (found both)",
		},
		{
			name: "Only Name provided",
			input: &v1alpha1.Events{
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "example-name",
					},
					Selector: "",
				},
			},
			expectErr: false,
		},
		{
			name: "Only Selector provided",
			input: &v1alpha1.Events{
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "",
					},
					Selector: "example-selector",
				},
			},
			expectErr: false,
		},
		{
			name:      "Neither Name nor Selector provided",
			input:     &v1alpha1.Events{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateEvents(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
