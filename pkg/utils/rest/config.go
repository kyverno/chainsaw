package rest

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Config(overrides clientcmd.ConfigOverrides) (*rest.Config, error) {
	loader := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, &overrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	config.QPS = 300
	config.Burst = 300
	return config, nil
}
