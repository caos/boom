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

func GetOIDC(spec *toolsetsv1beta1.ArgocdAuth) (string, error) {
	if spec == nil || spec.OIDC == nil {
		return "", nil
	}

	clientID, err := helper.GetSecretValue(spec.OIDC.ClientID, spec.OIDC.ExistingClientIDSecret)
	if err != nil {
		return "", err
	}

	clientSecret, err := helper.GetSecretValue(spec.OIDC.ClientSecret, spec.OIDC.ExistingClientSecretSecret)
	if err != nil {
		return "", err
	}

	var claims map[string]*Claim
	if len(spec.OIDC.RequestedIDTokenClaims) > 0 {
		claims = make(map[string]*Claim, 0)
		for k, v := range spec.OIDC.RequestedIDTokenClaims {
			claims[k] = &Claim{
				Essential: v.Essential,
				Values:    v.Values,
			}
		}
	}

	oidc := &oidc{
		Name:                   spec.OIDC.Name,
		Issuer:                 spec.OIDC.Issuer,
		ClientID:               clientID,
		ClientSecret:           clientSecret,
		RequestedScopes:        spec.OIDC.RequestedScopes,
		RequestedIDTokenClaims: claims,
	}

	data, err := yaml.Marshal(oidc)
	return string(data), errors.Wrap(err, "Error while generating argocd oidc configuration")
}
