package config

import (
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type Config struct {
	Monitor mntr.Monitor
	Version name.Version
}
