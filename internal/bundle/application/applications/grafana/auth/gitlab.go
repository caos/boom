package auth

import (
	"errors"
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetGitlabAuthConfig(spec *toolsetsv1beta1.GitlabAuth) (map[string]string, error) {
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
