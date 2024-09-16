package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kyverno/pkg/ext/output/color"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
)

const PollInterval = 50 * time.Millisecond

func Key(obj metav1.Object) ObjectKey {
	return ObjectKey{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
}

func Name(key ObjectKey) string {
	return ColouredName(key, nil)
}

func ColouredName(key ObjectKey, color *color.Color) string {
	sprint := fmt.Sprint
	if color != nil {
		sprint = color.Sprint
	}
	name := key.Name
	if name == "" {
		name = "*"
	}
	name = sprint(name)
	if key.Namespace != "" {
		name = sprint(key.Namespace) + "/" + name
	}
	return name
}

func PatchObject(actual, expected runtime.Object) (runtime.Object, error) {
	if actual == nil || expected == nil {
		return nil, errors.New("actual and expected objects must not be nil")
	}
	actualMeta, err := meta.Accessor(actual)
	if err != nil {
		return nil, err
	}
	copy := expected.DeepCopyObject()
	expectedMeta, err := meta.Accessor(copy)
	if err != nil {
		return nil, err
	}
	expectedMeta.SetResourceVersion(actualMeta.GetResourceVersion())
	return copy, nil
}

func WaitForDeletion(ctx context.Context, client Client, object Object) error {
	key := Key(object)
	return wait.PollUntilContextCancel(ctx, PollInterval, true, func(ctx context.Context) (bool, error) {
		if err := client.Get(ctx, key, object); err != nil {
			if kerrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	})
}
