package auth

import (
	"errors"
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetGithubAuthConfig(spec *toolsetsv1beta1.GithubAuth) (map[string]string, error) {
	conf, err := helper.GetClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := helper.GetClientSet(conf)
	if err != nil {
		return nil, err
	}

	secret, err := clientset.CoreV1().Secrets("caos-system").Get(spec.SecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, errors.New("Secret not found")
	}

	clientID := string(secret.Data[spec.ClientIDKey])
	clientSecret := string(secret.Data[spec.ClientSecret])
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
