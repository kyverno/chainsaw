package rest

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func DefaultConfig(overrides clientcmd.ConfigOverrides) (*rest.Config, error) {
	return load(clientcmd.NewDefaultClientConfigLoadingRules(), overrides)
}

func Config(kubeconfigPath string, overrides clientcmd.ConfigOverrides) (*rest.Config, error) {
	loader := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	return load(loader, overrides)
}

func load(loader clientcmd.ClientConfigLoader, overrides clientcmd.ConfigOverrides) (*rest.Config, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, &overrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	config.QPS = 300
	config.Burst = 300
	return config, nil
}
