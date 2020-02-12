package logs

import (
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator/logging"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/info"
	"github.com/caos/boom/internal/labels"
)

func GetFlow(outputs []string) *logging.FlowConfig {
	ls := labels.GetApplicationLabels(info.GetName())

	return &logging.FlowConfig{
		Name:         "flow-prometheus-node-exporter",
		Namespace:    "caos-system",
		SelectLabels: ls,
		Outputs:      outputs,
		ParserType:   "logfmt",
	}
}
