package kubestatemetrics

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/helm"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/mntr"
)

func (k *KubeStateMetrics) SpecToHelmValues(monitor mntr.Monitor, toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolset.KubeStateMetrics
	values := helm.DefaultValues(k.GetImageTags())

	if spec.ReplicaCount != 0 {
		values.Replicas = spec.ReplicaCount
	}

	return values
}

func (k *KubeStateMetrics) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (k *KubeStateMetrics) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
