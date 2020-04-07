package prometheussystemdexporter

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheussystemdexporter/info"
	"github.com/caos/boom/internal/name"
)

type prometheusSystemdExporter struct{}

func New() *prometheusSystemdExporter {
	return &prometheusSystemdExporter{}
}

func (*prometheusSystemdExporter) GetName() name.Application {
	return info.GetName()
}

func (*prometheusSystemdExporter) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.PrometheusSystemdExporter.Deploy
}

func (*prometheusSystemdExporter) GetNamespace() string {
	return info.GetNamespace()
}
