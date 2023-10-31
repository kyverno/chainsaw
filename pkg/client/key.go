package client

import (
	"fmt"

	"github.com/kyverno/kyverno/ext/output/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func ObjectKey(obj metav1.Object) ctrlclient.ObjectKey {
	return ctrlclient.ObjectKey{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
}

func ObjectName(obj metav1.Object) string {
	name := obj.GetName()
	if ns := obj.GetNamespace(); ns != "" {
		name = ns + "/" + name
	}
	return name
}

func Name(key ctrlclient.ObjectKey) string {
	name := key.Name
	if name == "" {
		name = "*"
	}
	if key.Namespace != "" {
		name = key.Namespace + "/" + name
	}
	return name
}

func ColouredName(key ctrlclient.ObjectKey, color *color.Color) string {
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
