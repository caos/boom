package prometheus

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	"github.com/caos/boom/api/v1beta1"
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/config"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/helm"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
)

const (
	applicationName name.Application = "prometheus"
)

func GetName() name.Application {
	return applicationName
}

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
	config := config.ScrapeMetricsCrdsConfig(toolsetCRDSpec)
	if config == nil {
		return false
	}
	return true
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

func (p *Prometheus) SpecToHelmValues(toolsetCRDSpec *v1beta1.ToolsetSpec) interface{} {
	config := config.ScrapeMetricsCrdsConfig(toolsetCRDSpec)

	values := helm.DefaultValues(p.GetImageTags())

	if config.StorageSpec != nil {
		storageSpec := &helm.StorageSpec{
			VolumeClaimTemplate: &helm.VolumeClaimTemplate{
				Spec: &helm.VolumeClaimTemplateSpec{
					StorageClassName: config.StorageSpec.StorageClass,
					AccessModes:      config.StorageSpec.AccessModes,
					Resources: &helm.Resources{
						Requests: &helm.Request{
							Storage: config.StorageSpec.Storage,
						},
					},
				},
			},
		}

		values.Prometheus.PrometheusSpec.StorageSpec = storageSpec
	}

	if config.MonitorLabels != nil {
		values.Prometheus.PrometheusSpec.ServiceMonitorSelector = &helm.MonitorSelector{
			MatchLabels: config.MonitorLabels,
		}
	}

	if config.ServiceMonitors != nil {
		additionalServiceMonitors := make([]*servicemonitor.Values, 0)
		for _, specServiceMonitor := range config.ServiceMonitors {
			valuesServiceMonitor := servicemonitor.SpecToValues(specServiceMonitor)
			additionalServiceMonitors = append(additionalServiceMonitors, valuesServiceMonitor)
		}

		values.Prometheus.AdditionalServiceMonitors = additionalServiceMonitors
	}

	if config.AdditionalScrapeConfigs != nil {
		values.Prometheus.PrometheusSpec.AdditionalScrapeConfigs = config.AdditionalScrapeConfigs
	}

	ruleLabels := map[string]string{"prometheus": "caos", "instance-name": "operated"}
	rules, _ := helm.GetDefaultRules(ruleLabels)

	values.Prometheus.PrometheusSpec.RuleSelector = &helm.RuleSelector{MatchLabels: ruleLabels}
	values.DefaultRules.Labels = ruleLabels
	values.KubeTargetVersionOverride = config.KubeVersion
	values.AdditionalPrometheusRules = []*helm.AdditionalPrometheusRules{rules}

	return values
}

func (p *Prometheus) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (p *Prometheus) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
