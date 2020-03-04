package logs

import (
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"
)

func GetFlow(outputs []string) *logging.FlowConfig {
	ls := map[string]string{
		"app":        "prometheus",
		"prometheus": "caos-prometheus",
	}

	return &logging.FlowConfig{
		Name:         "flow-prometheus",
		Namespace:    "caos-system",
		SelectLabels: ls,
		Outputs:      outputs,
		ParserType:   "none",
	}
}
