package kubestatemetrics

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
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
	return applicationName
}

func (k *KubeStateMetrics) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.KubeStateMetrics.Deploy
}

func (k *KubeStateMetrics) Initial() bool {
	return k.spec == nil
}

func (k *KubeStateMetrics) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.KubeStateMetrics, k.spec)
}

func (k *KubeStateMetrics) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	k.spec = toolsetCRDSpec.KubeStateMetrics
}

func (k *KubeStateMetrics) GetNamespace() string {
	return "caos-system"
}
