package metrics

import (
	"strings"

	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/labels"
)

func GetServicemonitors(instanceName string) []*servicemonitor.Config {
	ret := make([]*servicemonitor.Config, 0)
	ret = append(ret, getFluentd(instanceName))
	ret = append(ret, getFluentbit(instanceName))
	return ret
}

func getFluentd(instanceName string) *servicemonitor.Config {
	appName := info.GetName()
	monitorlabels := labels.GetMonitorLabels(instanceName)
	ls := make(map[string]string, 0)

	endpoint := &servicemonitor.ConfigEndpoint{
		Port:        "metrics",
		HonorLabels: true,
	}

	ls["app.kubernetes.io/name"] = "fluentd"
	ls["app.kubernetes.io/managed-by"] = "logging"

	jobname := strings.Join([]string{appName.String(), "fluentd-metrics"}, "-")
	return &servicemonitor.Config{
		Name:                  "logging-operator-fluentd-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               jobname,
	}
}

func getFluentbit(instanceName string) *servicemonitor.Config {
	appName := info.GetName()
	monitorlabels := labels.GetMonitorLabels(instanceName)
	ls := make(map[string]string, 0)

	endpoint := &servicemonitor.ConfigEndpoint{
		Port:        "metrics",
		Path:        "/api/v1/metrics/prometheus",
		HonorLabels: true,
	}

	ls["app.kubernetes.io/name"] = "fluentbit"
	ls["app.kubernetes.io/managed-by"] = "logging"

	jobname := strings.Join([]string{appName.String(), "fluentbit-metrics"}, "-")
	return &servicemonitor.Config{
		Name:                  "logging-operator-fluentbit-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               jobname,
	}
}
