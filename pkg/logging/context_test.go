package logging

import (
	"context"
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestWithSink(t *testing.T) {
	var f SinkFunc = func(string, string, string, Operation, Status, client.Object, *color.Color, ...fmt.Stringer) {
	}
	ctx := context.Background()
	assert.NotNil(t, ctx)
	ctx = WithSink(ctx, f)
	assert.NotNil(t, ctx)
}

func Test_getSink(t *testing.T) {
	var f SinkFunc = func(string, string, string, Operation, Status, client.Object, *color.Color, ...fmt.Stringer) {
	}
	ctx := context.Background()
	assert.NotNil(t, ctx)
	assert.Nil(t, getSink(ctx))
	ctx = WithSink(ctx, f)
	assert.NotNil(t, ctx)
	assert.NotNil(t, getSink(ctx))
}

func TestWithLogger(t *testing.T) {
	var f LoggerFunc = func(context.Context, Operation, Status, client.Object, *color.Color, ...fmt.Stringer) {
	}
	ctx := context.Background()
	assert.NotNil(t, ctx)
	ctx = WithLogger(ctx, f)
	assert.NotNil(t, ctx)
}

func Test_getLogger(t *testing.T) {
	var f LoggerFunc = func(context.Context, Operation, Status, client.Object, *color.Color, ...fmt.Stringer) {
	}
	ctx := context.Background()
	assert.NotNil(t, ctx)
	assert.Nil(t, getLogger(ctx))
	ctx = WithLogger(ctx, f)
	assert.NotNil(t, ctx)
	assert.NotNil(t, getLogger(ctx))
}
