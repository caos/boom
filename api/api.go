package api

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/secret"
)

func NewV1beta1Toolset(masterkey string) *v1beta1.Toolset {
	desiredKind := &v1beta1.Toolset{
		Spec: &v1beta1.ToolsetSpec{
			Grafana: &v1beta1.Grafana{
				Admin: &v1beta1.Admin{
					Username: &secret.Secret{Masterkey: masterkey},
					Password: &secret.Secret{Masterkey: masterkey},
				},
				Auth: &v1beta1.GrafanaAuth{
					Google: &v1beta1.GrafanaGoogleAuth{
						ClientID:     &secret.Secret{Masterkey: masterkey},
						ClientSecret: &secret.Secret{Masterkey: masterkey},
					},
					Github: &v1beta1.GrafanaGithubAuth{
						ClientID:     &secret.Secret{Masterkey: masterkey},
						ClientSecret: &secret.Secret{Masterkey: masterkey},
					},
					Gitlab: &v1beta1.GrafanaGitlabAuth{
						ClientID:     &secret.Secret{Masterkey: masterkey},
						ClientSecret: &secret.Secret{Masterkey: masterkey},
					},
					GenericOAuth: &v1beta1.GrafanaGenericOAuth{
						ClientID:     &secret.Secret{Masterkey: masterkey},
						ClientSecret: &secret.Secret{Masterkey: masterkey},
					},
				},
			},
			Argocd: &v1beta1.Argocd{
				Auth: &v1beta1.ArgocdAuth{
					OIDC: &v1beta1.ArgocdOIDC{
						ClientID:     &secret.Secret{Masterkey: masterkey},
						ClientSecret: &secret.Secret{Masterkey: masterkey},
					},
					GithubConnector: &v1beta1.ArgocdGithubConnector{
						Config: &v1beta1.ArgocdGithubConfig{
							ClientID:     &secret.Secret{Masterkey: masterkey},
							ClientSecret: &secret.Secret{Masterkey: masterkey},
						},
					},
					GitlabConnector: &v1beta1.ArgocdGitlabConnector{
						Config: &v1beta1.ArgocdGitlabConfig{
							ClientID:     &secret.Secret{Masterkey: masterkey},
							ClientSecret: &secret.Secret{Masterkey: masterkey},
						},
					},
					GoogleConnector: &v1beta1.ArgocdGoogleConnector{
						Config: &v1beta1.ArgocdGoogleConfig{
							ClientID:           &secret.Secret{Masterkey: masterkey},
							ClientSecret:       &secret.Secret{Masterkey: masterkey},
							ServiceAccountJSON: &secret.Secret{Masterkey: masterkey},
						},
					},
				},
				//Repositories: []*v1beta1.ArgocdRepository{{
				//	Username:    &secret.Secret{Masterkey: masterkey},
				//	Password:    &secret.Secret{Masterkey: masterkey},
				//	Certificate: &secret.Secret{Masterkey: masterkey},
				//}},
				//CustomImage: &v1beta1.ArgocdCustomImage{
				//	GopassStores: []*v1beta1.ArgocdGopassStore{{
				//		SSHKey: &secret.Secret{Masterkey: masterkey},
				//		GPGKey: &secret.Secret{Masterkey: masterkey},
				//	},
				//	},
				//},
			},
		},
	}
	return desiredKind
}
