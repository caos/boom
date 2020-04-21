package info

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "prometheus-systemd-exporter"
	namespace       string           = "caos-system"
)

func GetName() name.Application {
	return applicationName
}

func GetNamespace() string {
	return namespace
}
