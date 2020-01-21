package v1beta1

type Grafana struct {
	Deploy             bool          `json:"deploy,omitempty"`
	Admin              *Admin        `json:"admin,omitempty"`
	Datasources        []*Datasource `json:"datasources,omitempty"`
	DashboardProviders []*Provider   `json:"dashboardproviders,omitempty"`
}
type Admin struct {
	ExistingSecret string `json:"existingSecret,omitempty"`
	UserKey        string `json:"userKey,omitempty"`
	PasswordKey    string `json:"passwordKey,omitempty"`
}
type Datasource struct {
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Url       string `json:"url,omitempty"`
	Access    string `json:"access,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

type Provider struct {
	ConfigMaps []string `json:"configMap,omitempty"`
	Folder     string   `json:"folder,omitempty"`
}
