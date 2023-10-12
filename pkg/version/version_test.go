package version

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	require.NotEqual(t, "---", Version())
}

func TestBuildVersion(t *testing.T) {
	BuildVersion = "test"
	require.Equal(t, "test", Version())
}

func TestTime(t *testing.T) {
	require.Equal(t, "---", Time())
}

func TestHash(t *testing.T) {
	require.Equal(t, "---", Hash())
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
