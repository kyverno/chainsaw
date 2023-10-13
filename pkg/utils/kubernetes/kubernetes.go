package kubernetes

import (
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// This determine if the resource is clustered or namespaced
// if the resource is namespaced and doesn't have a namespace set, use the pet namespace
// if the resource is namespaced and has a namespace set, use the namespace set
// return namespace, name, error
func Namespaced(obj ctrlclient.Object, petNamespace string) (string, string, error) {

	return "", "", nil
}
