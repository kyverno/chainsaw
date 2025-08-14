package logging

import (
	"context"
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestLog(t *testing.T) {
	called := false
	_operation := Create
	_status := BeginStatus
	_obj := &unstructured.Unstructured{}
	_color := color.New(color.FgBlue)
	_args := []fmt.Stringer{Section("arg1"), Section("arg2")}
	var f LoggerFunc = func(_ context.Context, operation Operation, status Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
		assert.Equal(t, _operation, operation)
		assert.Equal(t, _status, status)
		assert.Equal(t, _obj, obj)
		assert.Equal(t, _color, color)
		assert.Equal(t, _args, args)
		called = true
	}
	ctx := context.Background()
	assert.NotNil(t, ctx)
	assert.Nil(t, getLogger(ctx))
	Log(ctx, _operation, _status, _obj, _color, _args...)
	assert.False(t, called)
	ctx = WithLogger(ctx, f)
	Log(ctx, _operation, _status, _obj, _color, _args...)
	assert.True(t, called)
}
