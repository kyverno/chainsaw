package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_testName(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	tests := []struct {
		name    string
		config  v1alpha1.ConfigurationSpec
		test    discovery.Test
		want    string
		wantErr bool
	}{{
		name: "nil test",
		config: v1alpha1.ConfigurationSpec{
			FullName: false,
		},
		test: discovery.Test{
			BasePath: cwd,
			Test:     nil,
		},
		wantErr: true,
	}, {
		name: "no full name",
		config: v1alpha1.ConfigurationSpec{
			FullName: false,
		},
		test: discovery.Test{
			BasePath: cwd,
			Test: &v1alpha1.Test{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
		},
		wantErr: false,
		want:    "foo",
	}, {
		name: "full name",
		config: v1alpha1.ConfigurationSpec{
			FullName: true,
		},
		test: discovery.Test{
			BasePath: cwd,
			Test: &v1alpha1.Test{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
		},
		wantErr: false,
		want:    ".[foo]",
	}, {
		name: "full name",
		config: v1alpha1.ConfigurationSpec{
			FullName: true,
		},
		test: discovery.Test{
			BasePath: filepath.Join(cwd, "..", "dir", "dir"),
			Test: &v1alpha1.Test{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
		},
		wantErr: false,
		want:    "../dir/dir[foo]",
	}, {
		name: "full name",
		config: v1alpha1.ConfigurationSpec{
			FullName: true,
		},
		test: discovery.Test{
			BasePath: filepath.Join(cwd, "dir", "dir"),
			Test: &v1alpha1.Test{
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
			got, err := testName(tt.config, tt.test)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
