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
	"github.com/caos/boom/internal/bundle/application/applications/postapply"
	postapplyinfo "github.com/caos/boom/internal/bundle/application/applications/postapply/info"
	"github.com/caos/boom/internal/bundle/application/applications/preapply"
	preapplyinfo "github.com/caos/boom/internal/bundle/application/applications/preapply/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus"
	prometheusinfo "github.com/caos/boom/internal/bundle/application/applications/prometheus/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter"
	prometheusnodeexporterinfo "github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator"
	prometheusoperatorinfo "github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/info"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/logging"
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
	SpecToHelmValues(logging.Logger, *v1beta1.ToolsetSpec) interface{}
}

type YAMLApplication interface {
	Application
	GetYaml(logging.Logger, *v1beta1.ToolsetSpec) interface{}
}

func New(logger logging.Logger, appName name.Application) Application {
	switch appName {
	case ambassadorinfo.GetName():
		return ambassador.New(logger)
	case argocdinfo.GetName():
		return argocd.New(logger)
	case grafanainfo.GetName():
		return grafana.New(logger)
	case kubestatemetricsinfo.GetName():
		return kubestatemetrics.New(logger)
	case prometheusoperatorinfo.GetName():
		return prometheusoperator.New(logger)
	case loggingoperatorinfo.GetName():
		return loggingoperator.New(logger)
	case prometheusnodeexporterinfo.GetName():
		return prometheusnodeexporter.New(logger)
	case prometheusinfo.GetName():
		return prometheus.New(logger)
	case lokiinfo.GetName():
		return loki.New(logger)
	case preapplyinfo.GetName():
		return preapply.New(logger)
	case postapplyinfo.GetName():
		return postapply.New(logger)
	}

	return nil
}

func GetOrderNumber(appName name.Application) int {
	switch appName {
	case prometheusinfo.GetName():
		return prometheusinfo.GetOrderNumber()
	case lokiinfo.GetName():
		return lokiinfo.GetOrderNumber()
	case preapplyinfo.GetName():
		return preapplyinfo.GetOrderNumber()
	case postapplyinfo.GetName():
		return postapplyinfo.GetOrderNumber()
	}

	return 1
}
