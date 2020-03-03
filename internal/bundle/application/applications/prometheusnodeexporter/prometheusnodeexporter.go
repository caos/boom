package prometheusnodeexporter

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type PrometheusNodeExporter struct {
	monitor mntr.Monitor
	spec    *toolsetsv1beta1.PrometheusNodeExporter
}

func New(monitor mntr.Monitor) *PrometheusNodeExporter {
	pne := &PrometheusNodeExporter{
		monitor: monitor,
	}

	return pne
}

func (pne *PrometheusNodeExporter) GetName() name.Application {
	return info.GetName()
}

func (pne *PrometheusNodeExporter) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.PrometheusNodeExporter.Deploy
}

func (pne *PrometheusNodeExporter) GetNamespace() string {
	return info.GetNamespace()
}
