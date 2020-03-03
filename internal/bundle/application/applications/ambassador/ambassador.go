package ambassador

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/ambassador/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type Ambassador struct {
	monitor mntr.Monitor
	spec    *toolsetsv1beta1.Ambassador
}

func New(monitor mntr.Monitor) *Ambassador {
	return &Ambassador{
		monitor: monitor,
	}
}

func (a *Ambassador) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.Ambassador.Deploy
}

func (a *Ambassador) GetName() name.Application {
	return info.GetName()
}

func (a *Ambassador) GetNamespace() string {
	return info.GetNamespace()
}
