package logs

import "github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"

func GetFlow(outputs []string) *logging.FlowConfig {
	lables := map[string]string{"app": "prometheus", "prometheus": "operated-prometheus"}

	return &logging.FlowConfig{
		Name:         "flow-prometheus",
		Namespace:    "caos-system",
		SelectLabels: lables,
		Outputs:      outputs,
		ParserType:   "none",
	}
}
