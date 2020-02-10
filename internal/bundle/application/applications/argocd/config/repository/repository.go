package repository

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/orbiter/logging"
)

type Repository struct {
	URL                 string
	UsernameSecret      *secret `yaml:"usernameSecret,omitempty"`
	PasswordSecret      *secret `yaml:"passwordSecret,omitempty"`
	SSHPrivateKeySecret *secret `yaml:"sshPrivateKeySecret,omitempty"`
}

type secret struct {
	Name string
	Key  string
}

func GetFromSpec(logger logging.Logger, spec *toolsetsv1beta1.Argocd) []*Repository {
	repositories := make([]*Repository, 0)

	if spec.Repositories == nil || len(spec.Repositories) == 0 {
		return repositories
	}

	for _, v := range spec.Repositories {
		var us, ps, ssh *secret
		if v.UsernameSecret != nil {
			us = &secret{
				Name: v.UsernameSecret.Name,
				Key:  v.UsernameSecret.Key,
			}
		}
		if v.PasswordSecret != nil {
			ps = &secret{
				Name: v.PasswordSecret.Name,
				Key:  v.PasswordSecret.Key,
			}
		}
		if v.CertificateSecret != nil {
			ssh = &secret{
				Name: v.CertificateSecret.Name,
				Key:  v.CertificateSecret.Key,
			}
		}

		repo := &Repository{
			URL:                 v.URL,
			UsernameSecret:      us,
			PasswordSecret:      ps,
			SSHPrivateKeySecret: ssh,
		}
		repositories = append(repositories, repo)
	}

	return repositories
}
