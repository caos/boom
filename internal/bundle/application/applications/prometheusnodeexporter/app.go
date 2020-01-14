package prometheusnodeexporter

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/helm"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
)

const (
	applicationName name.Application = "prometheus-node-exporter"
)

func GetName() name.Application {
	return applicationName
}

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

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
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

func (p *PrometheusNodeExporter) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.PrometheusNodeExporter
	values := helm.DefaultValues(p.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	return values
}

func (p *PrometheusNodeExporter) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (c *PrometheusNodeExporter) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
