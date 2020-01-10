package prometheusnodeexporter

import "github.com/caos/boom/internal/app/bundle/application/chart"

func (p *PrometheusNodeExporter) GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "prometheus-node-exporter",
		Version: "1.8.0",
	}
}

func (c *PrometheusNodeExporter) GetImageTags() map[string]string {
	return map[string]string{
		"quay.io/prometheus/node-exporter": "v0.18.1",
	}
}
