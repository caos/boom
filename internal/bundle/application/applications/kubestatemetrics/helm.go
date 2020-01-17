package kubestatemetrics

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/helm"
	"github.com/caos/boom/internal/templator/helm/chart"
)

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
