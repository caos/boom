package templator

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/mntr"
)

type Templator interface {
	Template(interface{}, *v1beta1.ToolsetSpec, func(string, string) error) error
	GetResultsFilePath(name.Application, string, string) string
	CleanUp() error
}

type BaseApplication interface {
	GetName() name.Application
}

type YamlApplication interface {
	BaseApplication
	GetYaml(mntr.Monitor, *v1beta1.ToolsetSpec) interface{}
}

type HelmApplication interface {
	BaseApplication
	GetNamespace() string
	SpecToHelmValues(mntr.Monitor, *v1beta1.ToolsetSpec) interface{}
	GetChartInfo() *chart.Chart
}
