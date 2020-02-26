package argocd

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type Argocd struct {
	monitor mntr.Monitor
	spec    *toolsetsv1beta1.Argocd
}

func New(monitor mntr.Monitor) *Argocd {
	c := &Argocd{
		monitor: monitor,
	}

	return c
}

func (a *Argocd) GetName() name.Application {
	return info.GetName()
}

func (a *Argocd) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.Argocd.Deploy
}

func (a *Argocd) GetNamespace() string {
	return info.GetNamespace()
}
