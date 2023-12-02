package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DefaultApplyTimeout   = 5 * time.Second
	DefaultAssertTimeout  = 30 * time.Second
	DefaultCleanupTimeout = 30 * time.Second
	DefaultDeleteTimeout  = 15 * time.Second
	DefaultErrorTimeout   = 30 * time.Second
	DefaultExecTimeout    = 5 * time.Second
)

// Timeouts contains timeouts per operation.
type Timeouts struct {
	// Apply defines the timeout for the apply operation
	Apply *metav1.Duration `json:"apply,omitempty"`

	// Assert defines the timeout for the assert operation
	Assert *metav1.Duration `json:"assert,omitempty"`

	// Cleanup defines the timeout for the cleanup operation
	Cleanup *metav1.Duration `json:"cleanup,omitempty"`

	// Delete defines the timeout for the delete operation
	Delete *metav1.Duration `json:"delete,omitempty"`

	// Error defines the timeout for the error operation
	Error *metav1.Duration `json:"error,omitempty"`

	// Exec defines the timeout for exec operations
	Exec *metav1.Duration `json:"exec,omitempty"`
}

func durationOrDefault(to *metav1.Duration, def time.Duration) time.Duration {
	if to != nil {
		return to.Duration
	}
	return def
}

func (t Timeouts) ApplyDuration() time.Duration {
	return durationOrDefault(t.Apply, DefaultApplyTimeout)
}

func (t Timeouts) AssertDuration() time.Duration {
	return durationOrDefault(t.Assert, DefaultAssertTimeout)
}

func (t Timeouts) CleanupDuration() time.Duration {
	return durationOrDefault(t.Cleanup, DefaultCleanupTimeout)
}

func (t Timeouts) DeleteDuration() time.Duration {
	return durationOrDefault(t.Delete, DefaultDeleteTimeout)
}

func (t Timeouts) ErrorDuration() time.Duration {
	return durationOrDefault(t.Error, DefaultErrorTimeout)
}

func (t Timeouts) ExecDuration() time.Duration {
	return durationOrDefault(t.Exec, DefaultExecTimeout)
}

func (t Timeouts) Combine(override *Timeouts) Timeouts {
	if override == nil {
		return t
	}
	if override.Apply != nil {
		t.Apply = override.Apply
	}
	if override.Assert != nil {
		t.Assert = override.Assert
	}
	if override.Error != nil {
		t.Error = override.Error
	}
	if override.Delete != nil {
		t.Delete = override.Delete
	}
	if override.Cleanup != nil {
		t.Cleanup = override.Cleanup
	}
	if override.Exec != nil {
		t.Exec = override.Exec
	}
	return t
}
