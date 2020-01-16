package ambassador

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/ambassador/helm"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
)

func (a *Ambassador) GetName() name.Application {
	return GetName()
}

func (a *Ambassador) GetNamespace() string {
	return "caos-system"
}

func (a *Ambassador) SpecToHelmValues(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) interface{} {
	a.spec = toolsetCRDSpec.Ambassador
	imageTags := helm.GetImageTags()

	values := helm.DefaultValues(imageTags)
	if a.spec.ReplicaCount != 0 {
		values.ReplicaCount = a.spec.ReplicaCount
	}

	return values
}

func (a *Ambassador) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (a *Ambassador) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
