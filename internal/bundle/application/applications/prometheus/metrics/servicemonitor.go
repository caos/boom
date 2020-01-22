package metrics

import (
	"github.com/caos/boom/internal/bundle/application/applications/prometheus"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/labels"
)

func GetServicemonitor(instanceName string) *servicemonitor.Config {
	appName := prometheus.GetName()
	monitorlabels := labels.GetMonitorLabels(instanceName)
	ls := labels.GetApplicationLabels(appName)

	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "web",
		Path: "/metrics",
	}

	return &servicemonitor.Config{
		Name:                  "prometheus-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               appName.String(),
	}
}
