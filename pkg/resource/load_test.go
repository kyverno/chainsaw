package resource

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	baseDir := filepath.Join("..", "..", "testdata", "resource")
	tests := []struct {
		fileName    string
		expectError bool
		expectedLen int
	}{
		{filepath.Join(baseDir, "valid.yaml"), false, 2},
		{filepath.Join(baseDir, "empty.yaml"), true, 0},
		{filepath.Join(baseDir, "invalid.yaml"), true, 0},
		{filepath.Join(baseDir, "nonexistent.yaml"), true, 0},
	}

	for _, tt := range tests {
		resources, err := Load(tt.fileName)
		if !tt.expectError {
			assert.NoError(t, err)
			assert.Len(t, resources, tt.expectedLen)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParse(t *testing.T) {
	baseDir := filepath.Join("..", "..", "testdata", "resource")
	tests := []struct {
		fileName    string
		expectError bool
		expectedLen int
	}{
		{filepath.Join(baseDir, "valid.yaml"), false, 2},
		{filepath.Join(baseDir, "empty.yaml"), false, 0},
		{filepath.Join(baseDir, "invalid.yaml"), true, 0},
	}

	for _, tt := range tests {
		content, readErr := os.ReadFile(tt.fileName)
		assert.NoError(t, readErr)

		resources, err := Parse(content)
		if !tt.expectError {
			assert.NoError(t, err)
			assert.Len(t, resources, tt.expectedLen)
		} else {
			assert.Error(t, err)
		}
	}
}
