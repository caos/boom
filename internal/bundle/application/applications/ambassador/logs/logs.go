package logs

import (
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"
)

func GetFlow(outputs []string) *logging.FlowConfig {
	ls := map[string]string{
		"app.kubernetes.io/instance": "ambassador",
		"app.kubernetes.io/name":     "ambassador",
	}

	return &logging.FlowConfig{
		Name:         "flow-ambassador",
		Namespace:    "caos-system",
		SelectLabels: ls,
		Outputs:      outputs,
		ParserType:   "logfmt",
	}
}
