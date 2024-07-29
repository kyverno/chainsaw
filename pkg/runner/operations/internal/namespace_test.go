package internal

import (
	"errors"
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	tnamespacer "github.com/kyverno/chainsaw/pkg/engine/namespacer/testing"
	"github.com/stretchr/testify/assert"
)

func TestApplyNamespacer(t *testing.T) {
	tests := []struct {
		name       string
		namespacer namespacer.Namespacer
		obj        client.Object
		wantErr    bool
	}{{
		name:       "nil namespacer",
		namespacer: nil,
		obj:        nil,
		wantErr:    false,
	}, {
		name: "namespacer err",
		namespacer: &tnamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return errors.New("namespacer err")
			},
		},
		obj:     nil,
		wantErr: true,
	}, {
		name: "namespacer ok",
		namespacer: &tnamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		obj:     nil,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApplyNamespacer(tt.namespacer, nil, tt.obj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
