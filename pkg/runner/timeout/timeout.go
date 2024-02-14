package timeout

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Get(operation *metav1.Duration, fallback time.Duration) *time.Duration {
	if operation != nil {
		return &operation.Duration
	}
	return &fallback
}
