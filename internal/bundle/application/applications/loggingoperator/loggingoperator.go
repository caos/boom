package loggingoperator

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
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
	return applicationName
}

func (lo *LoggingOperator) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.LoggingOperator.Deploy
}

func (l *LoggingOperator) Initial() bool {
	return l.spec == nil
}

func (l *LoggingOperator) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.LoggingOperator, l.spec)
}

func (l *LoggingOperator) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	l.spec = toolsetCRDSpec.LoggingOperator
}

func (l *LoggingOperator) GetNamespace() string {
	return "caos-system"
}
