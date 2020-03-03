package metrics

import (
	"github.com/caos/boom/internal/bundle/application/applications/ambassador/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/labels"
)

func GetServicemonitor(instanceName string) *servicemonitor.Config {
	appName := info.GetName()
	monitorlabels := labels.GetMonitorLabels(instanceName, appName)
	ls := labels.GetApplicationLabels(appName)

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
