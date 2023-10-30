package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
)

func TestNew(t *testing.T) {
	cfg := &rest.Config{
		Host:    "http://localhost",
		APIPath: "/api",
	}

	client, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client, err = New(nil)
	assert.Error(t, err)
	assert.Nil(t, client)
}
