package metrics

import (
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/labels"
)

func GetServicemonitor(instanceName string) *servicemonitor.Config {
	monitorlabels := labels.GetMonitorLabels(instanceName)
	ls := make(map[string]string, 0)

	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "ambassador-admin",
		Path: "/metrics",
	}

	ls["service"] = "ambassador-admin"
	ls["app.kubernetes.io/part-of"] = "ambassador"
	ls["app.kubernetes.io/name"] = "ambassador"
	ls["app.kubernetes.io/instance"] = "ambassador"

	return &servicemonitor.Config{
		Name:                  "ambassador-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               "ambassador",
	}
}
