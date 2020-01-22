package prometheusnodeexporter

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/helm"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/info"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/boom/internal/templator/helm/chart"
)

func (p *PrometheusNodeExporter) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.PrometheusNodeExporter
	values := helm.DefaultValues(p.GetImageTags())

	values.PodLabels = labels.GetApplicationLabels(p.GetName())
	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	values.PodLabels = labels.GetApplicationLabels(info.GetName())
	return values
}

func (p *PrometheusNodeExporter) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (c *PrometheusNodeExporter) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
