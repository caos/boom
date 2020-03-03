package prometheusoperator

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type PrometheusOperator struct {
	monitor mntr.Monitor
	spec    *toolsetsv1beta1.PrometheusOperator
}

func New(monitor mntr.Monitor) *PrometheusOperator {
	po := &PrometheusOperator{
		monitor: monitor,
	}

	return po
}

func (po *PrometheusOperator) GetName() name.Application {
	return info.GetName()
}

func (po *PrometheusOperator) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.PrometheusOperator.Deploy
}

func (po *PrometheusOperator) GetNamespace() string {
	return info.GetNamespace()
}
