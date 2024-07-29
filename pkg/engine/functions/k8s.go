package functions

import (
	"context"
	"errors"

	"github.com/kyverno/chainsaw/pkg/client"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
)

func jpKubernetesResourceExists(arguments []any) (any, error) {
	var client client.Client
	var apiVersion, kind string
	if err := getArg(arguments, 0, &client); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return nil, err
	}
	mapper := client.RESTMapper()
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return nil, err
	}
	gvk := gv.WithKind(kind)
	if _, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version); err == nil {
		return true, nil
	} else if meta.IsNoMatchError(err) {
		return false, nil
	} else {
		return nil, err
	}
}

func jpKubernetesExists(arguments []any) (any, error) {
	var apiVersion, kind string
	var key client.ObjectKey
	var client client.Client
	if err := getArg(arguments, 0, &client); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 3, &key.Namespace); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 4, &key.Name); err != nil {
		return nil, err
	}
	var obj unstructured.Unstructured
	obj.SetAPIVersion(apiVersion)
	obj.SetKind(kind)
	err := client.Get(context.TODO(), key, &obj)
	if err == nil {
		return true, nil
	}
	if apierrors.IsNotFound(err) {
		return false, nil
	}
	return nil, err
}

func jpKubernetesGet(arguments []any) (any, error) {
	var apiVersion, kind string
	var key client.ObjectKey
	var client client.Client
	if err := getArg(arguments, 0, &client); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 3, &key.Namespace); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 4, &key.Name); err != nil {
		return nil, err
	}
	var obj unstructured.Unstructured
	obj.SetAPIVersion(apiVersion)
	obj.SetKind(kind)
	if err := client.Get(context.TODO(), key, &obj); err != nil {
		return nil, err
	}
	return obj.UnstructuredContent(), nil
}

func jpKubernetesList(arguments []any) (any, error) {
	var c client.Client
	var apiVersion, kind, namespace string
	if err := getArg(arguments, 0, &c); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return nil, err
	}
	if len(arguments) == 4 {
		if err := getArg(arguments, 3, &namespace); err != nil {
			return nil, err
		}
	}
	var list unstructured.UnstructuredList
	list.SetAPIVersion(apiVersion)
	list.SetKind(kind)
	var listOptions []client.ListOption
	if namespace != "" {
		listOptions = append(listOptions, client.InNamespace(namespace))
	}
	if err := c.List(context.TODO(), &list, listOptions...); err != nil {
		return nil, err
	}
	return list.UnstructuredContent(), nil
}

func jpKubernetesServerVersion(arguments []any) (any, error) {
	var config *rest.Config
	if err := getArg(arguments, 0, &config); err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("rest config is nil")
	}
	client, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	version, err := client.ServerVersion()
	if err != nil {
		return nil, err
	}
	return version, nil
}
