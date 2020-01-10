package prometheusoperator

import "github.com/caos/boom/internal/app/bundle/application/applications/prometheus/servicemonitor"

func GetServicemonitor(monitorlabels map[string]string) *servicemonitor.Config {
	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "http",
		Path: "/metrics",
	}
	labels := map[string]string{"app": "prometheus-operator-operator"}

	return &servicemonitor.Config{
		Name:                  "prometheus-operator-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
	}
}
