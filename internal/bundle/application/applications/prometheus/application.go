package prometheus

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "prometheus"
	instanceName    string           = "caos"
	orderNumber     int              = 1
)

func GetName() name.Application {
	return applicationName
}

func GetOrderNumber() int {
	return orderNumber
}
