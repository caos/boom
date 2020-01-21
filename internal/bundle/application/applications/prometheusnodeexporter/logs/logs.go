package logs

import "github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"

func GetFlow(outputs []string) *logging.FlowConfig {
	lables := map[string]string{"release": "prometheus-node-exporter", "app": "prometheus-node-exporter"}

	return &logging.FlowConfig{
		Name:         "flow-prometheus-node-exporter",
		Namespace:    "caos-system",
		SelectLabels: lables,
		Outputs:      outputs,
		ParserType:   "none",
	}
}
