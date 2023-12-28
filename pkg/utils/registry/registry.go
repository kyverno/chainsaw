package registry

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeConfigRegistry map[string]client.Client

func New() *KubeConfigRegistry {
	return &KubeConfigRegistry{}
}

func (r *KubeConfigRegistry) AddToRegistry(name, kubeconfigPath string) error {
	if r == nil {
		return fmt.Errorf("kubeconfig registry is nil")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	client, err := client.New(config)
	if err != nil {
		return err
	}
	(*r)[name] = client
	return nil
}

func (r *KubeConfigRegistry) GetFromRegistry(name string) (client.Client, bool) {
	if r == nil {
		return nil, false
	}
	client, exists := (*r)[name]
	return client, exists
}
