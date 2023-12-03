package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestExpectations(t *testing.T) {
	tests := []struct {
		name       string
		config     bool
		test       *bool
		step       *bool
		wantResult bool
		wantError  bool
	}{
		{
			name:       "Valid Expectations",
			config:     true,
			test:       nil,
			step:       nil,
			wantResult: true,
			wantError:  false,
		},
		{
			name:       "Valid Expectations from Test",
			config:     false,
			test:       ptr.To(true),
			step:       nil,
			wantResult: true,
			wantError:  false,
		},
		{
			name:       "Invalid Expectations from Test",
			config:     true,
			test:       ptr.To(false),
			step:       nil,
			wantResult: false,
			wantError:  false,
		},
		{
			name:       "Valid Expectations from Step",
			config:     false,
			test:       ptr.To(false),
			step:       ptr.To(true),
			wantResult: true,
			wantError:  false,
		},
		{
			name:       "Invalid Expectations from Step",
			config:     true,
			test:       ptr.To(true),
			step:       ptr.To(false),
			wantResult: false,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotError := Expectations(tt.config, tt.test, tt.step)

			// Validate the result
			assert.Equal(t, tt.wantResult, gotResult)

			// Validate the error
			if tt.wantError {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}
