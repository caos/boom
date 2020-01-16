package config

import (
	"github.com/caos/orbiter/logging"
)

type Config struct {
	Logger           logging.Logger
	CrdUrl           string
	CrdDirectoryPath string
	CrdPath          string
	PrivateKey       []byte
}
