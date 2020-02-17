package kubestatemetrics

import (
	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/info"
	"github.com/caos/boom/internal/name"
)

type KubeStateMetrics struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.KubeStateMetrics
}

func New(logger logging.Logger) *KubeStateMetrics {
	lo := &KubeStateMetrics{
		logger: logger,
	}

	return lo
}

func (k *KubeStateMetrics) GetName() name.Application {
	return info.GetName()
}

func (k *KubeStateMetrics) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.KubeStateMetrics.Deploy
}

func (k *KubeStateMetrics) GetNamespace() string {
	return info.GetNamespace()
}
