package info

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "postapply"
	orderNumber     int              = 100
)

func GetName() name.Application {
	return applicationName
}

func GetOrderNumber() int {
	return orderNumber
}
