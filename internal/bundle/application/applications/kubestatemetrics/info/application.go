package info

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "kube-state-metrics"
	namespace       string           = "caos-system"
)

func GetName() name.Application {
	return applicationName
}

func GetNamespace() string {
	return namespace
}
