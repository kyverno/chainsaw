package timeout

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// func Combine(config v1alpha1.Timeouts, next ...*v1alpha1.Timeouts) v1alpha1.Timeouts {
// 	for _, next := range next {
// 		if next != nil {
// 			if next.Apply != nil {
// 				config.Apply = next.Apply
// 			}
// 			if next.Assert != nil {
// 				config.Assert = next.Assert
// 			}
// 			if next.Error != nil {
// 				config.Error = next.Error
// 			}
// 			if next.Delete != nil {
// 				config.Delete = next.Delete
// 			}
// 			if next.Cleanup != nil {
// 				config.Cleanup = next.Cleanup
// 			}
// 			if next.Exec != nil {
// 				config.Exec = next.Exec
// 			}
// 		}
// 	}
// 	return config
// }

func Get(operation *metav1.Duration, fallback time.Duration) *time.Duration {
	if operation != nil {
		return &operation.Duration
	}
	return &fallback
}
