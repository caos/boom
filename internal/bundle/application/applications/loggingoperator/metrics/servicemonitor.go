package metrics

import "github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"

func GetServicemonitors(monitorlabels map[string]string) []*servicemonitor.Config {
	ret := make([]*servicemonitor.Config, 0)
	ret = append(ret, getFluentd(monitorlabels))
	ret = append(ret, getFluentbit(monitorlabels))
	return ret
}

func getFluentd(monitorlabels map[string]string) *servicemonitor.Config {

	endpoint := &servicemonitor.ConfigEndpoint{
		Port:        "metrics",
		HonorLabels: true,
	}

	labels := map[string]string{
		"app.kubernetes.io/name":       "fluentd",
		"app.kubernetes.io/managed-by": "logging",
	}

	return &servicemonitor.Config{
		Name:                  "logging-operator-fluentd-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
		JobName:               "logging-operator-fluentd-metrics",
	}
}

func getFluentbit(monitorlabels map[string]string) *servicemonitor.Config {

	endpoint := &servicemonitor.ConfigEndpoint{
		Port:        "metrics",
		Path:        "/api/v1/metrics/prometheus",
		HonorLabels: true,
	}

	labels := map[string]string{
		"app.kubernetes.io/name":       "fluentbit",
		"app.kubernetes.io/managed-by": "logging",
	}

	return &servicemonitor.Config{
		Name:                  "logging-operator-fluentbit-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
		JobName:               "logging-operator-fluentbit-metrics",
	}
}
