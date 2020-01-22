package config

import (
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/logging"
)

type Config struct {
	Logger            logging.Logger
	CrdName           string
	BundleName        name.Bundle
	BaseDirectoryPath string
	Templator         name.Templator
}
