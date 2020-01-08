package prometheus

import "github.com/caos/boom/internal/app/v1beta1/crd/prometheus/servicemonitor"

func GetServicemonitor(monitorlabels map[string]string) *servicemonitor.Config {
	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "web",
		Path: "/metrics",
	}

	labels := map[string]string{
		"operated-prometheus": "true",
	}

	return &servicemonitor.Config{
		Name:                  "prometheus-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
	}
}
