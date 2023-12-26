package registry

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeConfigRegistry map[string]*kubernetes.Clientset

func New() *KubeConfigRegistry {
	return new(KubeConfigRegistry)
}

func (r *KubeConfigRegistry) AddToRegistry(name, kubeconfigPath string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	(*r)[name] = clientset
	return nil
}

func (r KubeConfigRegistry) GetFromRegistry(name string) (*kubernetes.Clientset, bool) {
	clientset, exists := r[name]
	return clientset, exists
}
