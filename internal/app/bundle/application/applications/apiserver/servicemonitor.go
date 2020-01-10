package apiserver

import "github.com/caos/boom/internal/app/bundle/application/applications/prometheus/servicemonitor"

func GetServicemonitor(monitorlabels map[string]string) *servicemonitor.Config {
	metricsRelabelings := make([]*servicemonitor.ConfigMetricRelabeling, 0)
	relabeling := &servicemonitor.ConfigMetricRelabeling{
		Action:       "keep",
		Regex:        "default;kubernetes;https",
		SourceLabels: []string{"__meta_kubernetes_namespace", "__meta_kubernetes_service_name", "__meta_kubernetes_endpoint_port_name"},
	}
	metricsRelabelings = append(metricsRelabelings, relabeling)

	endpoints := make([]*servicemonitor.ConfigEndpoint, 0)
	endpoint := &servicemonitor.ConfigEndpoint{
		Scheme:          "https",
		BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
		Port:            "https",
		Path:            "/metrics",
		TLSConfig: &servicemonitor.ConfigTLSConfig{
			CaFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		},
		MetricsRelabelings: metricsRelabelings,
	}
	endpoints = append(endpoints, endpoint)

	labels := map[string]string{
		"component": "apiserver",
		"provider":  "kubernetes",
	}

	return &servicemonitor.Config{
		Name:                  "kubernetes-apiservers-servicemonitor",
		Endpoints:             endpoints,
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
		NamespaceSelector:     []string{"default"},
		JobName:               "kubernetes-apiservers",
	}
}