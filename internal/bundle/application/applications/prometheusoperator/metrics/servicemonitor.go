package metrics

import (
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator"
	"github.com/caos/boom/internal/labels"
)

func GetServicemonitor(instanceName string) *servicemonitor.Config {
	appName := prometheusoperator.GetName()
	monitorlabels := labels.GetMonitorLabels(instanceName)
	ls := labels.GetApplicationLabels(appName)

	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "http",
		Path: "/metrics",
	}

	return &servicemonitor.Config{
		Name:                  "prometheus-operator-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               appName.String(),
	}
}
