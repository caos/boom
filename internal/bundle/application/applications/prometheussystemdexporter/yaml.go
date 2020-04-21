package prometheussystemdexporter

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheussystemdexporter/yaml"
	"github.com/caos/orbiter/mntr"
)

// var _ application.YAMLApplication = (*prometheusSystemdExporter)(nil)

func (*prometheusSystemdExporter) GetYaml(_ mntr.Monitor, _ *v1beta1.ToolsetSpec) interface{} {
	return yaml.GetDefault()
}
