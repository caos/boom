package helm

import "github.com/caos/boom/internal/templator/helm/chart"

func GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "prometheus-node-exporter",
		Version: "1.8.0",
	}
}

func GetImageTags() map[string]string {
	return map[string]string{
		"quay.io/prometheus/node-exporter": "v0.18.1",
	}
}
