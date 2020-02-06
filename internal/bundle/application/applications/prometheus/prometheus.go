package prometheus

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/config"
	"github.com/caos/boom/internal/name"
)

type Prometheus struct {
	logger logging.Logger
	config *config.Config
}

func New(logger logging.Logger) *Prometheus {
	return &Prometheus{
		logger: logger,
	}
}

func (p *Prometheus) GetName() name.Application {
	return applicationName
}

func (p *Prometheus) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	//not possible to deploy when prometheus operator is not deployed
	if !toolsetCRDSpec.PrometheusOperator.Deploy {
		return false
	}

	return toolsetCRDSpec.Prometheus.Deploy
}

func (p *Prometheus) Initial() bool {
	return p.config == nil
}

func (p *Prometheus) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	config := config.ScrapeMetricsCrdsConfig(toolsetCRDSpec)
	return !reflect.DeepEqual(config, p.config)
}

func (p *Prometheus) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	if toolsetCRDSpec == nil {
		return
	}
	p.config = config.ScrapeMetricsCrdsConfig(toolsetCRDSpec)
}

func (p *Prometheus) GetNamespace() string {
	return "caos-system"
}
