package data

import (
	"embed"
	"io/fs"
	"sync"
)

const (
	crdsFolder    = "crds"
	configFolder  = "config"
	schemasFolder = "schemas/json"
)

//go:embed crds
var crdsFs embed.FS

//go:embed config
var configFs embed.FS

//go:embed schemas
var schemasFs embed.FS

func sub(f embed.FS, dir string) (fs.FS, error) {
	return fs.Sub(f, dir)
}

func crds() (fs.FS, error) {
	return sub(crdsFs, crdsFolder)
}

func config() (fs.FS, error) {
	return sub(configFs, configFolder)
}

func schemas() (fs.FS, error) {
	return sub(schemasFs, schemasFolder)
}

func configFile() ([]byte, error) {
	configFs, err := Config()
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(configFs, "default.yaml")
}

var (
	Crds       = sync.OnceValues(crds)
	Config     = sync.OnceValues(config)
	Schemas    = sync.OnceValues(schemas)
	ConfigFile = sync.OnceValues(configFile)
)
