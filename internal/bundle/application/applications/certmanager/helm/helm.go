package helm

import "github.com/caos/boom/internal/templator/helm/chart"

func GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "cert-manager",
		Version: "v0.12.0",
		Index: &chart.Index{
			Name: "jetstack",
			URL:  "charts.jetstack.io",
		},
	}
}

func GetImageTags() map[string]string {
	return map[string]string{
		"quay.io/jetstack/cert-manager-controller": "v0.12.0",
		"quay.io/jetstack/cert-manager-webhook":    "v0.12.0",
		"quay.io/jetstack/cert-manager-cainjector": "v0.12.0",
	}
}