package test

import (
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateObjectReference(t *testing.T) {
	tests := []struct {
		name      string
		input     v1alpha1.ObjectReference
		expectErr bool
		errMsgs   []string
	}{
		{
			name: "Both Kind and APIVersion are empty",
			input: v1alpha1.ObjectReference{
				Kind:       "",
				APIVersion: "",
			},
			expectErr: true,
			errMsgs:   []string{"kind must be specified", "apiVersion must be specified"},
		},
		{
			name: "Kind is provided, APIVersion is empty",
			input: v1alpha1.ObjectReference{
				Kind:       "Pod",
				APIVersion: "",
			},
			expectErr: true,
			errMsgs:   []string{"apiVersion must be specified"},
		},
		{
			name: "APIVersion is provided, Kind is empty",
			input: v1alpha1.ObjectReference{
				Kind:       "",
				APIVersion: "v1",
			},
			expectErr: true,
			errMsgs:   []string{"kind must be specified"},
		},
		{
			name: "Both Kind and APIVersion are provided",
			input: v1alpha1.ObjectReference{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateObjectReference(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				for _, msg := range tt.errMsgs {
					assert.Contains(t, errs.ToAggregate().Error(), msg)
				}
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
