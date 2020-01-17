package config

import (
	"path/filepath"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

var orgID = 0

func getGrafanaDashboards(dashboardsfolder string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) []*Provider {
	providers := make([]*Provider, 0)
	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy &&
		(toolsetCRDSpec.Metrics == nil || toolsetCRDSpec.Metrics.Ambassador) {
		provider := &Provider{
			ConfigMaps: []string{
				"grafana-dashboard-ambassador-envoy-global",
				"grafana-dashboard-ambassador-envoy-ingress",
				"grafana-dashboard-ambassador-envoy-service",
			},
			Folder: filepath.Join(dashboardsfolder, "ambassador"),
		}
		providers = append(providers, provider)
	}

	if toolsetCRDSpec.Argocd != nil && toolsetCRDSpec.Argocd.Deploy &&
		(toolsetCRDSpec.Metrics == nil || toolsetCRDSpec.Metrics.Argocd) {
		provider := &Provider{
			ConfigMaps: []string{
				"grafana-dashboard-argocd",
			},
			Folder: filepath.Join(dashboardsfolder, "argocd"),
		}
		providers = append(providers, provider)
	}

	if toolsetCRDSpec.PrometheusNodeExporter != nil && toolsetCRDSpec.PrometheusNodeExporter.Deploy &&
		(toolsetCRDSpec.Metrics == nil || toolsetCRDSpec.Metrics.PrometheusNodeExporter) {
		provider := &Provider{
			ConfigMaps: []string{
				"grafana-dashboard-node-cluster-rsrc-use",
				"grafana-dashboard-node-rsrc-use",
			},
			Folder: filepath.Join(dashboardsfolder, "prometheusnodeexporter"),
		}
		providers = append(providers, provider)
	}

	// provider := &ConfigProvider{
	// 	ConfigMaps: []string{
	// 		"grafana-dashboard-kubelet",
	// 	},
	// 	Folder: filepath.Join(dashboardsfolder, "kubelet"),
	// }
	// providers = append(providers, provider)
	return providers
}
