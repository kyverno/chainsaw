package internal

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// IsRetriable checks if an error from the API server is a transient error that should be retried.
func IsRetriable(err error) bool {
	if err == nil {
		return false
	}
	return kerrors.IsInternalError(err) ||
		kerrors.IsServiceUnavailable(err) ||
		kerrors.IsTimeout(err) ||
		kerrors.IsServerTimeout(err) ||
		kerrors.IsTooManyRequests(err) ||
		kerrors.IsUnexpectedServerError(err)
}
