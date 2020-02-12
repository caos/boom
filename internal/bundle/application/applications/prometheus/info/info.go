package info

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "prometheus"
	instanceName    string           = "caos"
	orderNumber     int              = 1
	namespace       string           = "caos-system"
)

func GetName() name.Application {
	return applicationName
}

func GetNamespace() string {
	return namespace
}
func GetOrderNumber() int {
	return orderNumber
}

func GetInstanceName() string {
	return instanceName
}
