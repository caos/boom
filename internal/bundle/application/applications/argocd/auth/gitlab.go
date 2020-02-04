package auth

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
)

type gitlab struct {
	ClientID     string   `yaml:"clientID,omitempty"`
	ClientSecret string   `yaml:"clientSecret,omitempty"`
	RedirectURI  string   `yaml:"redirectURI,omitempty"`
	BaseURL      string   `yaml:"baseURL,omitempty"`
	Groups       []string `yaml:"groups,omitempty"`
	UseLoginAsID bool     `yaml:"useLoginAsID,omitempty"`
}

func GetGitlab(spec *toolsetsv1beta1.ArgocdGitlabConnector) (interface{}, error) {
	secret, err := helper.GetSecret(spec.Config.SecretName, "caos-system")
	if err != nil {
		return "", err
	}

	clientID := string(secret.Data[spec.Config.ClientIDKey])
	clientSecret := string(secret.Data[spec.Config.ClientSecretKey])

	var groups []string
	if len(spec.Config.Groups) > 0 {
		groups = make([]string, len(spec.Config.Groups))
		for k, v := range spec.Config.Groups {
			groups[k] = v
		}
	}

	gitlab := &gitlab{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  spec.Config.RedirectURI,
		Groups:       groups,
		UseLoginAsID: spec.Config.UseLoginAsID,
	}

	return gitlab, nil
}
