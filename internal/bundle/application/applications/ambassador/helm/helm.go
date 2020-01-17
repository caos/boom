package helm

import "github.com/caos/boom/internal/templator/helm/chart"

func GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "ambassador",
		Version: "6.0.0",
		Index: &chart.Index{
			Name: "datawire",
			URL:  "www.getambassador.io",
		},
	}
}

func GetImageTags() map[string]string {
	return map[string]string{
		"quay.io/datawire/aes": "1.0.0",
		"prom/statsd-exporter": "v0.8.1",
	}
}
