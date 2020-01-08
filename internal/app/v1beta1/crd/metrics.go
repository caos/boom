package crd

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/v1beta1/crd/ambassador"
	"github.com/caos/boom/internal/app/v1beta1/crd/apiserver"
	"github.com/caos/boom/internal/app/v1beta1/crd/argocd"
	"github.com/caos/boom/internal/app/v1beta1/crd/kubelet"
	"github.com/caos/boom/internal/app/v1beta1/crd/kubestatemetrics"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheus"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheus/servicemonitor"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheusnodeexporter"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheusoperator"
)

func (c *Crd) ScrapeMetricsCrdsConfig(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) (*prometheus.Config, string, error) {

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
	servicemonitors = append(servicemonitors, prometheus.GetServicemonitor(monitorlabels))

	adconfigs := make([]*prometheus.AdditionalScrapeConfig, 0)
	adconfigs = append(adconfigs, kubelet.GetScrapeConfigs()...)

	prom := &prometheus.Config{
		Prefix:                  "",
		Namespace:               "caos-system",
		MonitorLabels:           monitorlabels,
		ServiceMonitors:         servicemonitors,
		AdditionalScrapeConfigs: adconfigs,
		KubeVersion:             toolsetCRDSpec.KubeVersion,
	}

	datasource := ""
	if prom.Prefix != "" {
		datasource = strings.Join([]string{"http://", prom.Prefix, "-prometheus-operated.", prom.Namespace, ":9090"}, "")
	} else {
		datasource = strings.Join([]string{"http://prometheus-operated.", prom.Namespace, ":9090"}, "")
	}

	return prom, datasource, nil
}
