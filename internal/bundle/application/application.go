package application

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/ambassador"
	"github.com/caos/boom/internal/bundle/application/applications/argocd"
	"github.com/caos/boom/internal/bundle/application/applications/certmanager"
	"github.com/caos/boom/internal/bundle/application/applications/grafana"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics"
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/logging"
)

type Application interface {
	Initial() bool
	Changed(*v1beta1.ToolsetSpec) bool
	Deploy(*v1beta1.ToolsetSpec) bool
	SetAppliedSpec(*v1beta1.ToolsetSpec)
	GetName() name.Application
	GetNamespace() string
}

type HelmApplication interface {
	Application
	GetChartInfo() *chart.Chart
	GetImageTags() map[string]string
	SpecToHelmValues(spec *v1beta1.ToolsetSpec) interface{}
}

type YAMLApplication interface {
	Application
	GetYaml() interface{}
}

func New(logger logging.Logger, appName name.Application) Application {
	switch appName {
	case ambassador.GetName():
		return ambassador.New(logger)
	case argocd.GetName():
		return argocd.New(logger)
	case certmanager.GetName():
		return certmanager.New(logger)
	case grafana.GetName():
		return grafana.New(logger)
	case kubestatemetrics.GetName():
		return kubestatemetrics.New(logger)
	case prometheusoperator.GetName():
		return prometheusoperator.New(logger)
	case loggingoperator.GetName():
		return loggingoperator.New(logger)
	case prometheusnodeexporter.GetName():
		return prometheusnodeexporter.New(logger)
	case prometheus.GetName():
		return prometheus.New(logger)
	}

	return nil
}
