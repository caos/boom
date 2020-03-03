package prometheus

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type Prometheus struct {
	monitor mntr.Monitor
}

func New(monitor mntr.Monitor) *Prometheus {
	return &Prometheus{
		monitor: monitor,
	}
}

func (p *Prometheus) GetName() name.Application {
	return info.GetName()
}

func (p *Prometheus) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	//not possible to deploy when prometheus operator is not deployed

	po := prometheusoperator.New(p.monitor)
	if !po.Deploy(toolsetCRDSpec) {
		return false
	}

	return toolsetCRDSpec.Prometheus.Deploy
}

func (p *Prometheus) GetNamespace() string {
	return info.GetNamespace()
}
