package data

import (
	"embed"
	"io/fs"
)

const CrdsFolder = "crds"

//go:embed crds
var crdsFs embed.FS

//go:embed config
var configFs embed.FS

func Crds() fs.FS {
	return crdsFs
}

func Config() fs.FS {
	return configFs
}
