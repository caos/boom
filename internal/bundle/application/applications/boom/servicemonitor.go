package boom

import (
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/boom/internal/name"
)

func GetServicemonitor(instanceName string) *servicemonitor.Config {
	var appName name.Application
	appName = "boom"
	monitorlabels := labels.GetMonitorLabels(instanceName, appName)
	ls := map[string]string{
		"app.kubernetes.io/instance": "boom",
		"app.kubernetes.io/part-of": "orbos",
		"app.kubernetes.io/component": "boom",
	}

	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "metrics",
		Path: "/metrics",
	}


	return &servicemonitor.Config{
		Name:                  "boom-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               "boom",
	}
}
