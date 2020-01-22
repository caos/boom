package argocd

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/info"
	"github.com/caos/boom/internal/name"
)

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
	return info.GetName()
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
	return info.GetNamespace()
}
