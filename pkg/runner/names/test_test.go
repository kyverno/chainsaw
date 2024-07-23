package names

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTest(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	tests := []struct {
		name    string
		full    bool
		test    discovery.Test
		want    string
		wantErr bool
	}{{
		name: "nil test",
		full: false,
		test: discovery.Test{
			BasePath: cwd,
			Test:     nil,
		},
		wantErr: true,
	}, {
		name: "no full name",
		full: false,
		test: discovery.Test{
			BasePath: cwd,
			Test: &model.Test{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
		},
		wantErr: false,
		want:    "foo",
	}, {
		name: "full name",
		full: true,
		test: discovery.Test{
			BasePath: cwd,
			Test: &model.Test{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
		},
		wantErr: false,
		want:    ".[foo]",
	}, {
		name: "full name",
		full: true,
		test: discovery.Test{
			BasePath: filepath.Join(cwd, "..", "dir", "dir"),
			Test: &model.Test{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
		},
		wantErr: false,
		want:    "../dir/dir[foo]",
	}, {
		name: "full name",
		full: true,
		test: discovery.Test{
			BasePath: filepath.Join(cwd, "dir", "dir"),
			Test: &model.Test{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
		},
		wantErr: false,
		want:    "dir/dir[foo]",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Test(tt.full, tt.test)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestHelpTest(t *testing.T) {
	testCases := []struct {
		name           string
		workingDirFunc workignDirInterface
		absPathFunc    absolutePathInterface
		relPathFunc    relativePathInterface
		expectedErrMsg string
	}{
		{
			name:           "ErrorGettingWorkingDirectory",
			workingDirFunc: func() (string, error) { return "", errors.New("working directory error") },
			expectedErrMsg: "failed to get current working dir (working directory error)",
		},
		{
			name:           "ErrorGettingAbsolutePath",
			absPathFunc:    func(string) (string, error) { return "", errors.New("absolute path error") },
			expectedErrMsg: "failed to compute absolute path for /some/path (absolute path error)",
		},
		{
			name:           "ErrorGettingRelativePath",
			workingDirFunc: func() (string, error) { return "/correct/path", nil },
			absPathFunc:    func(string) (string, error) { return "/correct/path/subdir", nil },
			relPathFunc:    func(string, string) (string, error) { return "", errors.New("relative path error") },
			expectedErrMsg: "failed to compute relative path from /correct/path to /correct/path/subdir (relative path error)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			test := discovery.Test{
				BasePath: "/some/path",
				Test: &model.Test{
					ObjectMeta: metav1.ObjectMeta{
						Name: "foo",
					},
				},
			}

			_, err := helpTest(test, tc.workingDirFunc, tc.absPathFunc, tc.relPathFunc)

			if err == nil {
				t.Fatalf("%s: expected an error but got none", tc.name)
			}
			if err.Error() != tc.expectedErrMsg {
				t.Errorf("%s: expected error message '%v', got '%v'", tc.name, tc.expectedErrMsg, err.Error())
			}
		})
	}
}
