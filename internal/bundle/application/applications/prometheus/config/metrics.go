package config

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	ambassadormetrics "github.com/caos/boom/internal/bundle/application/applications/ambassador/metrics"
	"github.com/caos/boom/internal/bundle/application/applications/apiserver"
	argocdmetrics "github.com/caos/boom/internal/bundle/application/applications/argocd/metrics"
	"github.com/caos/boom/internal/bundle/application/applications/boom"
	kubestatemetrics "github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/metrics"
	lometrics "github.com/caos/boom/internal/bundle/application/applications/loggingoperator/metrics"
	lokimetrics "github.com/caos/boom/internal/bundle/application/applications/loki/metrics"
	"github.com/caos/boom/internal/bundle/application/applications/orbiter"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/metrics"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	pnemetrics "github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/metrics"
	pometrics "github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/metrics"
	psemetrics "github.com/caos/boom/internal/bundle/application/applications/prometheussystemdexporter/metrics"
	"github.com/caos/boom/internal/labels"
)

func ScrapeMetricsCrdsConfig(instanceName string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) *Config {
	servicemonitors := make([]*servicemonitor.Config, 0)

	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.Ambassador) {
		servicemonitors = append(servicemonitors, ambassadormetrics.GetServicemonitor(instanceName))
	}

	if toolsetCRDSpec.PrometheusOperator != nil && toolsetCRDSpec.PrometheusOperator.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.PrometheusOperator) {
		servicemonitors = append(servicemonitors, pometrics.GetServicemonitor(instanceName))
	}

	if toolsetCRDSpec.PrometheusNodeExporter != nil && toolsetCRDSpec.PrometheusNodeExporter.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.PrometheusNodeExporter) {
		servicemonitors = append(servicemonitors, pnemetrics.GetServicemonitors(instanceName)...)
	}

	if toolsetCRDSpec.PrometheusSystemdExporter != nil && toolsetCRDSpec.PrometheusSystemdExporter.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.PrometheusSystemdExporter) {
		servicemonitors = append(servicemonitors, psemetrics.GetServicemonitor(instanceName))
	}

	if toolsetCRDSpec.KubeStateMetrics != nil && toolsetCRDSpec.KubeStateMetrics.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.KubeStateMetrics) {
		servicemonitors = append(servicemonitors, kubestatemetrics.GetServicemonitors(instanceName)...)
	}

	if toolsetCRDSpec.Argocd != nil && toolsetCRDSpec.Argocd.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.Argocd) {
		servicemonitors = append(servicemonitors, argocdmetrics.GetServicemonitors(instanceName)...)
	}

	if toolsetCRDSpec.LoggingOperator != nil && toolsetCRDSpec.LoggingOperator.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.LoggingOperator) {
		servicemonitors = append(servicemonitors, lometrics.GetServicemonitors(instanceName)...)
	}

	if toolsetCRDSpec.Loki != nil && toolsetCRDSpec.Loki.Deploy &&
		(toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.Loki) {
		servicemonitors = append(servicemonitors, lokimetrics.GetServicemonitor(instanceName))
	}

	if toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.APIServer {
		servicemonitors = append(servicemonitors, apiserver.GetServicemonitor(instanceName))
	}

	if toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.Boom {
		servicemonitors = append(servicemonitors, boom.GetServicemonitor(instanceName))
	}

	if toolsetCRDSpec.Prometheus.Metrics == nil || toolsetCRDSpec.Prometheus.Metrics.Orbiter {
		servicemonitors = append(servicemonitors, orbiter.GetServicemonitor(instanceName))
	}


	if len(servicemonitors) > 0 {

		servicemonitors = append(servicemonitors, metrics.GetServicemonitor(instanceName))

		prom := &Config{
			Prefix:                  "",
			Namespace:               "caos-system",
			MonitorLabels:           labels.GetMonitorSelectorLabels(instanceName),
			ServiceMonitors:         servicemonitors,
			AdditionalScrapeConfigs: getScrapeConfigs(),
		}

		if toolsetCRDSpec.Prometheus.Storage != nil {
			prom.StorageSpec = &StorageSpec{
				StorageClass: toolsetCRDSpec.Prometheus.Storage.StorageClass,
				Storage:      toolsetCRDSpec.Prometheus.Storage.Size,
			}

			if toolsetCRDSpec.Prometheus.Storage.AccessModes != nil {
				prom.StorageSpec.AccessModes = toolsetCRDSpec.Prometheus.Storage.AccessModes
			}
		}

		return prom
	}
	return nil
}
