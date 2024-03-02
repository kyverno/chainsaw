package kubectl

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Describe(collector *v1alpha1.Describe) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, errors.New("collector is null")
	}
	if collector.Resource == "" {
		return nil, errors.New("a resource must be specified")
	}
	if collector.Name != "" && collector.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	cmd := v1alpha1.Command{
		Cluster:    collector.Cluster,
		Timeout:    collector.Timeout,
		Entrypoint: "kubectl",
		Args:       []string{"describe", collector.Resource},
	}
	if collector.Name != "" {
		cmd.Args = append(cmd.Args, collector.Name)
	}
	if collector.Selector != "" {
		cmd.Args = append(cmd.Args, "-l", collector.Selector)
	}
	config, _ := readKubeConfig("")
	isClusterScoped, _ := isClusterScoped(collector.Resource, config)
	if !isClusterScoped {
		// TODO: what if cluster scoped resource ?
		namespace := collector.Namespace
		if collector.Namespace == "" {
			namespace = "$NAMESPACE"
		}
		cmd.Args = append(cmd.Args, "-n", namespace)
	}
	if collector.ShowEvents != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--show-events=%t", *collector.ShowEvents))
	}
	return &cmd, nil
}

func resourceToGroupVersion(kind string, config *rest.Config) (schema.GroupVersion, error) {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return schema.GroupVersion{}, err
	}

	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return schema.GroupVersion{}, err
	}

	for _, apiResourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			return schema.GroupVersion{}, err
		}

		for _, apiResource := range apiResourceList.APIResources {
			if apiResource.Kind == kind {
				return gv, nil
			}
		}
	}

	return schema.GroupVersion{}, fmt.Errorf("resource kind %s not found", kind)
}

func isClusterScoped(kind string, config *rest.Config) (bool, error) {
	groupVersion, err := resourceToGroupVersion(kind, config)
	if err != nil {
		return false, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return false, err
	}

	apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(groupVersion.String())
	if err != nil {
		return false, err
	}

	for _, apiResource := range apiResourceList.APIResources {
		if apiResource.Kind == kind && !apiResource.Namespaced {
			return true, nil
		}
	}

	return false, nil
}

func readKubeConfig(kubeConfigPath string) (*rest.Config, error) {
	if kubeConfigPath == "" {
		if envPath := os.Getenv("KUBECONFIG"); envPath != "" {
			kubeConfigPath = envPath
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			kubeConfigPath = filepath.Join(home, ".kube", "config")
		}
	}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath}, &clientcmd.ConfigOverrides{})
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return restConfig, nil
}
