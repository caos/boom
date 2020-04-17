package auth

import (
	"github.com/caos/boom/internal/helper"
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

func GetGoogleAuthConfig(spec *toolsetsv1beta1.GrafanaGoogleAuth) (map[string]string, error) {
	clientID, err := helper.GetSecretValue(spec.ClientID, spec.ExistingClientIDSecret)
	if err != nil {
		return nil, err
	}

	clientSecret, err := helper.GetSecretValue(spec.ClientSecret, spec.ExistingClientSecretSecret)
	if err != nil {
		return nil, err
	}

	domains := strings.Join(spec.AllowedDomains, " ")

	return map[string]string{
		"enabled":         "true",
		"client_id":       string(clientID),
		"client_secret":   string(clientSecret),
		"scopes":          "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
		"auth_url":        "https://accounts.google.com/o/oauth2/auth",
		"token_url":       "https://accounts.google.com/o/oauth2/token",
		"allowed_domains": domains,
		"allow_sign_up":   "true",
	}, nil
}
