package kubestatemetrics

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/helm"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
)

const (
	applicationName name.Application = "kube-state-metrics"
)

func GetName() name.Application {
	return applicationName
}

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

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
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

func (k *KubeStateMetrics) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.CertManager
	values := helm.DefaultValues(k.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	return values
}

func (k *KubeStateMetrics) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (k *KubeStateMetrics) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
