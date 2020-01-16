package helm

import "github.com/caos/boom/internal/templator/helm/chart"

func GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "kube-state-metrics",
		Version: "2.4.1",
	}
}

func GetImageTags() map[string]string {
	return map[string]string{
		"quay.io/coreos/kube-state-metrics": "v1.8.0",
	}
}
