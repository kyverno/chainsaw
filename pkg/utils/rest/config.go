package rest

import (
	"io"
	"os"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	api "k8s.io/client-go/tools/clientcmd/api/v1"
)

func Config(overrides clientcmd.ConfigOverrides) (*rest.Config, error) {
	// TODO: in cluster config
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &overrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	config.QPS = 300
	config.Burst = 300
	config.WarningHandler = rest.NewWarningWriter(os.Stderr, rest.WarningWriterOptions{Deduplicate: true})
	return config, nil
}

// Kubeconfig converts a rest.Config into a YAML kubeconfig and writes it to w
func Kubeconfig(cfg *rest.Config, w io.Writer) error {
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
			Command:    cfg.ExecProvider.Command,
			Args:       cfg.ExecProvider.Args,
			APIVersion: cfg.ExecProvider.APIVersion,
			Env:        []api.ExecEnvVar{},
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
		CurrentContext: "cluster",
		Clusters: []api.NamedCluster{
			{
				Name: "cluster",
				Cluster: api.Cluster{
					Server:                   cfg.Host,
					CertificateAuthorityData: cfg.TLSClientConfig.CAData,
					InsecureSkipTLSVerify:    cfg.TLSClientConfig.Insecure,
				},
			},
		},
		Contexts: []api.NamedContext{
			{
				Name: "cluster",
				Context: api.Context{
					Cluster:  "cluster",
					AuthInfo: "user",
				},
			},
		},
		AuthInfos: []api.NamedAuthInfo{
			{
				Name: "user",
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
