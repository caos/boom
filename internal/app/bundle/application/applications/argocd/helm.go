package argocd

import "github.com/caos/boom/internal/app/bundle/application/chart"

func (a *Argocd) GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "argo-cd",
		Version: "1.5.3",
		Index: &chart.Index{
			Name: "argo",
			URL:  "argoproj.github.io/argo-helm",
		},
	}
}

func (c *Argocd) GetImageTags() map[string]string {
	return map[string]string{
		"argoproj/argocd":    "v1.3.6",
		"quay.io/dexidp/dex": "v2.14.0",
		"redis":              "5.0.3",
	}
}
