package helm

import "github.com/caos/boom/internal/templator/helm/chart"

func GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "loki",
		Version: "0.25.1",
		Index: &chart.Index{
			Name: "loki",
			URL:  "grafana.github.io/loki/charts",
		},
	}
}

func GetImageTags() map[string]string {
	return map[string]string{
		"grafana/loki": "v1.3.0",
	}
}
