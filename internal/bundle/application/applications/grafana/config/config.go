package config

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

type Datasource struct {
	Name      string
	Type      string
	Url       string
	Access    string
	IsDefault bool
}

type Provider struct {
	ConfigMaps []string
	Folder     string
}

type Config struct {
	Deploy             bool
	Datasources        []*Datasource
	DashboardProviders []*Provider
	KubeVersion        string
}

func New(kubeVersion string, spec *toolsetsv1beta1.ToolsetSpec) *Config {
	dashboardProviders := make([]*Provider, 0)
	if spec.Grafana.DashboardProviders != nil {
		for _, provider := range spec.Grafana.DashboardProviders {
			confProvider := &Provider{
				ConfigMaps: provider.ConfigMaps,
				Folder:     provider.Folder,
			}
			dashboardProviders = append(dashboardProviders, confProvider)
		}
	}

	datasources := make([]*Datasource, 0)
	if spec.Grafana.Datasources != nil {
		for _, datasource := range spec.Grafana.Datasources {
			confDatasource := &Datasource{
				Name:      datasource.Name,
				Type:      datasource.Type,
				Url:       datasource.Url,
				Access:    datasource.Access,
				IsDefault: datasource.IsDefault,
			}
			datasources = append(datasources, confDatasource)
		}
	}

	conf := &Config{
		Deploy:             spec.Grafana.Deploy,
		DashboardProviders: dashboardProviders,
		Datasources:        datasources,
		KubeVersion:        kubeVersion,
	}

	providers := getGrafanaDashboards("../../dashboards", spec)

	for _, provider := range providers {
		conf.AddDashboardProvider(provider)
	}

	datasourceProm := strings.Join([]string{"http://prometheus-operated.caos-system:9090"}, "")
	conf.AddDatasourceURL("prometheus", "prometheus", datasourceProm)

	datasourceLoki := strings.Join([]string{"http://loki.caos-system:3100"}, "")
	conf.AddDatasourceURL("loki", "loki", datasourceLoki)

	return conf
}

func (c *Config) AddDashboardProvider(provider *Provider) {
	c.DashboardProviders = append(c.DashboardProviders, provider)
}

func (c *Config) AddDatasourceURL(name, datatype, url string) {
	c.Datasources = append(c.Datasources, &Datasource{
		Name:   name,
		Type:   datatype,
		Url:    url,
		Access: "proxy",
	})
}
