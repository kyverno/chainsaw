package yaml

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_remarshal(t *testing.T) {
	in, err := os.ReadFile("with-anchors.yaml")
	assert.NoError(t, err)
	out, err := os.ReadFile("without-anchors.yaml")
	assert.NoError(t, err)
	tests := []struct {
		name      string
		document  string
		unmarshal func(in []byte, out any) (err error)
		want      string
		wantErr   bool
	}{{
		name:     "ok",
		document: string(in),
		want:     string(out),
		wantErr:  false,
	}, {
		name: "error",
		unmarshal: func(in []byte, out any) (err error) {
			return errors.New("dummy")
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := remarshal([]byte(tt.document), tt.unmarshal)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, strings.ReplaceAll(tt.want, "\r\n", "\n"), strings.ReplaceAll(string(got), "\r\n", "\n"))
		})
	}
}

func TestRemarshal(t *testing.T) {
	in, err := os.ReadFile("with-anchors.yaml")
	assert.NoError(t, err)
	out, err := os.ReadFile("without-anchors.yaml")
	assert.NoError(t, err)
	tests := []struct {
		name     string
		document string
		want     string
		wantErr  bool
	}{{
		name:     "ok",
		document: string(in),
		want:     string(out),
		wantErr:  false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Remarshal([]byte(tt.document))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, strings.ReplaceAll(tt.want, "\r\n", "\n"), strings.ReplaceAll(string(got), "\r\n", "\n"))
		})
	}
}
