package loki

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/loki/helm"
	"github.com/caos/boom/internal/bundle/application/applications/loki/info"
	"github.com/caos/boom/internal/bundle/application/applications/loki/logs"
	"github.com/caos/orbiter/logging"

	"github.com/caos/boom/internal/templator/helm/chart"
)

func (l *Loki) HelmPreApplySteps(logger logging.Logger, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) ([]interface{}, error) {
	return logs.GetAllResources(toolsetCRDSpec), nil
}

func (l *Loki) SpecToHelmValues(logger logging.Logger, toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolset.Loki
	values := helm.DefaultValues(l.GetImageTags())

	if spec.Storage != nil {
		values.Persistence.Enabled = true
		values.Persistence.Size = spec.Storage.Size
		values.Persistence.StorageClassName = spec.Storage.StorageClass
		if spec.Storage.AccessModes != nil {
			values.Persistence.AccessModes = spec.Storage.AccessModes
		}
	}

	values.FullNameOverride = info.GetName().String()
	return values
}

func (l *Loki) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (l *Loki) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
