package config

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	ambassadormetrics "github.com/caos/boom/internal/bundle/application/applications/ambassador/metrics"
	"github.com/caos/boom/internal/bundle/application/applications/apiserver"
	argocdmetrics "github.com/caos/boom/internal/bundle/application/applications/argocd/metrics"
	kubestatemetrics "github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/metrics"
	lometrics "github.com/caos/boom/internal/bundle/application/applications/loggingoperator/metrics"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/metrics"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	pnemetrics "github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/metrics"
	pometrics "github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/metrics"
)

func ScrapeMetricsCrdsConfig(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) *Config {

	monitorlabels := make(map[string]string, 0)
	monitorlabels["app.kubernetes.io/managed-by"] = "boom.caos.ch"

	servicemonitors := make([]*servicemonitor.Config, 0)

	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.Ambassador) {
		servicemonitors = append(servicemonitors, ambassadormetrics.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.PrometheusOperator != nil && toolsetCRDSpec.PrometheusOperator.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.PrometheusOperator) {
		servicemonitors = append(servicemonitors, pometrics.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.PrometheusNodeExporter != nil && toolsetCRDSpec.PrometheusNodeExporter.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.PrometheusNodeExporter) {
		servicemonitors = append(servicemonitors, pnemetrics.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.KubeStateMetrics != nil && toolsetCRDSpec.KubeStateMetrics.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.KubeStateMetrics) {
		servicemonitors = append(servicemonitors, kubestatemetrics.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.Argocd != nil && toolsetCRDSpec.Argocd.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.Argocd) {
		servicemonitors = append(servicemonitors, argocdmetrics.GetServicemonitors(monitorlabels)...)
	}

	if toolsetCRDSpec.LoggingOperator != nil && toolsetCRDSpec.LoggingOperator.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.LoggingOperator) {
		servicemonitors = append(servicemonitors, lometrics.GetServicemonitors(monitorlabels)...)
	}

	if toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.APIServer {
		servicemonitors = append(servicemonitors, apiserver.GetServicemonitor(monitorlabels))
	}

	if len(servicemonitors) > 0 {
		servicemonitors = append(servicemonitors, metrics.GetServicemonitor(monitorlabels))

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
