package loggingoperator

import (
	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator/info"
	"github.com/caos/boom/internal/name"
)

type LoggingOperator struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.LoggingOperator
}

func New(logger logging.Logger) *LoggingOperator {
	lo := &LoggingOperator{
		logger: logger,
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
