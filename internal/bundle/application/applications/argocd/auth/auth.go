package auth

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

type connectorsStruct struct {
	Connectors []*connector
}

type connector struct {
	Type   string
	Name   string
	ID     string
	Config interface{}
}

func GetDexConfigFromSpec(logger logging.Logger, spec *toolsetsv1beta1.Argocd) *connectorsStruct {
	logFields := map[string]interface{}{
		"logID":       "AUTH-yYqnPDhTdTjQWiJ",
		"application": "argocd",
	}

	connectors := make([]*connector, 0)
	if spec.Auth.RootURL == "" {
		logger.WithFields(logFields).Info("No auth connectors configured as no rootUrl is defined")
		return nil
	}
	redirect := strings.Join([]string{spec.Auth.RootURL, "/api/dex/callback"}, "")

	if spec.Auth.GithubConnector != nil {
		github, err := getGithub(spec.Auth.GithubConnector, redirect)
		if err == nil {
			connectors = append(connectors, &connector{
				Name:   spec.Auth.GithubConnector.Name,
				ID:     spec.Auth.GithubConnector.ID,
				Type:   "github",
				Config: github,
			})
		} else {
			logger.WithFields(logFields).Error(errors.Wrap(err, "Error while creating configuration for github connector"))
		}
	}

	if spec.Auth.GitlabConnector != nil {
		gitlab, err := getGitlab(spec.Auth.GitlabConnector, redirect)
		if err == nil {
			connectors = append(connectors, &connector{
				Name:   spec.Auth.GitlabConnector.Name,
				ID:     spec.Auth.GitlabConnector.ID,
				Type:   "gitlab",
				Config: gitlab,
			})
		} else {
			logger.WithFields(logFields).Error(errors.Wrap(err, "Error while creating configuration for gitlab connector"))
		}
	}

	if spec.Auth.GoogleConnector != nil {
		google, err := getGoogle(spec.Auth.GoogleConnector, redirect)
		if err == nil {
			connectors = append(connectors, &connector{
				Name:   spec.Auth.GoogleConnector.Name,
				ID:     spec.Auth.GoogleConnector.ID,
				Type:   "oidc",
				Config: google,
			})
		} else {
			logger.WithFields(logFields).Error(errors.Wrap(err, "Error while creating configuration for google connector"))
		}
	}

	if len(connectors) > 0 {
		logFields["connectors"] = len(connectors)
		logger.WithFields(logFields).Debug("Created dex configuration")
		return &connectorsStruct{Connectors: connectors}
	}
	return nil
}
