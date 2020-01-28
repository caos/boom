package logs

import (
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"
)

func GetFlow(outputs []string) *logging.FlowConfig {
	lables := map[string]string{"release": "grafana", "app": "grafana"}

	return &logging.FlowConfig{
		Name:         "flow-ambassador",
		Namespace:    "caos-system",
		SelectLabels: lables,
		Outputs:      outputs,
		ParserType:   "logfmt",
	}

}
