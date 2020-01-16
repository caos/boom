package certmanager

import "github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"

func GetServicemonitor(monitorlabels map[string]string) *servicemonitor.Config {
	endpoint := &servicemonitor.ConfigEndpoint{
		TargetPort: "http",
		Path:       "/metrics",
	}
	labels := map[string]string{"app": "cert-manager"}

	return &servicemonitor.Config{
		Name:                  "cert-manager-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
	}
}
