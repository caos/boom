package labels

import "github.com/caos/boom/internal/name"

var (
	instanceName = "boom"
)

func getGlobalLabels() map[string]string {
	labels := make(map[string]string, 0)
	labels["app.kubernetes.io/managed-by"] = "boom.caos.ch"
	labels["boom.caos.ch/instance"] = instanceName

	return labels
}

func GetApplicationLabels(appName name.Application) map[string]string {
	labels := getGlobalLabels()
	labels["boom.caos.ch/application"] = appName.String()
	return labels
}

func GetMonitorLabels(instanceName string) map[string]string {
	labels := getGlobalLabels()
	labels["boom.caos.ch/prometheus"] = instanceName
	return labels
}

func GetRuleLabels(instanceName string) map[string]string {
	labels := getGlobalLabels()
	labels["boom.caos.ch/prometheus"] = instanceName
	return labels
}
