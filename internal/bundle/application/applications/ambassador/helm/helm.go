package helm

import "github.com/caos/boom/internal/templator/helm/chart"

func GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "ambassador",
		Version: "6.1.3",
		Index: &chart.Index{
			Name: "datawire",
			URL:  "www.getambassador.io",
		},
	}
}

func GetImageTags() map[string]string {
	return map[string]string{
		"quay.io/datawire/aes": "1.2.1",
		"prom/statsd-exporter": "v0.8.1",
	}
}
