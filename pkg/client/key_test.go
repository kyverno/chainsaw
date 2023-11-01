package client

import (
	"testing"

	"github.com/fatih/color"
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

func TestColouredName(t *testing.T) {
	disabled := color.New(color.FgBlue)
	disabled.DisableColor()
	enabled := color.New(color.FgBlue)
	enabled.EnableColor()
	tests := []struct {
		name  string
		key   ctrlclient.ObjectKey
		color *color.Color
		want  string
	}{{
		name:  "empty",
		key:   ctrlclient.ObjectKey{},
		color: nil,
		want:  "*",
	}, {
		name:  "name only",
		key:   ctrlclient.ObjectKey{Name: "test-name"},
		color: nil,
		want:  "test-name",
	}, {
		name:  "name and namespace",
		key:   ctrlclient.ObjectKey{Name: "test-name", Namespace: "test-namespace"},
		color: nil,
		want:  "test-namespace/test-name",
	}, {
		name:  "empty",
		key:   ctrlclient.ObjectKey{},
		color: disabled,
		want:  "*",
	}, {
		name:  "name only",
		key:   ctrlclient.ObjectKey{Name: "test-name"},
		color: disabled,
		want:  "test-name",
	}, {
		name:  "name and namespace",
		key:   ctrlclient.ObjectKey{Name: "test-name", Namespace: "test-namespace"},
		color: disabled,
		want:  "test-namespace/test-name",
	}, {
		name:  "empty",
		key:   ctrlclient.ObjectKey{},
		color: enabled,
		want:  "\x1b[34m*\x1b[0m",
	}, {
		name:  "name only",
		key:   ctrlclient.ObjectKey{Name: "test-name"},
		color: enabled,
		want:  "\x1b[34mtest-name\x1b[0m",
	}, {
		name:  "name and namespace",
		key:   ctrlclient.ObjectKey{Name: "test-name", Namespace: "test-namespace"},
		color: enabled,
		want:  "\x1b[34mtest-namespace\x1b[0m/\x1b[34mtest-name\x1b[0m",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColouredName(tt.key, tt.color)
			assert.Equal(t, tt.want, got)
		})
	}
}
