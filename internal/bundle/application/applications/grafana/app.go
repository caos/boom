package grafana

import (
	"path/filepath"
	"reflect"

	"github.com/caos/boom/internal/bundle/application/applications/grafana/config"
	"github.com/caos/boom/internal/kustomize"
	"github.com/caos/boom/internal/templator/helm/chart"

	"github.com/caos/orbiter/logging"

	"github.com/caos/boom/api/v1beta1"
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/helm"
	"github.com/caos/boom/internal/bundle/application/applications/grafanastandalone"
	"github.com/caos/boom/internal/name"
)

const (
	applicationName name.Application = "grafana"
)

func GetName() name.Application {
	return applicationName
}

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
	return applicationName
}

func (g *Grafana) Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.Grafana.Deploy
}

func (g *Grafana) Initial() bool {
	return g.spec == nil
}

func (g *Grafana) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.Grafana, g.spec)
}

func (g *Grafana) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	g.spec = toolsetCRDSpec.Grafana
}

func (g *Grafana) GetNamespace() string {
	return "caos-system"
}

func (g *Grafana) HelmPreApplySteps(spec *v1beta1.ToolsetSpec) ([]interface{}, error) {
	config := config.New(spec.KubeVersion, spec)

	folders := make([]string, 0)
	for _, provider := range config.DashboardProviders {
		folders = append(folders, provider.Folder)
	}

	outs, err := getKustomizeOutput(folders)
	if err != nil {
		return nil, err
	}

	ret := make([]interface{}, len(outs))
	for k, v := range outs {
		ret[k] = v
	}
	return ret, nil
}

func (g *Grafana) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	conf := config.New(toolset.KubeVersion, toolset)
	values := helm.DefaultValues(g.GetImageTags())

	values.KubeTargetVersionOverride = conf.KubeVersion

	if conf.Datasources != nil {
		datasources := make([]*grafanastandalone.Datasource, 0)
		for _, datasource := range conf.Datasources {
			valuesDatasource := &grafanastandalone.Datasource{
				Name:      datasource.Name,
				Type:      datasource.Type,
				URL:       datasource.Url,
				Access:    datasource.Access,
				IsDefault: datasource.IsDefault,
			}
			datasources = append(datasources, valuesDatasource)
		}
		values.Grafana.AdditionalDataSources = datasources
	}

	if conf.DashboardProviders != nil {
		providers := make([]*helm.Provider, 0)
		dashboards := make(map[string]string, 0)
		for _, provider := range conf.DashboardProviders {
			for _, configmap := range provider.ConfigMaps {
				providers = append(providers, getProvider(configmap))
				dashboards[configmap] = configmap
			}
		}
		values.Grafana.DashboardProviders = &helm.DashboardProviders{
			Providers: &helm.Providersyaml{
				APIVersion: 1,
				Providers:  providers,
			},
		}
		values.Grafana.DashboardsConfigMaps = dashboards
	}

	return values
}

func getKustomizeOutput(folders []string) ([]string, error) {
	ret := make([]string, len(folders))
	for n, folder := range folders {

		cmd, err := kustomize.New(folder, false)
		if err != nil {
			return nil, err
		}
		execcmd := cmd.Build()

		out, err := execcmd.Output()
		if err != nil {
			return nil, err
		}
		ret[n] = string(out)
	}
	return ret, nil
}

func getProvider(appName string) *helm.Provider {
	return &helm.Provider{
		Name:            appName,
		Type:            "file",
		DisableDeletion: false,
		Editable:        true,
		Options: map[string]string{
			"path": filepath.Join("/var/lib/grafana/dashboards", appName),
		},
	}
}

func (g *Grafana) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (g *Grafana) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
