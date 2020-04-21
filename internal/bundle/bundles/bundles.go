package bundles

import (
	ambassadorinfo "github.com/caos/boom/internal/bundle/application/applications/ambassador/info"
	argocdinfo "github.com/caos/boom/internal/bundle/application/applications/argocd/info"
	grafanainfo "github.com/caos/boom/internal/bundle/application/applications/grafana/info"
	kubestatemetricsinfo "github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics/info"
	loggingoperatorinfo "github.com/caos/boom/internal/bundle/application/applications/loggingoperator/info"
	lokiinfo "github.com/caos/boom/internal/bundle/application/applications/loki/info"
	prometheusinfo "github.com/caos/boom/internal/bundle/application/applications/prometheus/info"
	prometheusnodeexporterinfo "github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter/info"
	prometheusoperatorinfo "github.com/caos/boom/internal/bundle/application/applications/prometheusoperator/info"
	prometheussystemdexporterinfo "github.com/caos/boom/internal/bundle/application/applications/prometheussystemdexporter/info"
	"github.com/caos/boom/internal/name"
)

const (
	Caos  name.Bundle = "caos"
	Empty name.Bundle = "empty"
)

func GetAll() []name.Application {
	apps := make([]name.Application, 0)
	apps = append(apps, GetCaos()...)
	return apps
}

func Get(bundle name.Bundle) []name.Application {
	switch bundle {
	case Caos:
		return GetCaos()
	case Empty:
		return make([]name.Application, 0)
	}

	return nil
}

func GetCaos() []name.Application {

	apps := make([]name.Application, 0)
	apps = append(apps, ambassadorinfo.GetName())
	apps = append(apps, argocdinfo.GetName())
	apps = append(apps, prometheusoperatorinfo.GetName())
	apps = append(apps, kubestatemetricsinfo.GetName())
	apps = append(apps, prometheusnodeexporterinfo.GetName())
	apps = append(apps, prometheussystemdexporterinfo.GetName())
	apps = append(apps, grafanainfo.GetName())
	apps = append(apps, prometheusinfo.GetName())
	apps = append(apps, loggingoperatorinfo.GetName())
	apps = append(apps, lokiinfo.GetName())

	return apps
}
