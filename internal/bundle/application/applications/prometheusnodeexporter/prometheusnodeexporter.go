package prometheusnodeexporter

import (
	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/info"
	"github.com/caos/boom/internal/name"
)

type PrometheusNodeExporter struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.PrometheusNodeExporter
}

func New(logger logging.Logger) *PrometheusNodeExporter {
	pne := &PrometheusNodeExporter{
		logger: logger,
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
