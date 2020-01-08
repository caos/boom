package crd

import (
	"path/filepath"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/v1beta1/crd/grafana"
)

var orgID = 0

func (c *Crd) GetGrafanaConfig(dashboardsfolder string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) *grafana.Config {
	providers := getGrafanaDashboards(dashboardsfolder, toolsetCRDSpec)
	conf := grafana.NewConfig(toolsetCRDSpec.KubeVersion, toolsetCRDSpec.Grafana)

	for _, provider := range providers {
		conf.AddDashboardProvider(provider)
	}
	return conf
}

func getGrafanaDashboards(dashboardsfolder string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) []*grafana.ConfigProvider {
	providers := make([]*grafana.ConfigProvider, 0)
	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy {
		provider := &grafana.ConfigProvider{
			ConfigMaps: []string{
				"grafana-dashboard-ambassador-envoy-global",
				"grafana-dashboard-ambassador-envoy-ingress",
				"grafana-dashboard-ambassador-envoy-service",
			},
			Folder: filepath.Join(dashboardsfolder, "ambassador"),
		}
		providers = append(providers, provider)
	}
	if toolsetCRDSpec.Argocd != nil && toolsetCRDSpec.Argocd.Deploy {
		provider := &grafana.ConfigProvider{
			ConfigMaps: []string{
				"grafana-dashboard-argocd",
			},
			Folder: filepath.Join(dashboardsfolder, "argocd"),
		}
		providers = append(providers, provider)
	}

	if toolsetCRDSpec.PrometheusNodeExporter != nil && toolsetCRDSpec.PrometheusNodeExporter.Deploy {
		provider := &grafana.ConfigProvider{
			ConfigMaps: []string{
				"grafana-dashboard-node-cluster-rsrc-use",
				"grafana-dashboard-node-rsrc-use",
			},
			Folder: filepath.Join(dashboardsfolder, "prometheusnodeexporter"),
		}
		providers = append(providers, provider)
	}

	// provider := &grafana.ConfigProvider{
	// 	ConfigMaps: []string{
	// 		"grafana-dashboard-kubelet",
	// 	},
	// 	Folder: filepath.Join(dashboardsfolder, "kubelet"),
	// }
	// providers = append(providers, provider)
	return providers
}
