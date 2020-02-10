package auth

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type oidc struct {
	Name                   string            `yaml:"name,omitempty"`
	Issuer                 string            `yaml:"issuer,omitempty"`
	ClientID               string            `yaml:"clientID,omitempty"`
	ClientSecret           string            `yaml:"clientSecret,omitempty"`
	RequestedScopes        []string          `yaml:"requestedScopes,omitempty"`
	RequestedIDTokenClaims map[string]*Claim `yaml:"requestedIDTokenClaims,omitempty"`
}
type Claim struct {
	Essential bool     `yaml:"essential,omitempty"`
	Values    []string `yaml:"values,omitempty"`
}

func GetOIDC(spec *toolsetsv1beta1.ArgocdOIDC) (string, error) {
	secret, err := helper.GetSecret(spec.SecretName, "caos-system")
	if err != nil {
		return "", err
	}
	clientID := string(secret.Data[spec.ClientIDKey])
	clientSecret := string(secret.Data[spec.ClientSecretKey])

	var claims map[string]*Claim
	if len(spec.RequestedIDTokenClaims) > 0 {
		claims = make(map[string]*Claim, 0)
		for k, v := range spec.RequestedIDTokenClaims {
			claims[k] = &Claim{
				Essential: v.Essential,
				Values:    v.Values,
			}
		}
	}

	oidc := &oidc{
		Name:                   spec.Name,
		Issuer:                 spec.Issuer,
		ClientID:               clientID,
		ClientSecret:           clientSecret,
		RequestedScopes:        spec.RequestedScopes,
		RequestedIDTokenClaims: claims,
	}

	data, err := yaml.Marshal(oidc)
	return string(data), errors.Wrap(err, "Error while generating argocd oidc configuration")
}
