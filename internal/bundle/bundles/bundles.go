package bundles

import (
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/bundle/application/applications/ambassador"
	"github.com/caos/boom/internal/bundle/application/applications/argocd"
	"github.com/caos/boom/internal/bundle/application/applications/certmanager"
	"github.com/caos/boom/internal/bundle/application/applications/grafana"
	"github.com/caos/boom/internal/bundle/application/applications/kubestatemetrics"
	"github.com/caos/boom/internal/bundle/application/applications/loggingoperator"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusnodeexporter"
	"github.com/caos/boom/internal/bundle/application/applications/prometheusoperator"
)

func GetAll() []name.Application {
	apps := make([]name.Application, 0)
	apps = append(apps, GetBasisset()...)
	return apps
}

func GetBasisset() []name.Application {

	apps := make([]name.Application, 0)
	apps = append(apps, ambassador.GetName())
	apps = append(apps, argocd.GetName())
	apps = append(apps, prometheusoperator.GetName())
	apps = append(apps, certmanager.GetName())
	apps = append(apps, kubestatemetrics.GetName())
	apps = append(apps, prometheusnodeexporter.GetName())
	apps = append(apps, grafana.GetName())
	apps = append(apps, prometheus.GetName())
	apps = append(apps, loggingoperator.GetName())

	return apps
}
