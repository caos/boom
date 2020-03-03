package loki

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/loki/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type Loki struct {
	monitor mntr.Monitor
	spec    *toolsetsv1beta1.Loki
}

func New(monitor mntr.Monitor) *Loki {
	lo := &Loki{
		monitor: monitor,
	}
	return lo
}

func (l *Loki) GetName() name.Application {
	return info.GetName()
}

func (lo *Loki) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.Loki.Deploy
}

func (l *Loki) GetNamespace() string {
	return info.GetNamespace()
}
