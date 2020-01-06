package crd

import (
	"path/filepath"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

var orgID = 0

func (c *Crd) GetGrafanaDashboards(dashboardsfolder string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) []*toolsetsv1beta1.Provider {
	providers := make([]*toolsetsv1beta1.Provider, 0)
	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy {
		provider := &toolsetsv1beta1.Provider{
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
		provider := &toolsetsv1beta1.Provider{
			ConfigMaps: []string{
				"grafana-dashboard-argocd",
			},
			Folder: filepath.Join(dashboardsfolder, "argocd"),
		}
		providers = append(providers, provider)
	}

	provider := &toolsetsv1beta1.Provider{
		ConfigMaps: []string{
			"grafana-dashboard-k8s-cluster-rsrc-use",
			"grafana-dashboard-k8s-node-rsrc-use",
			"grafana-dashboard-k8s-resources-cluster",
			"grafana-dashboard-k8s-resources-namespace",
			"grafana-dashboard-k8s-resources-pod",
			"grafana-dashboard-k8s-resources-workload",
			"grafana-dashboard-k8s-resources-workloads-namespace",
			"grafana-dashboard-nodes",
			"grafana-dashboard-persistentvolumesusage",
			"grafana-dashboard-pods",
			"grafana-dashboard-statefulset",
			"grafana-dashboard-cluster-overview",
		},
		Folder: filepath.Join(dashboardsfolder, "cluster"),
	}
	providers = append(providers, provider)
	return providers
}
