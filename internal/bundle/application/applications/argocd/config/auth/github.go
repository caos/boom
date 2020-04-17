package auth

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
)

type github struct {
	ClientID      string `yaml:"clientID,omitempty"`
	ClientSecret  string `yaml:"clientSecret,omitempty"`
	RedirectURI   string `yaml:"redirectURI,omitempty"`
	Orgs          []*org `yaml:"orgs,omitempty"`
	LoadAllGroups bool   `yaml:"loadAllGroups,omitempty"`
	TeamNameField string `yaml:"teamNameField,omitempty"`
	UseLoginAsID  bool   `yaml:"useLoginAsID,omitempty"`
}
type org struct {
	Name  string   `yaml:"name,omitempty"`
	Teams []string `yaml:"teams,omitempty"`
}

func getGithub(spec *toolsetsv1beta1.ArgocdGithubConnector, redirect string) (interface{}, error) {
	clientID, err := helper.GetSecretValue(spec.Config.ClientID, spec.Config.ExistingClientIDSecret)
	if err != nil {
		return nil, err
	}

	clientSecret, err := helper.GetSecretValue(spec.Config.ClientSecret, spec.Config.ExistingClientSecretSecret)
	if err != nil {
		return nil, err
	}

	var orgs []*org
	if len(spec.Config.Orgs) > 0 {
		orgs = make([]*org, len(spec.Config.Orgs))
		for k, v := range spec.Config.Orgs {
			orgs[k] = &org{
				Name:  v.Name,
				Teams: v.Teams,
			}
		}
	}

	github := &github{
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		RedirectURI:   redirect,
		Orgs:          orgs,
		LoadAllGroups: spec.Config.LoadAllGroups,
		TeamNameField: spec.Config.TeamNameField,
		UseLoginAsID:  spec.Config.UseLoginAsID,
	}

	return github, nil
}
