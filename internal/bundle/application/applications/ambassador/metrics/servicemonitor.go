package metrics

import "github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"

func GetServicemonitor(monitorlabels map[string]string) *servicemonitor.Config {
	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "ambassador-admin",
		Path: "/metrics",
	}
	labels := map[string]string{"app.kubernetes.io/name": "ambassador"}

	return &servicemonitor.Config{
		Name:                  "ambassador-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
	}
}
