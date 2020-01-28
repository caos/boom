package logs

import "github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"

func GetFlow(outputs []string) *logging.FlowConfig {
	lables := map[string]string{"app.kubernetes.io/instance": "kube-state-metrics", "app.kubernetes.io/name": "kube-state-metrics"}

	return &logging.FlowConfig{
		Name:         "flow-kube-state-metrics",
		Namespace:    "caos-system",
		SelectLabels: lables,
		Outputs:      outputs,
		ParserType:   "none",
	}
}
