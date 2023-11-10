package docs

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_options_validate(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		website    bool
		autogenTag bool
		root       *cobra.Command
		wantErr    bool
	}{{
		name:       "nil root",
		path:       "dummy",
		website:    false,
		autogenTag: false,
		root:       nil,
		wantErr:    true,
	}, {
		name:       "dempty path",
		path:       "",
		website:    false,
		autogenTag: false,
		root:       &cobra.Command{},
		wantErr:    true,
	}, {
		name:       "ok",
		path:       "dummy",
		website:    false,
		autogenTag: false,
		root:       &cobra.Command{},
		wantErr:    false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := options{
				path:       tt.path,
				website:    tt.website,
				autogenTag: tt.autogenTag,
			}
			err := o.validate(tt.root)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
