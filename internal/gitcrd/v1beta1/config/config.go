package config

import (
	"github.com/caos/boom/internal/git"
	"github.com/caos/orbiter/mntr"
)

type Config struct {
	Monitor          mntr.Monitor
	Git              *git.Client
	CrdDirectoryPath string
	CrdPath          string
}
