package metrics

import "github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"

func GetServicemonitor(monitorlabels map[string]string) *servicemonitor.Config {
	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "http-metrics",
		Path: "/metrics",
	}

	labels := map[string]string{
		"release": "loki",
		"app":     "loki",
	}

	return &servicemonitor.Config{
		Name:                  "loki-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
	}
}
