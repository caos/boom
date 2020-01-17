package kubestatemetrics

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "kube-state-metrics"
)

func GetName() name.Application {
	return applicationName
}
