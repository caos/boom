package bundles

import (
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
