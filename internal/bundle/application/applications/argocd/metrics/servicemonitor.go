package metrics

import "github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"

func GetServicemonitors(monitorlabels map[string]string) []*servicemonitor.Config {

	servicemonitors := make([]*servicemonitor.Config, 0)

	//argocd-application-controller
	endpoint := &servicemonitor.ConfigEndpoint{
		Port: "metrics",
		Path: "/metrics",
	}

	labels := map[string]string{
		"app.kubernetes.io/part-of":   "argocd",
		"app.kubernetes.io/component": "application-controller",
	}

	smconfig := &servicemonitor.Config{
		Name:                  "application-controller-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
		JobName:               "application-controller",
	}
	servicemonitors = append(servicemonitors, smconfig)
	// argocd-repo-server
	endpoint = &servicemonitor.ConfigEndpoint{
		Port: "metrics",
		Path: "/metrics",
	}

	labels = map[string]string{
		"app.kubernetes.io/part-of":   "argocd",
		"app.kubernetes.io/component": "repo-server",
	}

	smconfig = &servicemonitor.Config{
		Name:                  "argocd-repo-server-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
		JobName:               "argocd-repo-server",
	}
	servicemonitors = append(servicemonitors, smconfig)
	// argocd-server
	endpoint = &servicemonitor.ConfigEndpoint{
		Port: "metrics",
		Path: "/metrics",
	}

	labels = map[string]string{
		"app.kubernetes.io/part-of":   "argocd",
		"app.kubernetes.io/component": "server",
	}

	smconfig = &servicemonitor.Config{
		Name:                  "argocd-server-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		MonitorMatchingLabels: monitorlabels,
		ServiceMatchingLabels: labels,
		JobName:               "argocd-server",
	}
	return servicemonitors
}
