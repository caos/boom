package ambassador

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/ambassador/helm"
	"github.com/caos/boom/internal/templator/helm/chart"

	"github.com/caos/boom/internal/name"
)

const (
	applicationName name.Application = "ambassador"
)

func GetName() name.Application {
	return applicationName
}

type Ambassador struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.Ambassador
}

func New(logger logging.Logger) *Ambassador {
	return &Ambassador{
		logger: logger,
	}
}

func (a *Ambassador) GetName() name.Application {
	return applicationName
}

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.Ambassador.Deploy
}

func (a *Ambassador) Initial() bool {
	return a.spec == nil
}

func (a *Ambassador) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.Ambassador, a.spec)
}

func (a *Ambassador) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	a.spec = toolsetCRDSpec.Ambassador
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
