package version

import (
	"runtime/debug"
)

const (
	notFound    = "---"
	vcsTime     = "vcs.time"
	vcsRevision = "vcs.revision"
)

// BuildVersion is provided at compile-time
var BuildVersion string

func Version() string {
	if BuildVersion == "" {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			return notFound
		}
		BuildVersion = bi.Main.Version
	}
	return BuildVersion
}

func Time() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return notFound
	}
	return tryFindSetting(vcsTime, bi.Settings...)
}

func Hash() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return notFound
	}
	return tryFindSetting(vcsRevision, bi.Settings...)
}

func tryFindSetting(key string, settings ...debug.BuildSetting) string {
	for _, setting := range settings {
		if setting.Key == key {
			return setting.Value
		}
	}
	return notFound
}
