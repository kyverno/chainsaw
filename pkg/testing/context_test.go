package testing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	assert.Nil(t, FromContext(context.Background()))
	assert.NotNil(t, IntoContext(context.Background(), t))
	assert.Equal(t, t, FromContext(IntoContext(context.Background(), t)))
}
