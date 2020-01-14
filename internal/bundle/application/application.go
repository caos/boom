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
	"github.com/caos/boom/internal/templator/helm"
	"github.com/caos/orbiter/logging"
)

type Application interface {
	Initial() bool
	Changed(toolsetCRDSpec *v1beta1.ToolsetSpec) bool
	SetAppliedSpec(toolsetCRDSpec *v1beta1.ToolsetSpec)
}

type HelmApplication interface {
	Application
	GetImageTags() map[string]string
	helm.Templator
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

func Deploy(appName name.Application, toolsetCRDSpec *v1beta1.ToolsetSpec) bool {
	switch appName {
	case ambassador.GetName():
		return ambassador.Deploy(toolsetCRDSpec)
	case argocd.GetName():
		return argocd.Deploy(toolsetCRDSpec)
	case certmanager.GetName():
		return certmanager.Deploy(toolsetCRDSpec)
	case grafana.GetName():
		return grafana.Deploy(toolsetCRDSpec)
	case kubestatemetrics.GetName():
		return kubestatemetrics.Deploy(toolsetCRDSpec)
	case prometheusoperator.GetName():
		return prometheusoperator.Deploy(toolsetCRDSpec)
	case loggingoperator.GetName():
		return loggingoperator.Deploy(toolsetCRDSpec)
	case prometheusnodeexporter.GetName():
		return prometheusnodeexporter.Deploy(toolsetCRDSpec)
	case prometheus.GetName():
		return prometheus.Deploy(toolsetCRDSpec)
	}

	return false
}
