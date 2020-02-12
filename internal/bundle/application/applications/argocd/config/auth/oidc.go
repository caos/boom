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

func GetOIDC(spec *toolsetsv1beta1.Argocd) (string, error) {
	if spec.Auth == nil || spec.Auth.OIDC == nil {
		return "", nil
	}

	secret, err := helper.GetSecret(spec.Auth.OIDC.SecretName, "caos-system")
	if err != nil {
		return "", err
	}
	clientID := string(secret.Data[spec.Auth.OIDC.ClientIDKey])
	clientSecret := string(secret.Data[spec.Auth.OIDC.ClientSecretKey])

	var claims map[string]*Claim
	if len(spec.Auth.OIDC.RequestedIDTokenClaims) > 0 {
		claims = make(map[string]*Claim, 0)
		for k, v := range spec.Auth.OIDC.RequestedIDTokenClaims {
			claims[k] = &Claim{
				Essential: v.Essential,
				Values:    v.Values,
			}
		}
	}

	oidc := &oidc{
		Name:                   spec.Auth.OIDC.Name,
		Issuer:                 spec.Auth.OIDC.Issuer,
		ClientID:               clientID,
		ClientSecret:           clientSecret,
		RequestedScopes:        spec.Auth.OIDC.RequestedScopes,
		RequestedIDTokenClaims: claims,
	}

	data, err := yaml.Marshal(oidc)
	return string(data), errors.Wrap(err, "Error while generating argocd oidc configuration")
}
