package ambassador

import "github.com/caos/boom/internal/app/bundle/application/chart"

func (a *Ambassador) GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "ambassador",
		Version: "5.3.0",
	}
}

func (a *Ambassador) GetImageTags() map[string]string {
	return map[string]string{
		"quay.io/datawire/ambassador":     "0.86.1",
		"quay.io/datawire/ambassador_pro": "0.11.0",
		"prom/statsd-exporter":            "v0.9.0",
	}
}
