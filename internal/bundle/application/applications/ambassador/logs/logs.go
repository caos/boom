package logs

import "github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"

func GetFlow(outputs []string) *logging.FlowConfig {
	lables := map[string]string{"app.kubernetes.io/instance": "ambassador", "app.kubernetes.io/name": "ambassador"}

	return &logging.FlowConfig{
		Name:         "flow-ambassador",
		Namespace:    "caos-system",
		SelectLabels: lables,
		Outputs:      outputs,
		ParserType:   "logfmt",
	}
}
