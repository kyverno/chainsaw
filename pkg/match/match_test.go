package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		hasError bool
		errMsg   string
	}{
		{
			name:     "Exact Match",
			expected: "hello",
			actual:   "hello",
			hasError: false,
		},
		{
			name:     "Type Mismatch",
			expected: 1,
			actual:   "1",
			hasError: true,
			errMsg:   "type mismatch: int != string",
		},
		{
			name:     "Slice Length Mismatch",
			expected: []int{1, 2},
			actual:   []int{1, 2, 3},
			hasError: true,
			errMsg:   "slice length mismatch: 2 != 3",
		},
		{
			name:     "Missing Map Key",
			expected: map[string]int{"a": 1},
			actual:   map[string]int{"b": 2},
			hasError: true,
			errMsg:   ".a: key is missing from map",
		},
		{
			name:     "Nested Structure Match",
			expected: map[string]interface{}{"a": []int{1, 2}},
			actual:   map[string]interface{}{"a": []int{1, 2}},
			hasError: false,
		},
		{
			name:     "Nested Structure Mismatch",
			expected: map[string]interface{}{"a": []int{1, 2}},
			actual:   map[string]interface{}{"a": []int{1, 2, 3}},
			hasError: true,
			errMsg:   ".a: slice length mismatch: 2 != 3",
		},
		{
			name:     "Map Exact Match with Extra Keys",
			expected: map[string]int{"a": 1},
			actual:   map[string]int{"a": 1, "b": 2},
			hasError: false,
		},
		{
			name:     "Nested Maps with Missing Key",
			expected: map[string]interface{}{"a": map[string]int{"b": 2}},
			actual:   map[string]interface{}{"a": map[string]int{"c": 3}},
			hasError: true,
			errMsg:   ".a.b: key is missing from map",
		},
		{
			name:     "Slices with Different Elements",
			expected: []int{1, 2},
			actual:   []int{1, 3},
			hasError: true,
			errMsg:   "value mismatch, expected: 2 != actual: 3",
		},
		{
			name:     "Nested Slices Mismatch",
			expected: [][]int{{1, 2}, {3, 4}},
			actual:   [][]int{{1, 2}, {4, 5}},
			hasError: true,
			errMsg:   "value mismatch, expected: 3 != actual: 4",
		},
		{
			name:     "Bool Mismatch",
			expected: true,
			actual:   false,
			hasError: true,
			errMsg:   "value mismatch, expected: true != actual: false",
		},
		{
			name:     "String Mismatch",
			expected: "hello",
			actual:   "world",
			hasError: true,
			errMsg:   "value mismatch, expected: hello != actual: world",
		},
		{
			name:     "Float Mismatch",
			expected: 1.1,
			actual:   1.2,
			hasError: true,
			errMsg:   "value mismatch, expected: 1.1 != actual: 1.2",
		},
		{
			name:     "Int Mismatch",
			expected: 1,
			actual:   2,
			hasError: true,
			errMsg:   "value mismatch, expected: 1 != actual: 2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Match(test.expected, test.actual)
			if test.hasError {
				assert.Error(t, err)
				assert.Equal(t, test.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
