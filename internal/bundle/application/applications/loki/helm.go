package loki

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/loki/helm"
	"github.com/caos/boom/internal/bundle/application/applications/loki/info"
	"github.com/caos/boom/internal/bundle/application/applications/loki/logs"
	"github.com/caos/boom/internal/labels"

	"github.com/caos/boom/internal/templator/helm/chart"
)

func (l *Loki) HelmPreApplySteps(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) ([]interface{}, error) {
	return logs.GetAllResources(toolsetCRDSpec), nil
}

func (l *Loki) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolset.Loki
	values := helm.DefaultValues(l.GetImageTags())

	if spec.Storage != nil {
		values.Persistence.Enabled = true
		values.Persistence.AccessModes = spec.Storage.AccessModes
		values.Persistence.Size = spec.Storage.Size
		values.Persistence.StorageClassName = spec.Storage.StorageClass
	}

	appLabels := labels.GetApplicationLabels(info.GetName())
	values.PodLabels = appLabels
	values.Service.Labels = appLabels
	values.FullNameOverride = info.GetName().String()
	return values
}

func (l *Loki) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (l *Loki) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
