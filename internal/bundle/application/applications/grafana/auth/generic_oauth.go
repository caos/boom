package auth

import (
	"github.com/caos/boom/internal/helper"
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

func GetGenericOAuthConfig(spec *toolsetsv1beta1.GrafanaGenericOAuth) (map[string]string, error) {
	clientID, err := helper.GetSecretValue(spec.ClientID, spec.ExistingClientIDSecret)
	if err != nil {
		return nil, err
	}

	clientSecret, err := helper.GetSecretValue(spec.ClientSecret, spec.ExistingClientSecretSecret)
	if err != nil {
		return nil, err
	}

	allowedDomains := strings.Join(spec.AllowedDomains, " ")
	scopes := strings.Join(spec.Scopes, " ")

	return map[string]string{
		"enabled":         "true",
		"allow_sign_up":   "true",
		"client_id":       clientID,
		"client_secret":   clientSecret,
		"scopes":          scopes,
		"auth_url":        spec.AuthURL,
		"token_url":       spec.TokenURL,
		"api_url":         spec.APIURL,
		"allowed_domains": allowedDomains,
	}, nil
}
