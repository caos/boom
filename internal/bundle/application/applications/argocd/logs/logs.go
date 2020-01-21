package logs

import "github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"

func GetFlow(outputs []string) *logging.FlowConfig {
	lables := map[string]string{"app.kubernetes.io/instance": "argocd"}

	return &logging.FlowConfig{
		Name:         "flow-argocd",
		Namespace:    "caos-system",
		SelectLabels: lables,
		Outputs:      outputs,
		ParserType:   "none",
	}
}
