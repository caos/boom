package loki

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/loki/info"
	"github.com/caos/boom/internal/name"
)

type Loki struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.Loki
}

func New(logger logging.Logger) *Loki {
	lo := &Loki{
		logger: logger,
	}
	return lo
}

func (l *Loki) GetName() name.Application {
	return info.GetName()
}

func (lo *Loki) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.Loki.Deploy
}

func (l *Loki) Initial() bool {
	return l.spec == nil
}

func (l *Loki) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.Loki, l.spec)
}

func (l *Loki) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	l.spec = toolsetCRDSpec.Loki
}

func (l *Loki) GetNamespace() string {
	return info.GetNamespace()
}
