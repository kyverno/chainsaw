package data

import (
	"embed"
	"io/fs"
	"sync"
)

const (
	configFile    = "default.yaml"
	configFolder  = "config"
	crdsFolder    = "crds"
	schemasFolder = "schemas/json"
)

//go:embed crds
var crdsFs embed.FS

//go:embed config
var configFs embed.FS

//go:embed schemas
var schemasFs embed.FS

func _config() (fs.FS, error) {
	return _sub(configFs, configFolder)
}

func _configFile(_fs func() (fs.FS, error)) ([]byte, error) {
	if _fs == nil {
		_fs = config
	}
	configFs, err := _fs()
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(configFs, configFile)
}

func _crds() (fs.FS, error) {
	return _sub(crdsFs, crdsFolder)
}

func _schemas() (fs.FS, error) {
	return _sub(schemasFs, schemasFolder)
}

func _sub(f embed.FS, dir string) (fs.FS, error) {
	return fs.Sub(f, dir)
}

var (
	config     = sync.OnceValues(_config)
	ConfigFile = sync.OnceValues(func() ([]byte, error) { return _configFile(nil) })
	Crds       = sync.OnceValues(_crds)
	Schemas    = sync.OnceValues(_schemas)
)
