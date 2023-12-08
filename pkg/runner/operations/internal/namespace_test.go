package internal

import (
	"errors"
	"testing"

	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	tnamespacer "github.com/kyverno/chainsaw/pkg/runner/namespacer/testing"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
			ApplyFn: func(_ client.Object, call int) error {
				return errors.New("namespacer err")
			},
		},
		obj:     nil,
		wantErr: true,
	}, {
		name: "namespacer ok",
		namespacer: &tnamespacer.FakeNamespacer{
			ApplyFn: func(_ client.Object, call int) error {
				return nil
			},
		},
		obj:     nil,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApplyNamespacer(tt.namespacer, tt.obj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
