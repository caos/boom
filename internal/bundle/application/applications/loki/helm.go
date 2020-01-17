package loki

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/loki/helm"
	"github.com/caos/boom/internal/bundle/application/resources/logging"
	"github.com/caos/boom/internal/templator/helm/chart"
)

func (l *Loki) HelmPreApplySteps(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) ([]interface{}, error) {

	monitorlabels := make(map[string]string, 0)
	monitorlabels["app.kubernetes.io/managed-by"] = "boom.caos.ch"

	flows := make([]*logging.Flow, 0)
	outputs := make([]*logging.Output, 0)
	outputName := "output-loki"
	outputs = append(outputs, logging.NewOutput(outputName, "caos-system", "https://loki.caos-system:3100"))
	outputNames := make([]string, 0)
	outputNames = append(outputNames, outputName)

	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.Ambassador) {
		flow := logging.NewFlow("flow-ambassador", "caos-system", map[string]string{"app.kubernetes.io/part-of": "ambassador"}, outputNames)
		flows = append(flows, flow)
	}

	if toolsetCRDSpec.PrometheusOperator != nil && toolsetCRDSpec.PrometheusOperator.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.PrometheusOperator) {
		flow := logging.NewFlow("flow-prometheus-operator", "caos-system", map[string]string{"app.kubernetes.io/part-of": "prometheus-operator"}, outputNames)
		flows = append(flows, flow)
	}

	if toolsetCRDSpec.PrometheusNodeExporter != nil && toolsetCRDSpec.PrometheusNodeExporter.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.PrometheusNodeExporter) {
		flow := logging.NewFlow("flow-prometheus-node-exporter", "caos-system", map[string]string{"app.kubernetes.io/part-of": "prometheus-node-exporter"}, outputNames)
		flows = append(flows, flow)
	}

	if toolsetCRDSpec.KubeStateMetrics != nil && toolsetCRDSpec.KubeStateMetrics.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.KubeStateMetrics) {
		flow := logging.NewFlow("flow-kube-state-metrics", "caos-system", map[string]string{"app.kubernetes.io/part-of": "kube-state-metrics"}, outputNames)
		flows = append(flows, flow)
	}

	if toolsetCRDSpec.Argocd != nil && toolsetCRDSpec.Argocd.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.Argocd) {
		flow := logging.NewFlow("flow-argocd", "caos-system", map[string]string{"app.kubernetes.io/part-of": "argocd"}, outputNames)
		flows = append(flows, flow)
	}

	ret := make([]interface{}, 0)
	if len(flows) > 0 {
		for _, flow := range flows {
			ret = append(ret, flow)
		}
		for _, output := range outputs {
			ret = append(ret, output)
		}

		ret = append(ret, logging.New("logging", "caos-system", "caos-system"))
	}

	return ret, nil
}

func (l *Loki) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.LoggingOperator
	values := helm.DefaultValues(l.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	return values
}

func (l *Loki) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (l *Loki) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
