package metrics

import (
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/labels"
)

func getLocalServiceMonitor(monitorMatchingLabels, serviceMatchingLabels map[string]string) *servicemonitor.Config {
	relabelings := make([]*servicemonitor.ConfigRelabeling, 0)
	relabeling := &servicemonitor.ConfigRelabeling{
		Action: "labeldrop",
		Regex:  "(pod|service|endpoint|namespace)",
	}
	relabelings = append(relabelings, relabeling)

	endpoint := &servicemonitor.ConfigEndpoint{
		Port:        "http",
		Path:        "/metrics",
		HonorLabels: true,
		Relabelings: relabelings,
	}

	return &servicemonitor.Config{
		Name:                  "local-kube-state-metrics-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		JobName:               "local-kube-state-metrics",
		MonitorMatchingLabels: monitorMatchingLabels,
		ServiceMatchingLabels: serviceMatchingLabels,
	}
}

func getIngestionServiceMonitor(monitorMatchingLabels, serviceMatchingLabels map[string]string) *servicemonitor.Config {

	relabelings := []*servicemonitor.ConfigRelabeling{{
		Action: "labeldrop",
		Regex:  "(container|endpoint|namespace|pod)",
	}}

	metricRelabelings := []*servicemonitor.ConfigRelabeling{{
		Action: "keep",
		SourceLabels: []string{
			"__name__",
			"condition",
			"status",
		},
		Regex: "kube_node_status_condition;Ready;true|(kube_deployment_status_replicas|kube_deployment_spec_replicas|kube_deployment_status_replicas_available|kube_statefulset_status_replicas_current|kube_statefulset_replicas|kube_statefulset_status_replicas_ready|kube_daemonset_status_current_number_scheduled|kube_daemonset_status_desired_number_scheduled|kube_daemonset_status_number_available);;",
	}, {
		Action: "replace",
		SourceLabels: []string{
			"daemonset",
			"statefulset",
			"deployment",
		},
		TargetLabel: "controller",
		Regex:       "(.*);(.*);(.*)",
		Replacement: "${1}${2}${3}",
	}, {
		Action: "labelkeep",
		Regex:  "__.+|job|node|controller",
	}}

	endpoint := &servicemonitor.ConfigEndpoint{
		Port:              "http",
		Path:              "/metrics",
		HonorLabels:       false,
		Relabelings:       relabelings,
		MetricRelabelings: metricRelabelings,
	}

	return &servicemonitor.Config{
		Name:                  "ingestion-kube-state-metrics-servicemonitor",
		Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
		JobName:               "ingestion-kube-state-metrics",
		MonitorMatchingLabels: monitorMatchingLabels,
		ServiceMatchingLabels: serviceMatchingLabels,
	}
}

func GetServicemonitors(instanceName string) []*servicemonitor.Config {

	appName := info.GetName()
	monitorLabels := labels.GetMonitorLabels(instanceName, appName)
	ls := labels.GetApplicationLabels(appName)

	return []*servicemonitor.Config{
		getLocalServiceMonitor(monitorLabels, ls),
		getIngestionServiceMonitor(monitorLabels, ls),
	}
}
