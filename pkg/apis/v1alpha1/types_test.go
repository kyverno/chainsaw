package v1alpha1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBinding_CheckName(t *testing.T) {
	tests := []struct {
		name         string
		bindingName  Expression
		bindingValue Any
		wantErr      bool
	}{{
		name:    "empty",
		wantErr: true,
	}, {
		name:        "simple",
		bindingName: "simple",
		wantErr:     false,
	}, {
		name:        "with dollar",
		bindingName: "$simple",
		wantErr:     true,
	}, {
		name:        "with space",
		bindingName: "simple one",
		wantErr:     true,
	}, {
		name:        "with dot",
		bindingName: "simple.one",
		wantErr:     true,
	}, {
		name:        "good expression",
		bindingName: "('test')",
		wantErr:     false,
	}, {
		name:        "bad expression",
		bindingName: "('test'",
		wantErr:     true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Binding{
				Name:  tt.bindingName,
				Value: tt.bindingValue,
			}
			err := b.CheckName()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTimeouts_Combine(t *testing.T) {
	base := DefaultTimeouts{
		Apply:   metav1.Duration{Duration: 1 * time.Minute},
		Assert:  metav1.Duration{Duration: 1 * time.Minute},
		Cleanup: metav1.Duration{Duration: 1 * time.Minute},
		Delete:  metav1.Duration{Duration: 1 * time.Minute},
		Error:   metav1.Duration{Duration: 1 * time.Minute},
		Exec:    metav1.Duration{Duration: 1 * time.Minute},
	}
	override := Timeouts{
		Apply:   &metav1.Duration{Duration: 2 * time.Minute},
		Assert:  &metav1.Duration{Duration: 2 * time.Minute},
		Cleanup: &metav1.Duration{Duration: 2 * time.Minute},
		Delete:  &metav1.Duration{Duration: 2 * time.Minute},
		Error:   &metav1.Duration{Duration: 2 * time.Minute},
		Exec:    &metav1.Duration{Duration: 2 * time.Minute},
	}
	tests := []struct {
		name     string
		base     DefaultTimeouts
		override *Timeouts
		want     DefaultTimeouts
	}{{
		name:     "nil",
		base:     base,
		override: nil,
		want:     base,
	}, {
		name:     "override",
		base:     base,
		override: &override,
		want: DefaultTimeouts{
			Apply:   metav1.Duration{Duration: 2 * time.Minute},
			Assert:  metav1.Duration{Duration: 2 * time.Minute},
			Cleanup: metav1.Duration{Duration: 2 * time.Minute},
			Delete:  metav1.Duration{Duration: 2 * time.Minute},
			Error:   metav1.Duration{Duration: 2 * time.Minute},
			Exec:    metav1.Duration{Duration: 2 * time.Minute},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.Combine(tt.override)
			assert.Equal(t, tt.want, got)
		})
	}
}
