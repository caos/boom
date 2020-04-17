package auth

import (
	"github.com/caos/boom/internal/helper"
	"io/ioutil"
	"os"
	"path/filepath"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/pkg/errors"
)

type google struct {
	Issuer                 string   `yaml:"issuer,omitempty"`
	ClientID               string   `yaml:"clientID,omitempty"`
	ClientSecret           string   `yaml:"clientSecret,omitempty"`
	RedirectURI            string   `yaml:"redirectURI,omitempty"`
	HostedDomains          []string `yaml:"hostedDomains,omitempty"`
	Groups                 []string `yaml:"groups,omitempty"`
	ServiceAccountFilePath string   `yaml:"serviceAccountFilePath,omitempty"`
	AdminEmail             string   `yaml:"adminEmail,omitempty"`
}

func getGoogle(spec *toolsetsv1beta1.ArgocdGoogleConnector, redirect string) (interface{}, error) {
	clientID, err := helper.GetSecretValue(spec.Config.ClientID, spec.Config.ExistingClientIDSecret)
	if err != nil {
		return nil, err
	}

	clientSecret, err := helper.GetSecretValue(spec.Config.ClientSecret, spec.Config.ExistingClientSecretSecret)
	if err != nil {
		return nil, err
	}

	serviceAccountJSON, err := helper.GetSecretValue(spec.Config.ServiceAccountJSON, spec.Config.ExistingServiceAccountJSONSecret)
	if err != nil {
		return nil, err
	}

	// get base path
	base, err := filepath.Abs(spec.Config.ServiceAccountFilePath)
	if err != nil {
		return nil, err
	}

	// remove file if alread exists
	_, err = os.Stat(spec.Config.ServiceAccountFilePath)
	if !os.IsNotExist(err) {
		if err := os.Remove(spec.Config.ServiceAccountFilePath); err != nil {
			return nil, err
		}
	}

	// create all directories to the file
	if err := os.MkdirAll(base, os.ModePerm); err != nil {
		return nil, err
	}

	if serviceAccountJSON != "" {
		// write json to file
		err = ioutil.WriteFile(spec.Config.ServiceAccountFilePath, []byte(serviceAccountJSON), 0644)
		if err != nil {
			return nil, errors.Wrapf(err, "Error while writing json to file %s", spec.Config.ServiceAccountFilePath)
		}
	}

	google := &google{
		ClientID:               clientID,
		ClientSecret:           clientSecret,
		RedirectURI:            redirect,
		Groups:                 spec.Config.Groups,
		HostedDomains:          spec.Config.HostedDomains,
		ServiceAccountFilePath: spec.Config.ServiceAccountFilePath,
		AdminEmail:             spec.Config.AdminEmail,
		Issuer:                 "https://accounts.google.com",
	}

	return google, nil
}
