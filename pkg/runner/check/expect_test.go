package check

import (
	"context"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestExpectations(t *testing.T) {
	// Mock data for testing
	obj := unstructured.Unstructured{}
	bindings := binding.Bindings{}

	tests := []struct {
		name        string
		expectation v1alpha1.Expectation
		checkResult bool
		wantMatched bool
		wantErrors  int
	}{
		{
			name: "Valid Expectation",
			expectation: v1alpha1.Expectation{
				Check: "your_valid_check_value",
			},
			checkResult: true, // Set to true if the check should pass
			wantMatched: true,
			wantErrors:  0,
		},
		{
			name: "Invalid Expectation - Check Fails",
			expectation: v1alpha1.Expectation{
				Check: "your_invalid_check_value",
			},
			checkResult: false, // Set to false if the check should fail
			wantMatched: true,
			wantErrors:  1,
		},
		{
			name: "Expectation with Match",
			expectation: v1alpha1.Expectation{
				Match: &v1alpha1.Match{
					Value: "your_match_value",
				},
				Check: "your_check_value",
			},
			checkResult: true, // Set to true if the check should pass
			wantMatched: true,
			wantErrors:  0,
		},
		{
			name: "Expectation with Match - Match Fails",
			expectation: v1alpha1.Expectation{
				Match: &v1alpha1.Match{
					Value: "your_match_value",
				},
				Check: "your_check_value",
			},
			checkResult: false, // Set to false if the match should fail
			wantMatched: false,
			wantErrors:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the Check function result based on the test case
			mockCheck := func(ctx context.Context, obj interface{}, bindings binding.Bindings, check *v1alpha1.Check) (field.ErrorList, error) {
				if tt.checkResult {
					return nil, nil
				}
				return field.ErrorList{field.NewError(field.Invalid("field"), nil, "check failed")}, nil
			}

			// Replace the original Check function with the mockCheck function
			originalCheck := Check
			Check = mockCheck
			defer func() { Check = originalCheck }()

			// Call the Expectations function with the test case inputs
			matched, errs := Expectations(context.Background(), obj, bindings, tt.expectation)

			// Validate the results based on the expected values
			assert.Equal(t, tt.wantMatched, matched)
			assert.Equal(t, tt.wantErrors, len(errs))
		})
	}
}
