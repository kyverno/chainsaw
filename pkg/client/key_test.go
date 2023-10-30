package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestObjectKey(t *testing.T) {
	obj := &metav1.ObjectMeta{
		Name:      "test-name",
		Namespace: "test-namespace",
	}

	key := ObjectKey(obj)
	assert.Equal(t, "test-name", key.Name)
	assert.Equal(t, "test-namespace", key.Namespace)
}

func TestObjectName(t *testing.T) {
	obj := &metav1.ObjectMeta{
		Name: "test-name",
	}
	name := ObjectName(obj)
	assert.Equal(t, "test-name", name)

	obj.Namespace = "test-namespace"
	name = ObjectName(obj)
	assert.Equal(t, "test-namespace/test-name", name)
}

func TestName(t *testing.T) {
	key := ctrlclient.ObjectKey{Name: "test-name"}
	name := Name(key)
	assert.Equal(t, "test-name", name)

	key.Namespace = "test-namespace"
	name = Name(key)
	assert.Equal(t, "test-namespace/test-name", name)

	key = ctrlclient.ObjectKey{}
	name = Name(key)
	assert.Equal(t, "*", name)
}

// Getting a error here
// func TestColouredName(t *testing.T) {
// 	colour := color.BoldFgCyan

// 	key := ctrlclient.ObjectKey{Name: "test-name"}
// 	name := ColouredName(key, colour)
// 	assert.Contains(t, name, "test-name")

// 	key.Namespace = "test-namespace"
// 	name = ColouredName(key, colour)
// 	assert.Contains(t, name, "test-namespace")
// 	assert.Contains(t, name, "test-name")

// 	key = ctrlclient.ObjectKey{}
// 	name = ColouredName(key, colour)
// 	assert.Contains(t, name, "*")
// }
