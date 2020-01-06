package crd

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheus"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheus/servicemonitor"
)

func (c *Crd) ScrapeMetricsCrdsConfig(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) (*prometheus.Config, string, error) {

	monitorlabels := make(map[string]string, 0)
	monitorlabels["app.kubernetes.io/managed-by"] = "boom.caos.ch"

	servicemonitors := make([]*servicemonitor.Config, 0)

	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy {
		endpoint := &servicemonitor.ConfigEndpoint{
			Port: "ambassador-admin",
			Path: "/metrics",
		}
		labels := map[string]string{"service": "ambassador-admin"}

		smconfig := &servicemonitor.Config{
			Name:                  "ambassador-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	if toolsetCRDSpec.CertManager != nil && toolsetCRDSpec.CertManager.Deploy {
		endpoint := &servicemonitor.ConfigEndpoint{
			TargetPort: "9402",
			Path:       "/metrics",
		}
		labels := map[string]string{"app": "cert-manager"}

		smconfig := &servicemonitor.Config{
			Name:                  "cert-manager-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	if toolsetCRDSpec.PrometheusOperator != nil && toolsetCRDSpec.PrometheusOperator.Deploy {
		endpoint := &servicemonitor.ConfigEndpoint{
			Port: "http",
			Path: "/metrics",
		}
		labels := map[string]string{"app": "prometheus-operator-operator"}

		smconfig := &servicemonitor.Config{
			Name:                  "prometheus-operator-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	if toolsetCRDSpec.PrometheusNodeExporter != nil && toolsetCRDSpec.PrometheusNodeExporter.Deploy {
		relabelings := make([]*servicemonitor.ConfigRelabeling, 0)
		relabeling := &servicemonitor.ConfigRelabeling{
			Action:       "replace",
			Regex:        "(.*)",
			Replacement:  "$1",
			SourceLabels: []string{"__meta_kubernetes_pod_node_name"},
			TargetLabel:  "instance",
		}
		relabelings = append(relabelings, relabeling)

		endpoint := &servicemonitor.ConfigEndpoint{
			Port:        "metrics",
			Path:        "/metrics",
			Relabelings: relabelings,
		}
		labels := map[string]string{"app": "prometheus-node-exporter"}

		smconfig := &servicemonitor.Config{
			Name:                  "prometheus-node-exporter-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
			JobName:               "node-exporter",
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	if toolsetCRDSpec.KubeStateMetrics != nil && toolsetCRDSpec.KubeStateMetrics.Deploy {

		relabelings := make([]*servicemonitor.ConfigRelabeling, 0)
		relabeling := &servicemonitor.ConfigRelabeling{
			Action: "labeldrop",
			Regex:  "(pod|service|endpoint|namespace)",
		}
		relabelings = append(relabelings, relabeling)

		endpoint := &servicemonitor.ConfigEndpoint{
			Port:        "http",
			Path:        "/metrics",
			HonorLabels: true,
			Relabelings: relabelings,
		}

		labels := map[string]string{
			"app.kubernetes.io/name": "kube-state-metrics",
		}

		smconfig := &servicemonitor.Config{
			Name:                  "kube-state-metrics-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
			JobName:               "kube-state-metrics",
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	adconfigs := make([]*prometheus.AdditionalScrapeConfig, 0)
	//Metrics from the kubernetes service
	smconfig := getAPIServers(monitorlabels)
	servicemonitors = append(servicemonitors, smconfig)

	adconfig := getCadvisor()
	adconfigs = append(adconfigs, adconfig)
	adconfig = getNodes()
	adconfigs = append(adconfigs, adconfig)

	if len(servicemonitors) > 0 {
		endpoint := &servicemonitor.ConfigEndpoint{
			Port: "web",
			Path: "/metrics",
		}
		labels := map[string]string{"app": "prometheus-operator-prometheus", "release": "caos"}

		smconfig := &servicemonitor.Config{
			Name:                  "prometheus-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)

		prom := &prometheus.Config{
			Prefix:                  "",
			Namespace:               "caos-system",
			MonitorLabels:           monitorlabels,
			ServiceMonitors:         servicemonitors,
			AdditionalScrapeConfigs: adconfigs,
		}

		datasource := ""
		if prom.Prefix != "" {
			datasource = strings.Join([]string{"http://", prom.Prefix, "-prometheus-operated.", prom.Namespace, ":9090"}, "")
		} else {
			datasource = strings.Join([]string{"http://prometheus-operated.", prom.Namespace, ":9090"}, "")
		}

		return prom, datasource, nil
	}

	return nil, "", nil
}

func getAPIServers(monitorlabels map[string]string) *servicemonitor.Config {
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

	smconfig := &servicemonitor.Config{
		Name:                  "kubernetes-apiservers-servicemonitor",
		Endpoints:             endpoints,
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
		NamespaceSelector:     []string{"default"},
		JobName:               "kubernetes-apiservers",
	}
	return smconfig
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
