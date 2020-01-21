package loki

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "loki"
	orderNumber     int              = 1
)

func GetName() name.Application {
	return applicationName
}
func GetOrderNumber() int {
	return orderNumber
}
