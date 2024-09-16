package data

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrds(t *testing.T) {
	data, err := Crds()
	assert.NoError(t, err)
	files := []string{
		"chainsaw.kyverno.io_configurations.yaml",
		"chainsaw.kyverno.io_steptemplates.yaml",
		"chainsaw.kyverno.io_tests.yaml",
	}
	for _, file := range files {
		file, err := fs.Stat(data, file)
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.False(t, file.IsDir())
	}
}

func Test_config(t *testing.T) {
	data, err := config()
	assert.NoError(t, err)
	files := []string{
		"default.yaml",
	}
	for _, file := range files {
		file, err := fs.Stat(data, file)
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.False(t, file.IsDir())
	}
}

func TestSchemas(t *testing.T) {
	data, err := Schemas()
	assert.NoError(t, err)
	files := []string{
		"configuration-chainsaw-v1alpha1.json",
		"configuration-chainsaw-v1alpha2.json",
		"steptemplate-chainsaw-v1alpha1.json",
		"test-chainsaw-v1alpha1.json",
		"test-chainsaw-v1alpha2.json",
	}
	for _, file := range files {
		file, err := fs.Stat(data, file)
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.False(t, file.IsDir())
	}
}

func TestConfigFile(t *testing.T) {
	data, err := ConfigFile()
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func Test_configFile(t *testing.T) {
	data, err := _configFile(func() (fs.FS, error) {
		return nil, errors.New("dummy")
	})
	assert.Error(t, err)
	assert.Nil(t, data)
}
