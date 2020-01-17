package prometheusnodeexporter

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
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
	return applicationName
}

func (pne *PrometheusNodeExporter) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.PrometheusNodeExporter.Deploy
}

func (pne *PrometheusNodeExporter) Initial() bool {
	return pne.spec == nil
}

func (pne *PrometheusNodeExporter) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.PrometheusNodeExporter, pne.spec)
}

func (pne *PrometheusNodeExporter) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	pne.spec = toolsetCRDSpec.PrometheusNodeExporter
}

func (pne *PrometheusNodeExporter) GetNamespace() string {
	return "caos-system"
}
