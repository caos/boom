package metrics

import (
	"github.com/caos/boom/internal/bundle/application/applications/ambassador"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/labels"
)

func GetServicemonitor(instanceName string) *servicemonitor.Config {
	appName := ambassador.GetName()
	monitorlabels := labels.GetMonitorLabels(instanceName)
	ls := labels.GetApplicationLabels(appName)

	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "admin",
		Path: "/metrics",
	}

	return &servicemonitor.Config{
		Name:                  "ambassador-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               "ambassador",
	}
}
