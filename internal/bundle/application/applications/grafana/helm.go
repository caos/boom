package grafana

import (
	"path/filepath"

	"github.com/caos/boom/api/v1beta1"
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/config"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/helm"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/info"
	"github.com/caos/boom/internal/bundle/application/applications/grafanastandalone"
	"github.com/caos/boom/internal/kustomize"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/boom/internal/templator/helm/chart"
)

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

	providers := make([]*helm.Provider, 0)
	dashboards := make(map[string]string, 0)
	datasources := make([]*grafanastandalone.Datasource, 0)

	//internal datasources
	if conf.Datasources != nil {
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
	}

	//internal dashboards
	if conf.DashboardProviders != nil {
		for _, provider := range conf.DashboardProviders {
			for _, configmap := range provider.ConfigMaps {
				providers = append(providers, getProvider(configmap))
				dashboards[configmap] = configmap
			}
		}
	}

	if len(providers) > 0 {
		values.Grafana.DashboardProviders = &helm.DashboardProviders{
			Providers: &helm.Providersyaml{
				APIVersion: 1,
				Providers:  providers,
			},
		}
		values.Grafana.DashboardsConfigMaps = dashboards
	}
	if len(datasources) > 0 {
		values.Grafana.AdditionalDataSources = datasources
	}

	if toolset.Grafana.Admin != nil {
		values.Grafana.Admin.ExistingSecret = toolset.Grafana.Admin.ExistingSecret
		values.Grafana.Admin.UserKey = toolset.Grafana.Admin.UserKey
		values.Grafana.Admin.PasswordKey = toolset.Grafana.Admin.PasswordKey
	}

	appLabels := labels.GetApplicationLabels(info.GetName())
	values.Grafana.Labels = appLabels
	values.Grafana.PodLabels = appLabels
	service := &helm.Service{
		Labels: appLabels,
	}
	values.Grafana.Service = service

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
