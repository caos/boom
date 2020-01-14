package prometheusoperator

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/helm"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
)

const (
	applicationName name.Application = "prometheus-operator"
)

func GetName() name.Application {
	return applicationName
}

type PrometheusOperator struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.PrometheusOperator
}

func New(logger logging.Logger) *PrometheusOperator {
	po := &PrometheusOperator{
		logger: logger,
	}

	return po
}

func (po *PrometheusOperator) GetName() name.Application {
	return applicationName
}

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.PrometheusOperator.Deploy
}

func (po *PrometheusOperator) Initial() bool {
	return po.spec == nil
}

func (po *PrometheusOperator) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.PrometheusOperator, po.spec)
}

func (po *PrometheusOperator) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	po.spec = toolsetCRDSpec.PrometheusOperator
}

func (po *PrometheusOperator) GetNamespace() string {
	return "caos-system"
}

func (p *PrometheusOperator) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.PrometheusNodeExporter
	values := helm.DefaultValues(p.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	return values
}

func (p *PrometheusOperator) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()

}

func (p *PrometheusOperator) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
