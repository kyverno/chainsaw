package version

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	require.Equal(t, "(devel)", Version())
}

func Test_version(t *testing.T) {
	testCases := []struct {
		name            string
		reader          buildInfoReader
		expectedVersion string
	}{{
		name: "No Build Info",
		reader: func() (*debug.BuildInfo, bool) {
			return nil, false
		},
		expectedVersion: notFound,
	}, {
		name: "Valid Build Info",
		reader: func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{
				Main: debug.Module{Version: "v1.0.0"},
			}, true
		},
		expectedVersion: "v1.0.0",
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			BuildVersion = ""
			actualVersion := version(tc.reader)
			require.Equal(t, tc.expectedVersion, actualVersion, "Expected version to be '%s'", tc.expectedVersion)
		})
	}
}

func TestBuildVersion(t *testing.T) {
	BuildVersion = "test"
	require.Equal(t, "test", Version())
}

func TestTime(t *testing.T) {
	require.Equal(t, "---", Time())
}

func Test_time(t *testing.T) {
	vcsTime := "vcs.time"
	testCases := []struct {
		name         string
		reader       buildInfoReader
		expectedTime string
	}{{
		name: "No Build Info",
		reader: func() (*debug.BuildInfo, bool) {
			return nil, false
		},
		expectedTime: notFound,
	}, {
		name: "VCS Time Found",
		reader: func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: vcsTime, Value: "2021-04-01T12:34:56Z"},
				},
			}, true
		},
		expectedTime: "2021-04-01T12:34:56Z",
	}, {
		name: "VCS Time Not Found",
		reader: func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "some.other.setting", Value: "some-value"},
				},
			}, true
		},
		expectedTime: notFound,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualTime := time(tc.reader)
			require.Equal(t, tc.expectedTime, actualTime, "Expected time to be '%s'", tc.expectedTime)
		})
	}
}

func TestHash(t *testing.T) {
	require.Equal(t, "---", Hash())
}

func Test_hash(t *testing.T) {
	vcsRevision := "vcs.revision"
	testCases := []struct {
		name         string
		reader       buildInfoReader
		expectedHash string
	}{{
		name: "No Build Info",
		reader: func() (*debug.BuildInfo, bool) {
			return nil, false
		},
		expectedHash: notFound,
	}, {
		name: "VCS Revision Found",
		reader: func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: vcsRevision, Value: "abcdef123456"},
				},
			}, true
		},
		expectedHash: "abcdef123456",
	}, {
		name: "VCS Revision Not Found",
		reader: func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "some.other.setting", Value: "some-value"},
				},
			}, true
		},
		expectedHash: notFound,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualHash := hash(tc.reader)
			require.Equal(t, tc.expectedHash, actualHash, "Expected hash to be '%s'", tc.expectedHash)
		})
	}
}

func Test_tryFindSetting(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		settings []debug.BuildSetting
		want     string
	}{{
		name: "nil",
		key:  vcsTime,
		want: notFound,
	}, {
		name: "not found",
		key:  vcsTime,
		settings: []debug.BuildSetting{{
			Key:   "foo",
			Value: "bar",
		}},
		want: notFound,
	}, {
		name: "found",
		key:  vcsTime,
		settings: []debug.BuildSetting{{
			Key:   vcsTime,
			Value: "bar",
		}},
		want: "bar",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tryFindSetting(tt.key, tt.settings...); got != tt.want {
				t.Errorf("tryeFindSetting() = %v, want %v", got, tt.want)
			}
		})
	}
}
