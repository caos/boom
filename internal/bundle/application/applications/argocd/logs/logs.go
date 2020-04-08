package logs

import (
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"
)

func GetFlow(outputs []string) *logging.FlowConfig {
	ls := map[string]string{
		"app.kubernetes.io/instance": "argocd",
		"app.kubernetes.io/part-of":  "orbos",
	}

	return &logging.FlowConfig{
		Name:         "flow-argocd",
		Namespace:    "caos-system",
		SelectLabels: ls,
		Outputs:      outputs,
		ParserType:   "logfmt",
	}
}
