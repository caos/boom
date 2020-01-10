package loggingoperator

import "github.com/caos/boom/internal/app/bundle/application/chart"

func (l *LoggingOperator) GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "logging-operator",
		Version: "2.7.0",
		Index: &chart.Index{
			Name: "banzaicloud-stable",
			URL:  "kubernetes-charts.banzaicloud.com",
		},
	}
}

func (l *LoggingOperator) GetImageTags() map[string]string {
	return map[string]string{
		"banzaicloud/logging-operator": "2.7.0",
	}
}
