package kubestatemetrics

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type KubeStateMetrics struct {
	monitor mntr.Monitor
	spec    *toolsetsv1beta1.KubeStateMetrics
}

func New(monitor mntr.Monitor) *KubeStateMetrics {
	lo := &KubeStateMetrics{
		monitor: monitor,
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
