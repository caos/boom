package v1beta1

type Grafana struct {
	Deploy      bool          `json:"deploy,omitempty"`
	Prefix      string        `json:"prefix,omitempty"`
	Namespace   string        `json:"namespace,omitempty"`
	Admin       *Admin        `json:"admin,omitempty"`
	Datasources []*Datasource `json:"datasources,omitempty"`
	Dashboards  []*Dashboard  `json:"dashboards,omitempty"`
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

type Dashboard struct {
	ConfigMap string           `json:"configMap,omitempty"`
	FileNames []*DashboardFile `json:"files,omitempty"`
}
type DashboardFile struct {
	Name     string `json:"name,omitempty"`
	FileName string `json:"filename,omitempty"`
}
