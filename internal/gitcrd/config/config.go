package config

import "github.com/caos/orbiter/mntr"

type Config struct {
	Monitor          mntr.Monitor
	CrdUrl           string
	CrdDirectoryPath string
	CrdPath          string
	PrivateKey       []byte
	User             string
	Email            string
}
