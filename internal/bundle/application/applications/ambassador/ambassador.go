package ambassador

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/ambassador/info"
	"github.com/caos/boom/internal/name"
)

type Ambassador struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.Ambassador
}

func New(logger logging.Logger) *Ambassador {
	return &Ambassador{
		logger: logger,
	}
}

func (a *Ambassador) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.Ambassador.Deploy
}

func (a *Ambassador) Initial() bool {
	return a.spec == nil
}

func (a *Ambassador) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.Ambassador, a.spec)
}

func (a *Ambassador) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	a.spec = toolsetCRDSpec.Ambassador
}

func (a *Ambassador) GetName() name.Application {
	return info.GetName()
}

func (a *Ambassador) GetNamespace() string {
	return info.GetNamespace()
}
