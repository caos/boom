package prometheusnodeexporter

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "prometheus-node-exporter"
)

func GetName() name.Application {
	return applicationName
}
