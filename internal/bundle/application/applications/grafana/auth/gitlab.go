package auth

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
)

func GetGitlabAuthConfig(spec *toolsetsv1beta1.GrafanaGitlabAuth) (map[string]string, error) {
	secret, err := helper.GetSecret(spec.SecretName, "caos-system")
	if err != nil {
		return nil, err
	}

	clientID := string(secret.Data[spec.ClientIDKey])
	clientSecret := string(secret.Data[spec.ClientSecret])
	allowedGroups := strings.Join(spec.AllowedGroups, " ")

	return map[string]string{
		"enabled":        "true",
		"allow_sign_up":  "false",
		"client_id":      clientID,
		"client_secret":  clientSecret,
		"scopes":         "api",
		"auth_url":       "https://gitlab.com/oauth/authorize",
		"token_url":      "https://gitlab.com/oauth/token",
		"api_url":        "https://gitlab.com/api/v4",
		"allowed_groups": allowedGroups,
	}, nil
}
