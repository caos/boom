package prometheusoperator

import (
	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/info"
	"github.com/caos/boom/internal/name"
)

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
	return info.GetName()
}

func (po *PrometheusOperator) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.PrometheusOperator.Deploy
}

func (po *PrometheusOperator) GetNamespace() string {
	return info.GetNamespace()
}
