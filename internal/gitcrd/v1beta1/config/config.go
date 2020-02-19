package config

import (
	"github.com/caos/boom/internal/git"
	"github.com/caos/orbiter/logging"
)

type Config struct {
	Logger           logging.Logger
	Git              *git.Client
	CrdDirectoryPath string
	CrdPath          string
}
