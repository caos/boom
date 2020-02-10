package prometheusoperator

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/helm"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/logging"
)

func (p *PrometheusOperator) SpecToHelmValues(logger logging.Logger, toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.PrometheusNodeExporter
	values := helm.DefaultValues(p.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	return values
}

func (p *PrometheusOperator) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (p *PrometheusOperator) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
