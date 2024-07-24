package logging

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func TestSection(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected string
	}{
		{"empty", nil, "=== EMPTY"},
		{"single arg", []any{"example"}, "=== SINGLE ARG\nexample"},
		{"multiple args", []any{"one", 2, 3.0}, "=== MULTIPLE ARGS\none2 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sec := Section(tt.name, tt.args...)
			assert.Equal(t, tt.expected, sec.String(), "Output should match expected string")
		})
	}
}

func TestErrSection(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			"single error",
			errors.New("error 1"),
			"=== ERROR\nerror 1",
		},
		{
			"multi error",
			multierr.Combine(errors.New("error 1"), errors.New("error 2")),
			"=== ERROR\nerror 1\nerror 2",
		},
		{
			"aggregate error",
			utilerrors.NewAggregate([]error{errors.New("error 1"), errors.New("error 2")}),
			"=== ERROR\nerror 1\nerror 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sec := ErrSection(tt.err)
			assert.Equal(t, tt.expected, sec.String(), "Output should match expected string")
		})
	}
}
