package info

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "preapply"
	orderNumber     int              = 0
)

func GetName() name.Application {
	return applicationName
}

func GetOrderNumber() int {
	return orderNumber
}
