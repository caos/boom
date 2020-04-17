package auth

import (
	"github.com/caos/boom/internal/helper"
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

func GetGitlabAuthConfig(spec *toolsetsv1beta1.GrafanaGitlabAuth) (map[string]string, error) {
	clientID, err := helper.GetSecretValue(spec.ClientID, spec.ExistingClientIDSecret)
	if err != nil {
		return nil, err
	}

	clientSecret, err := helper.GetSecretValue(spec.ClientSecret, spec.ExistingClientSecretSecret)
	if err != nil {
		return nil, err
	}

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
