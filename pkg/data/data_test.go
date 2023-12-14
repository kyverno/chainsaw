package data

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrds(t *testing.T) {
	data := Crds()
	{
		file, err := fs.Stat(data, "crds/chainsaw.kyverno.io_configurations.yaml")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.False(t, file.IsDir())
	}
	{
		file, err := fs.Stat(data, "crds")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.True(t, file.IsDir())
	}
}

func TestConfig(t *testing.T) {
	data := Config()
	{
		file, err := fs.Stat(data, "config/default.yaml")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.False(t, file.IsDir())
	}
	{
		file, err := fs.Stat(data, "config")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.True(t, file.IsDir())
	}
}

func TestSchemas(t *testing.T) {
	data := Schemas()
	{
		file, err := fs.Stat(data, "schemas/json/test-chainsaw-v1alpha1.json")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.False(t, file.IsDir())
	}
	{
		file, err := fs.Stat(data, "schemas/json")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.True(t, file.IsDir())
	}
	{
		file, err := fs.Stat(data, "schemas")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.True(t, file.IsDir())
	}
}
