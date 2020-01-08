package kubelet

import "github.com/caos/boom/internal/app/v1beta1/crd/prometheus"

func GetScrapeConfigs() []*prometheus.AdditionalScrapeConfig {

	adconfigs := make([]*prometheus.AdditionalScrapeConfig, 0)

	adconfigs = append(adconfigs, getNodes())
	adconfigs = append(adconfigs, getCadvisor())

	return adconfigs
}

func getNodes() *prometheus.AdditionalScrapeConfig {
	relabelings := make([]*prometheus.RelabelConfig, 0)
	relabeling := &prometheus.RelabelConfig{
		Action: "labelmap",
		Regex:  "__meta_kubernetes_node_label_(.+)",
	}
	relabelings = append(relabelings, relabeling)
	relabeling = &prometheus.RelabelConfig{
		TargetLabel: "__address__",
		Replacement: "kubernetes.default.svc:443",
	}
	relabelings = append(relabelings, relabeling)
	relabeling = &prometheus.RelabelConfig{
		SourceLabels: []string{"__meta_kubernetes_node_name"},
		Regex:        "(.+)",
		TargetLabel:  "__metrics_path__",
		Replacement:  "/api/v1/nodes/${1}/proxy/metrics",
	}
	relabelings = append(relabelings, relabeling)

	sdconfig := &prometheus.KubernetesSdConfig{
		Role: "node",
	}

	return &prometheus.AdditionalScrapeConfig{
		JobName:             "kubernetes-nodes",
		Scheme:              "https",
		KubernetesSdConfigs: []*prometheus.KubernetesSdConfig{sdconfig},
		BearerTokenFile:     "/var/run/secrets/kubernetes.io/serviceaccount/token",
		TLSConfig: &prometheus.TLSConfig{
			CaFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		},
		RelabelConfigs: relabelings,
	}
}

func getCadvisor() *prometheus.AdditionalScrapeConfig {
	relabelings := make([]*prometheus.RelabelConfig, 0)
	relabeling := &prometheus.RelabelConfig{
		Action: "labelmap",
		Regex:  "__meta_kubernetes_node_label_(.+)",
	}
	relabelings = append(relabelings, relabeling)
	relabeling = &prometheus.RelabelConfig{
		TargetLabel: "__address__",
		Replacement: "kubernetes.default.svc:443",
	}
	relabelings = append(relabelings, relabeling)
	relabeling = &prometheus.RelabelConfig{
		SourceLabels: []string{"__meta_kubernetes_node_name"},
		Regex:        "(.+)",
		TargetLabel:  "__metrics_path__",
		Replacement:  "/api/v1/nodes/${1}/proxy/metrics/cadvisor",
	}
	relabelings = append(relabelings, relabeling)

	sdconfig := &prometheus.KubernetesSdConfig{
		Role: "node",
	}

	return &prometheus.AdditionalScrapeConfig{
		JobName:             "kubelet",
		Scheme:              "https",
		KubernetesSdConfigs: []*prometheus.KubernetesSdConfig{sdconfig},
		BearerTokenFile:     "/var/run/secrets/kubernetes.io/serviceaccount/token",
		TLSConfig: &prometheus.TLSConfig{
			CaFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		},
		RelabelConfigs: relabelings,
	}
}
