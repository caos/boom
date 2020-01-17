package loggingoperator

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "logging-operator"
)

func GetName() name.Application {
	return applicationName
}
