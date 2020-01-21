package loki

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/loki/helm"
	"github.com/caos/boom/internal/bundle/application/applications/loki/logs"

	"github.com/caos/boom/internal/templator/helm/chart"
)

func (l *Loki) HelmPreApplySteps(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) ([]interface{}, error) {
	return logs.GetAllResources(toolsetCRDSpec), nil
}

func (l *Loki) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.LoggingOperator
	values := helm.DefaultValues(l.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	return values
}

func (l *Loki) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (l *Loki) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
