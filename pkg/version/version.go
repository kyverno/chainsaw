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

type buildInfoReader = func() (*debug.BuildInfo, bool)

func Version() string {
	return version(debug.ReadBuildInfo)
}

func version(reader buildInfoReader) string {
	if BuildVersion == "" {
		bi, ok := reader()
		if !ok {
			return notFound
		}
		BuildVersion = bi.Main.Version
	}
	return BuildVersion
}

func Time() string {
	return time(debug.ReadBuildInfo)
}

func time(reader buildInfoReader) string {
	bi, ok := reader()
	if !ok {
		return notFound
	}
	return tryFindSetting(vcsTime, bi.Settings...)
}

func Hash() string {
	return hash(debug.ReadBuildInfo)
}

func hash(reader buildInfoReader) string {
	bi, ok := reader()
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
