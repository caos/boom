package helm

import "github.com/caos/boom/internal/templator/helm/chart"

func GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "loki",
		Version: "0.22.0",
		Index: &chart.Index{
			Name: "loki",
			URL:  "grafana.github.io/loki/charts",
		},
	}
}

func GetImageTags() map[string]string {
	return map[string]string{
		"grafana/loki": "v1.2.0",
	}
}
