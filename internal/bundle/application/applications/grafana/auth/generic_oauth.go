package auth

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/clientgo"
)

func GetGenericOAuthConfig(spec *toolsetsv1beta1.GrafanaGenericOAuth) (map[string]string, error) {
	secret, err := clientgo.GetSecret(spec.SecretName, "caos-system")
	if err != nil {
		return nil, err
	}

	
	clientID := string(secret.Data[spec.ClientIDKey])
	clientSecret := string(secret.Data[spec.ClientSecretKey])
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
