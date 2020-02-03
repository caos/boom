package auth

import (
	"errors"
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetGoogleAuthConfig(spec *toolsetsv1beta1.GoogleAuth) (map[string]string, error) {
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
