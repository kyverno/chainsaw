package rest

import (
	"io"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	api "k8s.io/client-go/tools/clientcmd/api/v1"
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

func Save(cfg *rest.Config, w io.Writer) error {
	var authProvider *api.AuthProviderConfig
	var execConfig *api.ExecConfig
	if cfg.AuthProvider != nil {
		authProvider = &api.AuthProviderConfig{
			Name:   cfg.AuthProvider.Name,
			Config: cfg.AuthProvider.Config,
		}
	}
	if cfg.ExecProvider != nil {
		execConfig = &api.ExecConfig{
			Command:         cfg.ExecProvider.Command,
			Args:            cfg.ExecProvider.Args,
			APIVersion:      cfg.ExecProvider.APIVersion,
			Env:             []api.ExecEnvVar{},
			InteractiveMode: api.ExecInteractiveMode(cfg.ExecProvider.InteractiveMode),
		}
		for _, envVar := range cfg.ExecProvider.Env {
			execConfig.Env = append(execConfig.Env, api.ExecEnvVar{
				Name:  envVar.Name,
				Value: envVar.Value,
			})
		}
	}
	err := rest.LoadTLSFiles(cfg)
	if err != nil {
		return err
	}
	return json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil).Encode(&api.Config{
		CurrentContext: "chainsaw",
		Clusters: []api.NamedCluster{
			{
				Name: "chainsaw",
				Cluster: api.Cluster{
					Server:                   cfg.Host,
					CertificateAuthorityData: cfg.TLSClientConfig.CAData,
					InsecureSkipTLSVerify:    cfg.TLSClientConfig.Insecure,
				},
			},
		},
		Contexts: []api.NamedContext{
			{
				Name: "chainsaw",
				Context: api.Context{
					Cluster:  "chainsaw",
					AuthInfo: "chainsaw",
				},
			},
		},
		AuthInfos: []api.NamedAuthInfo{
			{
				Name: "chainsaw",
				AuthInfo: api.AuthInfo{
					ClientCertificateData: cfg.TLSClientConfig.CertData,
					ClientKeyData:         cfg.TLSClientConfig.KeyData,
					Token:                 cfg.BearerToken,
					Username:              cfg.Username,
					Password:              cfg.Password,
					Impersonate:           cfg.Impersonate.UserName,
					ImpersonateGroups:     cfg.Impersonate.Groups,
					ImpersonateUserExtra:  cfg.Impersonate.Extra,
					AuthProvider:          authProvider,
					Exec:                  execConfig,
				},
			},
		},
	}, w)
}
