package v1beta1

type Grafana struct {
	Deploy             bool          `json:"deploy,omitempty"`
	Admin              *Admin        `json:"admin,omitempty"`
	Datasources        []*Datasource `json:"datasources,omitempty"`
	DashboardProviders []*Provider   `json:"dashboardproviders,omitempty"`
	Storage            *StorageSpec  `json:"storage,omitempty"`
}
type Admin struct {
	ExistingSecret string `json:"existingSecret,omitempty" yaml:"existingSecret,omitempty"`
	UserKey        string `json:"userKey,omitempty" yaml:"userKey,omitempty"`
	PasswordKey    string `json:"passwordKey,omitempty" yaml:"passwordKey,omitempty"`
}

type Datasource struct {
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Url       string `json:"url,omitempty"`
	Access    string `json:"access,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty" yaml:"isDefault,omitempty"`
}

type Provider struct {
	ConfigMaps []string `json:"configMaps,omitempty" yaml:"configMaps,omitempty"`
	Folder     string   `json:"folder,omitempty"`
}
