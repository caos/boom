package info

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "loki"
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
