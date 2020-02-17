package postapply

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/postapply/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/logging"
)

type PostApply struct {
	logger logging.Logger
}

func New(logger logging.Logger) *PostApply {
	return &PostApply{
		logger: logger,
	}
}

func (p *PostApply) GetName() name.Application {
	return info.GetName()
}

func (p *PostApply) Deploy(toolsetCRDSpec *v1beta1.ToolsetSpec) bool {
	if toolsetCRDSpec.PostApply == nil {
		return false
	}

	return toolsetCRDSpec.PostApply.Deploy
}
