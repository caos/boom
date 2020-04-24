package credential

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/orbiter/mntr"
)

type Credential struct {
	URL                 string
	UsernameSecret      *secret `yaml:"usernameSecret,omitempty"`
	PasswordSecret      *secret `yaml:"passwordSecret,omitempty"`
	SSHPrivateKeySecret *secret `yaml:"sshPrivateKeySecret,omitempty"`
}

type secret struct {
	Name string
	Key  string
}

func GetFromSpec(monitor mntr.Monitor, spec *toolsetsv1beta1.Argocd) []*Credential {
	creds := make([]*Credential, 0)

	if spec.Credentials == nil || len(spec.Credentials) == 0 {
		return creds
	}

	for _, v := range spec.Credentials {
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

		cred := &Credential{
			URL:                 v.URL,
			UsernameSecret:      us,
			PasswordSecret:      ps,
			SSHPrivateKeySecret: ssh,
		}
		creds = append(creds, cred)
	}

	return creds
}
