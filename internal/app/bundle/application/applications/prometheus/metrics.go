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

func ScrapeMetricsCrdsConfig(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) (*Config, error) {

	monitorlabels := make(map[string]string, 0)
	monitorlabels["app.kubernetes.io/managed-by"] = "boom.caos.ch"

	servicemonitors := make([]*servicemonitor.Config, 0)

	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy {
		servicemonitors = append(servicemonitors, ambassador.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.PrometheusOperator != nil && toolsetCRDSpec.PrometheusOperator.Deploy {
		servicemonitors = append(servicemonitors, prometheusoperator.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.PrometheusNodeExporter != nil && toolsetCRDSpec.PrometheusNodeExporter.Deploy {
		servicemonitors = append(servicemonitors, prometheusnodeexporter.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.KubeStateMetrics != nil && toolsetCRDSpec.KubeStateMetrics.Deploy {
		servicemonitors = append(servicemonitors, kubestatemetrics.GetServicemonitor(monitorlabels))
	}

	if toolsetCRDSpec.Argocd != nil && toolsetCRDSpec.Argocd.Deploy {
		servicemonitors = append(servicemonitors, argocd.GetServicemonitors(monitorlabels)...)
	}

	servicemonitors = append(servicemonitors, apiserver.GetServicemonitor(monitorlabels))
	servicemonitors = append(servicemonitors, GetServicemonitor(monitorlabels))

	prom := &Config{
		Prefix:                  "",
		Namespace:               "caos-system",
		MonitorLabels:           monitorlabels,
		ServiceMonitors:         servicemonitors,
		AdditionalScrapeConfigs: getScrapeConfigs(),
		KubeVersion:             toolsetCRDSpec.KubeVersion,
	}

	return prom, nil
}
