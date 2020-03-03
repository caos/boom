package auth

import (
	"io/ioutil"
	"os"
	"path/filepath"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/clientgo"
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
	secret, err := clientgo.GetSecret(spec.Config.SecretName, "caos-system")
	if err != nil {
		return "", err
	}

	clientID := string(secret.Data[spec.Config.ClientIDKey])
	clientSecret := string(secret.Data[spec.Config.ClientSecretKey])
	serviceAccountJSON := secret.Data[spec.Config.ServiceAccountJSONKey]

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

	if serviceAccountJSON != nil {
		// write json to file
		err = ioutil.WriteFile(spec.Config.ServiceAccountFilePath, serviceAccountJSON, 0644)
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
