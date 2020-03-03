package config

import (
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type Config struct {
	Monitor           mntr.Monitor
	CrdName           string
	BundleName        name.Bundle
	BaseDirectoryPath string
	Templator         name.Templator
}
