package prometheus

func getScrapeConfigs() []*AdditionalScrapeConfig {

	adconfigs := make([]*AdditionalScrapeConfig, 0)

	adconfigs = append(adconfigs, getNodes())
	adconfigs = append(adconfigs, getCadvisor())

	return adconfigs
}

func getNodes() *AdditionalScrapeConfig {
	relabelings := make([]*RelabelConfig, 0)
	relabeling := &RelabelConfig{
		Action: "labelmap",
		Regex:  "__meta_kubernetes_node_label_(.+)",
	}
	relabelings = append(relabelings, relabeling)
	relabeling = &RelabelConfig{
		TargetLabel: "__address__",
		Replacement: "kubernetes.default.svc:443",
	}
	relabelings = append(relabelings, relabeling)
	relabeling = &RelabelConfig{
		SourceLabels: []string{"__meta_kubernetes_node_name"},
		Regex:        "(.+)",
		TargetLabel:  "__metrics_path__",
		Replacement:  "/api/v1/nodes/${1}/proxy/metrics",
	}
	relabelings = append(relabelings, relabeling)

	sdconfig := &KubernetesSdConfig{
		Role: "node",
	}

	return &AdditionalScrapeConfig{
		JobName:             "kubernetes-nodes",
		Scheme:              "https",
		KubernetesSdConfigs: []*KubernetesSdConfig{sdconfig},
		BearerTokenFile:     "/var/run/secrets/kubernetes.io/serviceaccount/token",
		TLSConfig: &TLSConfig{
			CaFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		},
		RelabelConfigs: relabelings,
	}
}

func getCadvisor() *AdditionalScrapeConfig {
	relabelings := make([]*RelabelConfig, 0)
	relabeling := &RelabelConfig{
		Action: "labelmap",
		Regex:  "__meta_kubernetes_node_label_(.+)",
	}
	relabelings = append(relabelings, relabeling)
	relabeling = &RelabelConfig{
		TargetLabel: "__address__",
		Replacement: "kubernetes.default.svc:443",
	}
	relabelings = append(relabelings, relabeling)
	relabeling = &RelabelConfig{
		SourceLabels: []string{"__meta_kubernetes_node_name"},
		Regex:        "(.+)",
		TargetLabel:  "__metrics_path__",
		Replacement:  "/api/v1/nodes/${1}/proxy/metrics/cadvisor",
	}
	relabelings = append(relabelings, relabeling)

	sdconfig := &KubernetesSdConfig{
		Role: "node",
	}

	return &AdditionalScrapeConfig{
		JobName:             "kubelet",
		Scheme:              "https",
		KubernetesSdConfigs: []*KubernetesSdConfig{sdconfig},
		BearerTokenFile:     "/var/run/secrets/kubernetes.io/serviceaccount/token",
		TLSConfig: &TLSConfig{
			CaFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		},
		RelabelConfigs: relabelings,
	}
}
