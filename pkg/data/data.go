package data

import (
	"embed"
	"io/fs"
	"sync"
)

const (
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

func _configFile() ([]byte, error) {
	configFs, err := config()
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(configFs, "default.yaml")
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
	ConfigFile = sync.OnceValues(_configFile)
	Crds       = sync.OnceValues(_crds)
	Schemas    = sync.OnceValues(_schemas)
)
