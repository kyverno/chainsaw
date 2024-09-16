package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/kyverno/chainsaw/pkg/engine/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/engine/logging/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestGetLogger(t *testing.T) {
	assert.Nil(t, GetLogger(context.TODO(), nil))
	{
		logger := &tlogging.FakeLogger{}
		ctx := logging.IntoContext(context.TODO(), logger)
		assert.NotNil(t, GetLogger(ctx, nil))
	}
	{
		logger := &tlogging.FakeLogger{}
		ctx := logging.IntoContext(context.TODO(), logger)
		obj := unstructured.Unstructured{}
		obj.SetAPIVersion("foo/v1")
		obj.SetKind("bar")
		obj.SetName("hello")
		obj.SetNamespace("chainsaw")
		l := GetLogger(ctx, &obj)
		l.Log("aaa", "bbb", nil)
		assert.Equal(t, []string{"aaa: bbb - []"}, logger.Logs)
	}
}

func TestLogStart(t *testing.T) {
	logger := &tlogging.FakeLogger{}
	ctx := logging.IntoContext(context.TODO(), logger)
	obj := unstructured.Unstructured{}
	obj.SetAPIVersion("foo/v1")
	obj.SetKind("bar")
	obj.SetName("hello")
	obj.SetNamespace("chainsaw")
	l := GetLogger(ctx, &obj)
	LogStart(l, "aaa")
	assert.Equal(t, []string{"aaa: RUN - []"}, logger.Logs)
}

func TestLogEnd(t *testing.T) {
	{
		logger := &tlogging.FakeLogger{}
		ctx := logging.IntoContext(context.TODO(), logger)
		obj := unstructured.Unstructured{}
		obj.SetAPIVersion("foo/v1")
		obj.SetKind("bar")
		obj.SetName("hello")
		obj.SetNamespace("chainsaw")
		l := GetLogger(ctx, &obj)
		LogEnd(l, "aaa", nil)
		assert.Equal(t, []string{"aaa: DONE - []"}, logger.Logs)
	}
	{
		logger := &tlogging.FakeLogger{}
		ctx := logging.IntoContext(context.TODO(), logger)
		obj := unstructured.Unstructured{}
		obj.SetAPIVersion("foo/v1")
		obj.SetKind("bar")
		obj.SetName("hello")
		obj.SetNamespace("chainsaw")
		l := GetLogger(ctx, &obj)
		LogEnd(l, "aaa", errors.New("some error"))
		assert.Equal(t, []string{"aaa: ERROR - [=== ERROR\nsome error]"}, logger.Logs)
	}
}
