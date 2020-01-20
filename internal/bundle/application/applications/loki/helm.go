package loki

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	amlogs "github.com/caos/boom/internal/bundle/application/applications/ambassador/logs"
	aglogs "github.com/caos/boom/internal/bundle/application/applications/argocd/logs"
	ksmlogs "github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/logs"
	"github.com/caos/boom/internal/bundle/application/applications/loki/helm"
	pnelogs "github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/logs"
	pologs "github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/logs"

	"github.com/caos/boom/internal/bundle/application/resources/logging"
	"github.com/caos/boom/internal/templator/helm/chart"
)

func (l *Loki) HelmPreApplySteps(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) ([]interface{}, error) {

	flows := make([]*logging.Flow, 0)
	outputs := make([]*logging.Output, 0)
	outputName := "output-loki"
	outputs = append(outputs, logging.NewOutput(outputName, "caos-system", "http://loki.caos-system:3100"))
	outputNames := make([]string, 0)
	outputNames = append(outputNames, outputName)

	if toolsetCRDSpec.Ambassador != nil && toolsetCRDSpec.Ambassador.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.Ambassador) {
		flows = append(flows, logging.NewFlow(amlogs.GetFlow(outputNames)))
	}

	if toolsetCRDSpec.PrometheusOperator != nil && toolsetCRDSpec.PrometheusOperator.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.PrometheusOperator) {
		flows = append(flows, logging.NewFlow(pologs.GetFlow(outputNames)))
	}

	if toolsetCRDSpec.PrometheusNodeExporter != nil && toolsetCRDSpec.PrometheusNodeExporter.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.PrometheusNodeExporter) {
		flows = append(flows, logging.NewFlow(pnelogs.GetFlow(outputNames)))
	}

	if toolsetCRDSpec.KubeStateMetrics != nil && toolsetCRDSpec.KubeStateMetrics.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.KubeStateMetrics) {
		flows = append(flows, logging.NewFlow(ksmlogs.GetFlow(outputNames)))
	}

	if toolsetCRDSpec.Argocd != nil && toolsetCRDSpec.Argocd.Deploy &&
		(toolsetCRDSpec.Logs == nil || toolsetCRDSpec.Logs.Argocd) {
		flows = append(flows, logging.NewFlow(aglogs.GetFlow(outputNames)))
	}

	ret := make([]interface{}, 0)
	if len(flows) > 0 {
		for _, flow := range flows {
			ret = append(ret, flow)
		}
		for _, output := range outputs {
			ret = append(ret, output)
		}

		nameLogging := "logging"
		namespaceLogging := "caos-system"
		ret = append(ret, logging.New(nameLogging, namespaceLogging, "caos-system"))

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
