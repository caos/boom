package metrics

import (
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/labels"
)

func GetServicemonitors(instanceName string) []*servicemonitor.Config {

	servicemonitors := make([]*servicemonitor.Config, 0)

	servicemonitors = append(servicemonitors, getSMApplicationController(instanceName))
	servicemonitors = append(servicemonitors, getSMRepoServer(instanceName))
	servicemonitors = append(servicemonitors, getSMServer(instanceName))

	return servicemonitors
}

func getSMServer(instanceName string) *servicemonitor.Config {
	monitorlabels := labels.GetMonitorLabels(instanceName)
	ls := make(map[string]string, 0)

	// argocd-server
	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "metrics",
		Path: "/metrics",
	}

	ls["app.kubernetes.io/instance"] = "argocd"
	ls["app.kubernetes.io/part-of"] = "argocd"
	ls["app.kubernetes.io/component"] = "server"

	return &servicemonitor.Config{
		Name:                  "argocd-server-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               "argocd-server",
	}
}

func getSMRepoServer(instanceName string) *servicemonitor.Config {
	monitorlabels := labels.GetMonitorLabels(instanceName)
	ls := make(map[string]string, 0)

	// argocd-repo-server
	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "metrics",
		Path: "/metrics",
	}

	ls["app.kubernetes.io/instance"] = "argocd"
	ls["app.kubernetes.io/part-of"] = "argocd"
	ls["app.kubernetes.io/component"] = "repo-server"

	return &servicemonitor.Config{
		Name:                  "argocd-repo-server-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               "argocd-repo-server",
	}

}

func getSMApplicationController(instanceName string) *servicemonitor.Config {
	monitorlabels := labels.GetMonitorLabels(instanceName)
	ls := make(map[string]string, 0)

	//argocd-application-controller
	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "metrics",
		Path: "/metrics",
	}

	ls["app.kubernetes.io/instance"] = "argocd"
	ls["app.kubernetes.io/part-of"] = "argocd"
	ls["app.kubernetes.io/component"] = "application-controller"

	return &servicemonitor.Config{
		Name:                  "application-controller-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: ls,
		JobName:               "application-controller",
	}
}
