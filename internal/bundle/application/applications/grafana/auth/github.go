package auth

import (
	"github.com/caos/boom/internal/helper"
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

func GetGithubAuthConfig(spec *toolsetsv1beta1.GrafanaGithubAuth) (map[string]string, error) {
	clientID, err := helper.GetSecretValue(spec.ClientID, spec.ExistingClientIDSecret)
	if err != nil {
		return nil, err
	}

	clientSecret, err := helper.GetSecretValue(spec.ClientSecret, spec.ExistingClientSecretSecret)
	if err != nil {
		return nil, err
	}

	teamIds := strings.Join(spec.TeamIDs, " ")
	allowedOrganizations := strings.Join(spec.AllowedOrganizations, " ")

	return map[string]string{
		"enabled":               "true",
		"allow_sign_up":         "true",
		"client_id":             clientID,
		"client_secret":         clientSecret,
		"scopes":                "user:email,read:org",
		"auth_url":              "https://github.com/login/oauth/authorize",
		"token_url":             "https://github.com/login/oauth/access_token",
		"api_url":               "https://api.github.com/user",
		"team_ids":              teamIds,
		"allowed_organizations": allowedOrganizations,
	}, nil
}
