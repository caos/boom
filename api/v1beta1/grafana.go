package v1beta1

type Grafana struct {
	Deploy             bool          `json:"deploy,omitempty" yaml:"deploy,omitempty"`
	Admin              *Admin        `json:"admin,omitempty" yaml:"admin,omitempty"`
	Datasources        []*Datasource `json:"datasources,omitempty" yaml:"datasources,omitempty"`
	DashboardProviders []*Provider   `json:"dashboardproviders,omitempty" yaml:"dashboardproviders,omitempty"`
	Storage            *StorageSpec  `json:"storage,omitempty" yaml:"storage,omitempty"`
	Auth               *Auth         `json:"auth,omitempty" yaml:"auth,omitempty"`
}

type Admin struct {
	ExistingSecret string `json:"existingSecret,omitempty" yaml:"existingSecret,omitempty"`
	UserKey        string `json:"userKey,omitempty" yaml:"userKey,omitempty"`
	PasswordKey    string `json:"passwordKey,omitempty" yaml:"passwordKey,omitempty"`
}

type Datasource struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Type      string `json:"type,omitempty" yaml:"type,omitempty"`
	Url       string `json:"url,omitempty" yaml:"url,omitempty"`
	Access    string `json:"access,omitempty" yaml:"access,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty" yaml:"isDefault,omitempty"`
}

type Provider struct {
	ConfigMaps []string `json:"configMaps,omitempty" yaml:"configMaps,omitempty"`
	Folder     string   `json:"folder,omitempty" yaml:"folder,omitempty"`
}

type Auth struct {
	Google       *GoogleAuth   `json:"google,omitempty" yaml:"google,omitempty"`
	Github       *GithubAuth   `json:"github,omitempty" yaml:"github,omitempty"`
	Gitlab       *GitlabAuth   `json:"gitlab,omitempty" yaml:"gitlab,omitempty"`
	GenericOAuth *GenericOAuth `json:"genericOAuth,omitempty" yaml:"genericOAuth,omitempty"`
}

type GoogleAuth struct {
	SecretName     string   `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	ClientIDKey    string   `json:"clientIDKey,omitempty" yaml:"clientIDKey,omitempty"`
	ClientSecret   string   `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	AllowedDomains []string `json:"allowedDomains,omitempty" yaml:"allowedDomains,omitempty"`
}

type GithubAuth struct {
	SecretName           string   `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	ClientIDKey          string   `json:"clientIDKey,omitempty" yaml:"clientIDKey,omitempty"`
	ClientSecret         string   `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	AllowedOrganizations []string `json:"allowedOrganizations,omitempty" yaml:"allowedOrganizations,omitempty"`
	TeamIDs              []string `json:"teamIDs,omitempty" yaml:"teamIDs,omitempty"`
}

type GitlabAuth struct {
	SecretName    string   `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	ClientIDKey   string   `json:"clientIDKey,omitempty" yaml:"clientIDKey,omitempty"`
	ClientSecret  string   `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	AllowedGroups []string `json:"allowedGroups,omitempty" yaml:"allowedGroups,omitempty"`
}

type GenericOAuth struct {
	SecretName     string   `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	ClientIDKey    string   `json:"clientIDKey,omitempty" yaml:"clientIDKey,omitempty"`
	ClientSecret   string   `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	Scopes         []string `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	AuthURL        string   `json:"authURL,omitempty" yaml:"authURL,omitempty"`
	TokenURL       string   `json:"tokenURL,omitempty" yaml:"tokenURL,omitempty"`
	APIURL         string   `json:"apiURL,omitempty" yaml:"apiURL,omitempty"`
	AllowedDomains []string `json:"allowedDomains,omitempty" yaml:"allowedDomains,omitempty"`
}
