package auth

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

type connector struct {
	Type   string
	Name   string
	ID     string
	Config interface{}
}

func GetDexConfigFromSpec(spec *toolsetsv1beta1.Argocd) []*connector {

	dex := make([]*connector, 0)

	if spec.Auth.GithubConnector != nil {
		github, err := getGithub(spec.Auth.GithubConnector)
		if err == nil {
			dex = append(dex, &connector{
				Name:   spec.Auth.GithubConnector.Name,
				ID:     spec.Auth.GithubConnector.ID,
				Type:   "github",
				Config: github,
			})
		}
	}

	if spec.Auth.GitlabConnector != nil {
		gitlab, err := getGitlab(spec.Auth.GitlabConnector)
		if err == nil {
			dex = append(dex, &connector{
				Name:   spec.Auth.GitlabConnector.Name,
				ID:     spec.Auth.GitlabConnector.ID,
				Type:   "gitlab",
				Config: gitlab,
			})
		}
	}

	if spec.Auth.GoogleConnector != nil {
		google, err := getGoogle(spec.Auth.GoogleConnector)
		if err == nil {
			dex = append(dex, &connector{
				Name:   spec.Auth.GoogleConnector.Name,
				ID:     spec.Auth.GoogleConnector.ID,
				Type:   "google",
				Config: google,
			})
		}
	}

	return dex
}
