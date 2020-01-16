package config

import (
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/logging"
)

type Config struct {
	Logger  logging.Logger
	Version name.Version
}
