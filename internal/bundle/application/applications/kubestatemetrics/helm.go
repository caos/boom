package kubestatemetrics

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/helm"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/info"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/boom/internal/templator/helm/chart"
)

func (k *KubeStateMetrics) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolset.KubeStateMetrics
	values := helm.DefaultValues(k.GetImageTags())

	if spec.ReplicaCount != 0 {
		values.Replicas = spec.ReplicaCount
	}

	values.CustomLabels = labels.GetApplicationLabels(info.GetName())
	return values
}

func (k *KubeStateMetrics) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (k *KubeStateMetrics) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
