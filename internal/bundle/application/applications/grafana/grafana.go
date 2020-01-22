package grafana

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/info"
	"github.com/caos/boom/internal/name"
)

type Grafana struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.Grafana
}

func New(logger logging.Logger) *Grafana {
	return &Grafana{
		logger: logger,
	}
}
func (g *Grafana) GetName() name.Application {
	return info.GetName()
}

func (g *Grafana) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	// workaround as grafana always deploys new pods even when the spec of the deployment is not changed
	// due to the fact that kubernetes has an internal mapping from extensions/v1beta1 to apps/v1 in old k8s versions
	if g.Changed(toolsetCRDSpec) {
		return toolsetCRDSpec.Grafana.Deploy
	}
	return false
}

func (g *Grafana) Initial() bool {
	return g.spec == nil
}

func (g *Grafana) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	if g.spec == nil {
		return true
	}
	return !reflect.DeepEqual(toolsetCRDSpec.Grafana, g.spec)
}

func (g *Grafana) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	g.spec = toolsetCRDSpec.Grafana
}

func (g *Grafana) GetNamespace() string {
	return info.GetNamespace()
}
