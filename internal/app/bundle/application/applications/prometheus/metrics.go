package prometheus

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/bundle/application/applications/ambassador"
	"github.com/caos/boom/internal/app/bundle/application/applications/apiserver"
	"github.com/caos/boom/internal/app/bundle/application/applications/argocd"
	"github.com/caos/boom/internal/app/bundle/application/applications/kubestatemetrics"
	"github.com/caos/boom/internal/app/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/app/bundle/application/applications/prometheusnodeexporter"
	"github.com/caos/boom/internal/app/bundle/application/applications/prometheusoperator"
)

func ScrapeMetricsCrdsConfig(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) *Config {

	monitorlabels := make(map[string]string, 0)
	monitorlabels["app.kubernetes.io/managed-by"] = "boom.caos.ch"

	servicemonitors := make([]*servicemonitor.Config, 0)

	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy &&
		(toolsetCRDSpec.Metrics == nil || toolsetCRDSpec.Metrics.Ambassador) {
		servicemonitors = append(servicemonitors, ambassador.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.PrometheusOperator != nil && toolsetCRDSpec.PrometheusOperator.Deploy &&
		(toolsetCRDSpec.Metrics == nil || toolsetCRDSpec.Metrics.PrometheusOperator) {
		servicemonitors = append(servicemonitors, prometheusoperator.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.PrometheusNodeExporter != nil && toolsetCRDSpec.PrometheusNodeExporter.Deploy &&
		(toolsetCRDSpec.Metrics == nil || toolsetCRDSpec.Metrics.PrometheusNodeExporter) {
		servicemonitors = append(servicemonitors, prometheusnodeexporter.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.KubeStateMetrics != nil && toolsetCRDSpec.KubeStateMetrics.Deploy &&
		(toolsetCRDSpec.Metrics == nil || toolsetCRDSpec.Metrics.KubeStateMetrics) {
		servicemonitors = append(servicemonitors, kubestatemetrics.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.Argocd != nil && toolsetCRDSpec.Argocd.Deploy &&
		(toolsetCRDSpec.Metrics == nil || toolsetCRDSpec.Metrics.Argocd) {
		servicemonitors = append(servicemonitors, argocd.GetServicemonitors(monitorlabels)...)
	}

	if toolsetCRDSpec.Metrics == nil || toolsetCRDSpec.Metrics.APIServer {
		servicemonitors = append(servicemonitors, apiserver.GetServicemonitor(monitorlabels))
	}

	if len(servicemonitors) > 0 {
		servicemonitors = append(servicemonitors, GetServicemonitor(monitorlabels))

		prom := &Config{
			Prefix:                  "",
			Namespace:               "caos-system",
			MonitorLabels:           monitorlabels,
			ServiceMonitors:         servicemonitors,
			AdditionalScrapeConfigs: getScrapeConfigs(),
			KubeVersion:             toolsetCRDSpec.KubeVersion,
		}

		return prom
	}
	return nil
}
