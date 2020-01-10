package prometheusnodeexporter

import "github.com/caos/boom/internal/app/bundle/application/applications/prometheus/servicemonitor"

func GetServicemonitor(monitorlabels map[string]string) *servicemonitor.Config {
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

	return &servicemonitor.Config{
		Name:                  "prometheus-node-exporter-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
		JobName:               "node-exporter",
	}
}
