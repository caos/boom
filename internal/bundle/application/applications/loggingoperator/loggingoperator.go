package loggingoperator

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
)

type LoggingOperator struct {
	monitor mntr.Monitor
	spec    *toolsetsv1beta1.LoggingOperator
}

func New(monitor mntr.Monitor) *LoggingOperator {
	lo := &LoggingOperator{
		monitor: monitor,
	}

	return lo
}
func (l *LoggingOperator) GetName() name.Application {
	return info.GetName()
}

func (lo *LoggingOperator) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.LoggingOperator.Deploy
}

func (l *LoggingOperator) GetNamespace() string {
	return info.GetNamespace()
}
