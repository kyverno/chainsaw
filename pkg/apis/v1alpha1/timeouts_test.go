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
