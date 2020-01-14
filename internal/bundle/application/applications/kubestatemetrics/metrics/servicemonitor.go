package metrics

import "github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"

func GetServicemonitor(monitorlabels map[string]string) *servicemonitor.Config {

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

	return &servicemonitor.Config{
		Name:                  "kube-state-metrics-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
		JobName:               "kube-state-metrics",
	}
}
