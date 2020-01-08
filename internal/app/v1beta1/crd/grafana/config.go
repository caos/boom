package grafana

import toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"

type ConfigDatasource struct {
	Name      string
	Type      string
	Url       string
	Access    string
	IsDefault bool
}

type ConfigProvider struct {
	ConfigMaps []string
	Folder     string
}

type Config struct {
	Deploy             bool
	Prefix             string
	Namespace          string
	Datasources        []*ConfigDatasource
	DashboardProviders []*ConfigProvider
	KubeVersion        string
}

func NewConfig(kubeVersion string, crd *toolsetsv1beta1.Grafana) *Config {
	dashboardProviders := make([]*ConfigProvider, 0)
	if crd.DashboardProviders != nil {
		for _, provider := range crd.DashboardProviders {
			confProvider := &ConfigProvider{
				ConfigMaps: provider.ConfigMaps,
				Folder:     provider.Folder,
			}
			dashboardProviders = append(dashboardProviders, confProvider)
		}
	}
	datasources := make([]*ConfigDatasource, 0)
	if crd.Datasources != nil {
		for _, datasource := range crd.Datasources {
			confDatasource := &ConfigDatasource{
				Name:      datasource.Name,
				Type:      datasource.Type,
				Url:       datasource.Url,
				Access:    datasource.Access,
				IsDefault: datasource.IsDefault,
			}
			datasources = append(datasources, confDatasource)
		}
	}

	return &Config{
		Deploy:             crd.Deploy,
		Prefix:             crd.Prefix,
		Namespace:          crd.Namespace,
		DashboardProviders: dashboardProviders,
		Datasources:        datasources,
		KubeVersion:        kubeVersion,
	}
}

func (c *Config) AddDashboardProvider(provider *ConfigProvider) {
	c.DashboardProviders = append(c.DashboardProviders, provider)
}

func (c *Config) AddDatasourceURL(name, datatype, url string) {
	c.Datasources = append(c.Datasources, &ConfigDatasource{
		Name:   name,
		Type:   datatype,
		Url:    url,
		Access: "proxy",
	})
}
