package processors

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func withTimeouts(timeouts v1alpha1.DefaultTimeouts, override v1alpha1.Timeouts) v1alpha1.DefaultTimeouts {
	if new := override.Apply; new != nil {
		timeouts.Apply = *new
	}
	if new := override.Assert; new != nil {
		timeouts.Assert = *new
	}
	if new := override.Cleanup; new != nil {
		timeouts.Cleanup = *new
	}
	if new := override.Delete; new != nil {
		timeouts.Delete = *new
	}
	if new := override.Error; new != nil {
		timeouts.Error = *new
	}
	if new := override.Exec; new != nil {
		timeouts.Exec = *new
	}
	return timeouts
}
