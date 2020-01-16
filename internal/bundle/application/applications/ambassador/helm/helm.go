package helm

import "github.com/caos/boom/internal/templator/helm/chart"

func GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "ambassador",
		Version: "5.3.0",
	}
}

func GetImageTags() map[string]string {
	return map[string]string{
		"quay.io/datawire/ambassador":     "0.86.1",
		"quay.io/datawire/ambassador_pro": "0.11.0",
		"prom/statsd-exporter":            "v0.9.0",
	}
}
