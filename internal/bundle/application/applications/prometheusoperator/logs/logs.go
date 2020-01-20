package logs

import "github.com/caos/boom/internal/bundle/application/resources/logging"

func GetFlow(outputs []string) *logging.FlowConfig {
	lables := map[string]string{"release": "prometheus-operator", "app": "prometheus-operator-operator"}

	return &logging.FlowConfig{
		Name:         "flow-prometheus-operator",
		Namespace:    "caos-system",
		SelectLabels: lables,
		Outputs:      outputs,
		ParserType:   "none",
	}
}
