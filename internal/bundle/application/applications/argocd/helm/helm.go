package helm

import "github.com/caos/boom/internal/templator/helm/chart"

func GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "argo-cd",
		Version: "1.6.4",
		Index: &chart.Index{
			Name: "argo",
			URL:  "argoproj.github.io/argo-helm",
		},
	}
}

func GetImageTags() map[string]string {
	return map[string]string{
		"argoproj/argocd":    "v1.4.0",
		"quay.io/dexidp/dex": "v2.14.0",
		"redis":              "5.0.3",
	}
}
