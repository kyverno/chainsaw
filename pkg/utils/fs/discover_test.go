package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscoverFolders(t *testing.T) {
	root := t.TempDir()
	dirs := []string{"dir1", "dir2", "dir2/subdir1", "dir3", "dir3/subdir1", "dir3/subdir1/subsubdir1"}
	files := []string{"file1", "dir2/file2", "dir3/subdir1/file3"}
	for _, dir := range dirs {
		assert.NoError(t, os.MkdirAll(filepath.Join(root, dir), os.ModePerm))
	}
	for _, file := range files {
		assert.NoError(t, os.WriteFile(filepath.Join(root, file), []byte("test"), 0o600))
	}
	discovered, err := DiscoverFolders(root)
	assert.NoError(t, err)
	expectedDirs := []string{root}
	for _, dir := range dirs {
		expectedDirs = append(expectedDirs, filepath.Join(root, dir))
	}
	assert.ElementsMatch(t, expectedDirs, discovered)
}

func TestDiscoverFoldersWithError(t *testing.T) {
	root := t.TempDir()
	unreadableDir := filepath.Join(root, "unreadable")
	assert.NoError(t, os.MkdirAll(unreadableDir, os.ModePerm))
	assert.NoError(t, os.Chmod(unreadableDir, 0o000))
	_, err := DiscoverFolders(unreadableDir)
	assert.Error(t, err)
}
