package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestIsRetriable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "generic error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "forbidden error",
			err:      kerrors.NewForbidden(schema.GroupResource{Group: "", Resource: "pods"}, "test-pod", errors.New("forbidden")),
			expected: false,
		},
		{
			name:     "unauthorized error",
			err:      kerrors.NewUnauthorized("unauthorized"),
			expected: false,
		},
		{
			name:     "bad request error",
			err:      kerrors.NewBadRequest("bad request"),
			expected: false,
		},
		{
			name:     "internal error",
			err:      kerrors.NewInternalError(errors.New("internal server error")),
			expected: true,
		},
		{
			name:     "service unavailable error",
			err:      kerrors.NewServiceUnavailable("service unavailable"),
			expected: true,
		},
		{
			name:     "timeout error",
			err:      kerrors.NewTimeoutError("timeout", 10),
			expected: true,
		},
		{
			name:     "server timeout error",
			err:      kerrors.NewServerTimeout(schema.GroupResource{Group: "", Resource: "pods"}, "get", 10),
			expected: true,
		},
		{
			name:     "too many requests error",
			err:      kerrors.NewTooManyRequests("too many requests", 10),
			expected: true,
		},
		{
			name:     "unexpected server error",
			err:      kerrors.NewGenericServerResponse(500, "GET", schema.GroupResource{Group: "", Resource: "pods"}, "test-pod", "unexpected body", 0, false),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IsRetriable(tt.err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
