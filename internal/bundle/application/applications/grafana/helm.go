package grafana

import (
	"path/filepath"

	"github.com/caos/boom/api/v1beta1"
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/auth"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/config"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/helm"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/info"
	"github.com/caos/boom/internal/bundle/application/applications/grafanastandalone"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/kustomize"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/logging"
)

func (g *Grafana) HelmPreApplySteps(logger logging.Logger, spec *v1beta1.ToolsetSpec) ([]interface{}, error) {
	config := config.New(spec)

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

func (g *Grafana) SpecToHelmValues(logger logging.Logger, toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	version, err := kubectl.NewVersion().GetKubeVersion(logger)
	if err != nil {
		return nil
	}

	conf := config.New(toolset)
	values := helm.DefaultValues(g.GetImageTags())

	values.KubeTargetVersionOverride = version

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
	if toolset.Grafana.Storage != nil {
		values.Grafana.Persistence.Enabled = true
		values.Grafana.Persistence.Size = toolset.Grafana.Storage.Size
		values.Grafana.Persistence.StorageClassName = toolset.Grafana.Storage.StorageClass

		if toolset.Grafana.Storage.AccessModes != nil {
			values.Grafana.Persistence.AccessModes = toolset.Grafana.Storage.AccessModes
		}
	}

	if toolset.Grafana.Network != nil && toolset.Grafana.Network.Domain != "" {
		values.Grafana.Env["GF_SERVER_DOMAIN"] = toolset.Grafana.Network.Domain

		if toolset.Grafana.Auth != nil {
			if toolset.Grafana.Auth.Google != nil {
				google, err := auth.GetGoogleAuthConfig(toolset.Grafana.Auth.Google)
				if err == nil {
					values.Grafana.Ini.AuthGoogle = google
				}
			}

			if toolset.Grafana.Auth.Github != nil {
				github, err := auth.GetGithubAuthConfig(toolset.Grafana.Auth.Github)
				if err == nil {
					values.Grafana.Ini.AuthGithub = github
				}
			}

			if toolset.Grafana.Auth.Gitlab != nil {
				gitlab, err := auth.GetGitlabAuthConfig(toolset.Grafana.Auth.Gitlab)
				if err == nil {
					values.Grafana.Ini.AuthGitlab = gitlab
				}
			}

			if toolset.Grafana.Auth.GenericOAuth != nil {
				generic, err := auth.GetGenericOAuthConfig(toolset.Grafana.Auth.GenericOAuth)
				if err == nil {
					values.Grafana.Ini.AuthGeneric = generic
				}
			}
		}
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
