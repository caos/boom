package argocd

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/helm"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
)

const (
	applicationName name.Application = "argocd"
)

func GetName() name.Application {
	return applicationName
}

type Argocd struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.Argocd
}

func New(logger logging.Logger) *Argocd {
	c := &Argocd{
		logger: logger,
	}

	return c
}

func (a *Argocd) GetName() name.Application {
	return applicationName
}

func (a *Argocd) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.Argocd.Deploy
}

func (a *Argocd) Initial() bool {
	return a.spec == nil
}

func (a *Argocd) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.Argocd, a.spec)
}

func (a *Argocd) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	a.spec = toolsetCRDSpec.Argocd
}

func (a *Argocd) GetNamespace() string {
	return "caos-system"
}

func (a *Argocd) SpecToHelmValues(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.Argocd
	values := helm.DefaultValues(a.GetImageTags())

	return values
}

func (a *Argocd) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (a *Argocd) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
