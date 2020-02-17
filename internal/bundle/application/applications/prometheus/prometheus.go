package prometheus

import (
	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator"
	"github.com/caos/boom/internal/name"
)

type Prometheus struct {
	logger logging.Logger
}

func New(logger logging.Logger) *Prometheus {
	return &Prometheus{
		logger: logger,
	}
}

func (p *Prometheus) GetName() name.Application {
	return info.GetName()
}

func (p *Prometheus) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	//not possible to deploy when prometheus operator is not deployed

	po := prometheusoperator.New(p.logger)
	if !po.Deploy(toolsetCRDSpec) {
		return false
	}

	return toolsetCRDSpec.Prometheus.Deploy
}

func (p *Prometheus) GetNamespace() string {
	return info.GetNamespace()
}
