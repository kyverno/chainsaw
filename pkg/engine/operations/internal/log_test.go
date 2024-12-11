package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestLogStart(t *testing.T) {
	logger := &mocks.Logger{}
	ctx := logging.WithLogger(context.TODO(), logger)
	obj := unstructured.Unstructured{}
	obj.SetAPIVersion("foo/v1")
	obj.SetKind("bar")
	obj.SetName("hello")
	obj.SetNamespace("chainsaw")
	LogStart(ctx, "aaa", &obj)
	assert.Equal(t, []string{"aaa: RUN - []"}, logger.Logs)
}

func TestLogEnd(t *testing.T) {
	{
		logger := &mocks.Logger{}
		ctx := logging.WithLogger(context.TODO(), logger)
		obj := unstructured.Unstructured{}
		obj.SetAPIVersion("foo/v1")
		obj.SetKind("bar")
		obj.SetName("hello")
		obj.SetNamespace("chainsaw")
		LogEnd(ctx, "aaa", &obj, nil)
		assert.Equal(t, []string{"aaa: DONE - []"}, logger.Logs)
	}
	{
		logger := &mocks.Logger{}
		ctx := logging.WithLogger(context.TODO(), logger)
		obj := unstructured.Unstructured{}
		obj.SetAPIVersion("foo/v1")
		obj.SetKind("bar")
		obj.SetName("hello")
		obj.SetNamespace("chainsaw")
		LogEnd(ctx, "aaa", &obj, errors.New("some error"))
		assert.Equal(t, []string{"aaa: ERROR - [=== ERROR\nsome error]"}, logger.Logs)
	}
}
