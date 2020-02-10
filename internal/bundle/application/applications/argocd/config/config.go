package config

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/config/auth"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/config/repository"
	"github.com/caos/orbiter/logging"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Repositories string `yaml:"repositories,omitempty"`
	Connectors   string `yaml:"connectors,omitempty"`
	OIDC         string `yaml:"oidc,omitempty"`
}

func GetFromSpec(logger logging.Logger, spec *toolsetsv1beta1.Argocd) *Config {
	conf := &Config{}

	dexconfig := auth.GetDexConfigFromSpec(logger, spec)
	data, err := yaml.Marshal(dexconfig)
	if err == nil {
		conf.Connectors = string(data)
	}
	repos := repository.GetFromSpec(logger, spec)
	data2, err := yaml.Marshal(repos)
	if err == nil {
		conf.Repositories = string(data2)
	}

	oidc, err := auth.GetOIDC(spec)
	if err == nil {
		conf.OIDC = oidc
	}

	return conf
}
