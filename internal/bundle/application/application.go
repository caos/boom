package application

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/ambassador"
	ambassadorinfo "github.com/caos/boom/internal/bundle/application/applications/ambassador/info"
	"github.com/caos/boom/internal/bundle/application/applications/argocd"
	argocdinfo "github.com/caos/boom/internal/bundle/application/applications/argocd/info"
	"github.com/caos/boom/internal/bundle/application/applications/grafana"
	grafanainfo "github.com/caos/boom/internal/bundle/application/applications/grafana/info"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics"
	kubestatemetricsinfo "github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/info"
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator"
	loggingoperatorinfo "github.com/caos/boom/internal/bundle/application/applications/loggingoperator/info"
	"github.com/caos/boom/internal/bundle/application/applications/loki"
	lokiinfo "github.com/caos/boom/internal/bundle/application/applications/loki/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus"
	prometheusinfo "github.com/caos/boom/internal/bundle/application/applications/prometheus/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter"
	prometheusnodeexporterinfo "github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator"
	prometheusoperatorinfo "github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheussystemdexporter"
	prometheussystemdexporterinfo "github.com/caos/boom/internal/bundle/application/applications/prometheussystemdexporter/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/mntr"
)

type Application interface {
	Deploy(*v1beta1.ToolsetSpec) bool
	GetName() name.Application
}

type HelmApplication interface {
	Application
	GetNamespace() string
	GetChartInfo() *chart.Chart
	GetImageTags() map[string]string
	SpecToHelmValues(mntr.Monitor, *v1beta1.ToolsetSpec) interface{}
}

type YAMLApplication interface {
	Application
	GetYaml(mntr.Monitor, *v1beta1.ToolsetSpec) interface{}
}

func New(monitor mntr.Monitor, appName name.Application, orb string) Application {
	switch appName {
	case ambassadorinfo.GetName():
		return ambassador.New(monitor)
	case argocdinfo.GetName():
		return argocd.New(monitor)
	case grafanainfo.GetName():
		return grafana.New(monitor)
	case kubestatemetricsinfo.GetName():
		return kubestatemetrics.New(monitor)
	case prometheusoperatorinfo.GetName():
		return prometheusoperator.New(monitor)
	case loggingoperatorinfo.GetName():
		return loggingoperator.New(monitor)
	case prometheusnodeexporterinfo.GetName():
		return prometheusnodeexporter.New(monitor)
	case prometheussystemdexporterinfo.GetName():
		return prometheussystemdexporter.New()
	case prometheusinfo.GetName():
		return prometheus.New(monitor, orb)
	case lokiinfo.GetName():
		return loki.New(monitor)
	}

	return nil
}

func GetOrderNumber(appName name.Application) int {
	switch appName {
	case prometheusinfo.GetName():
		return prometheusinfo.GetOrderNumber()
	case lokiinfo.GetName():
		return lokiinfo.GetOrderNumber()
	}

	return 1
}
