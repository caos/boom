package config

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/config/auth"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/config/plugin"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/config/repository"
	"github.com/caos/orbiter/mntr"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Repositories            string `yaml:"repositories,omitempty"`
	Connectors              string `yaml:"connectors,omitempty"`
	OIDC                    string `yaml:"oidc,omitempty"`
	ConfigManagementPlugins string `yaml:"configManagementPlugins,omitempty"`
}

func GetFromSpec(monitor mntr.Monitor, spec *toolsetsv1beta1.Argocd) *Config {
	conf := &Config{}

	dexconfig := auth.GetDexConfigFromSpec(monitor, spec)
	data, err := yaml.Marshal(dexconfig)
	if err == nil {
		conf.Connectors = string(data)
	}
	repos := repository.GetFromSpec(monitor, spec)
	data2, err := yaml.Marshal(repos)
	if err == nil {
		conf.Repositories = string(data2)
	}

	oidc, err := auth.GetOIDC(spec.Auth)
	if err == nil {
		conf.OIDC = oidc
	}

	if spec.CustomImage != nil {
		plugins := make([]*plugin.Plugin, 0)
		init := &plugin.Command{
			Command: []string{"gopass", "sync"},
		}
		generate := &plugin.Command{
			Command: []string{"sh", "-c"},
			Args:    []string{"kustomize build && ./secrets.yaml.sh "},
		}
		plugins = append(plugins, plugin.New("getSecrets", init, generate))

		pluginsYaml, err := yaml.Marshal(plugins)
		if err == nil {
			conf.ConfigManagementPlugins = string(pluginsYaml)
		}
	}

	return conf
}
