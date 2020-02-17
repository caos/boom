package preapply

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/preapply/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/logging"
)

type PreApply struct {
	logger logging.Logger
}

func New(logger logging.Logger) *PreApply {
	return &PreApply{
		logger: logger,
	}
}

func (p *PreApply) GetName() name.Application {
	return info.GetName()
}

func (p *PreApply) Deploy(toolsetCRDSpec *v1beta1.ToolsetSpec) bool {
	if toolsetCRDSpec.PreApply == nil {
		return false
	}

	return toolsetCRDSpec.PreApply.Deploy
}
