package kubestatemetrics

import "github.com/caos/boom/internal/app/bundle/application/chart"

func (k *KubeStateMetrics) GetChartInfo() *chart.Chart {
	return &chart.Chart{
		Name:    "kube-state-metrics",
		Version: "2.4.1",
	}
}

func (k *KubeStateMetrics) GetImageTags() map[string]string {
	return map[string]string{
		"quay.io/coreos/kube-state-metrics": "v1.8.0",
	}
}
