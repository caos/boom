package prometheusoperator

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "prometheus-operator"
)

func GetName() name.Application {
	return applicationName
}
