package v1alpha1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_durationOrDefault(t *testing.T) {
	tests := []struct {
		name string
		to   *metav1.Duration
		def  time.Duration
		want time.Duration
	}{{
		name: "nil",
		to:   nil,
		def:  time.Second * 3,
		want: time.Second * 3,
	}, {
		name: "not nil",
		to:   &metav1.Duration{Duration: time.Second * 2},
		def:  time.Second * 3,
		want: time.Second * 2,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := durationOrDefault(tt.to, tt.def)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeouts_Defaults(t *testing.T) {
	var timeouts Timeouts
	assert.Equal(t, DefaultApplyTimeout, timeouts.ApplyDuration())
	assert.Equal(t, DefaultAssertTimeout, timeouts.AssertDuration())
	assert.Equal(t, DefaultCleanupTimeout, timeouts.CleanupDuration())
	assert.Equal(t, DefaultDeleteTimeout, timeouts.DeleteDuration())
	assert.Equal(t, DefaultErrorTimeout, timeouts.ErrorDuration())
	assert.Equal(t, DefaultExecTimeout, timeouts.ExecDuration())
}

func TestTimeouts_NoyDefaults(t *testing.T) {
	to := &metav1.Duration{Duration: time.Hour * 2}
	timeouts := Timeouts{
		Apply:   to,
		Assert:  to,
		Cleanup: to,
		Delete:  to,
		Error:   to,
		Exec:    to,
	}
	assert.Equal(t, time.Hour*2, timeouts.ApplyDuration())
	assert.Equal(t, time.Hour*2, timeouts.AssertDuration())
	assert.Equal(t, time.Hour*2, timeouts.CleanupDuration())
	assert.Equal(t, time.Hour*2, timeouts.DeleteDuration())
	assert.Equal(t, time.Hour*2, timeouts.ErrorDuration())
	assert.Equal(t, time.Hour*2, timeouts.ExecDuration())
}

func TestTimeouts_Combine(t *testing.T) {
	base := Timeouts{
		Apply:   &metav1.Duration{Duration: 1 * time.Minute},
		Assert:  &metav1.Duration{Duration: 1 * time.Minute},
		Cleanup: &metav1.Duration{Duration: 1 * time.Minute},
		Delete:  &metav1.Duration{Duration: 1 * time.Minute},
		Error:   &metav1.Duration{Duration: 1 * time.Minute},
		Exec:    &metav1.Duration{Duration: 1 * time.Minute},
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
		base     Timeouts
		override *Timeouts
		want     Timeouts
	}{{
		name:     "nil",
		base:     base,
		override: nil,
		want:     base,
	}, {
		name:     "override",
		base:     base,
		override: &override,
		want:     override,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.Combine(tt.override)
			assert.Equal(t, tt.want, got)
		})
	}
}
