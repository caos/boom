package v1beta1

type Argocd struct {
	Deploy       bool                `json:"deploy,omitempty"`
	CustomImage  *ArgocdCustomImage  `json:"customImage,omitempty" yaml:"customImage,omitempty"`
	Network      *Network            `json:"network,omitempty" yaml:"network,omitempty"`
	Auth         *ArgocdAuth         `json:"auth,omitempty" yaml:"auth,omitempty"`
	Repositories []*ArgocdRepository `json:"repositories,omitempty" yaml:"repositories,omitempty"`
}

type ArgocdRepository struct {
	URL               string        `json:"url,omitempty" yaml:"url,omitempty"`
	UsernameSecret    *ArgocdSecret `json:"usernameSecret,omitempty" yaml:"usernameSecret,omitempty"`
	PasswordSecret    *ArgocdSecret `json:"passwordSecret,omitempty" yaml:"passwordSecret,omitempty"`
	CertificateSecret *ArgocdSecret `json:"certificateSecret,omitempty" yaml:"certificateSecret,omitempty"`
}

type ArgocdSecret struct {
	Name string `json:"name" yaml:"name"`
	Key  string `json:"key" yaml:"key"`
}

type ArgocdCustomImage struct {
	Enabled         bool                 `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	ImagePullSecret string               `json:"imagePullSecret,omitempty" yaml:"imagePullSecret,omitempty"`
	GopassGPGKey    string               `json:"gopassPGPKey,omitempty" yaml:"gopassPGPKey,omitempty"`
	GopassSSHKey    string               `json:"gopassSSHKey,omitempty" yaml:"gopassSSHKey,omitempty"`
	GopassStores    []*ArgocdGopassStore `json:"gopassStores,omitempty" yaml:"gopassStores,omitempty"`
}

type ArgocdGopassStore struct {
	Directory string `json:"directory,omitempty" yaml:"directory,omitempty"`
	StoreName string `json:"storeName,omitempty" yaml:"storeName,omitempty"`
}

type ArgocdAuth struct {
	OIDC            *ArgocdOIDC            `json:"oidc,omitempty" yaml:"oidc,omitempty"`
	GithubConnector *ArgocdGithubConnector `json:"github,omitempty" yaml:"github,omitempty"`
	GitlabConnector *ArgocdGitlabConnector `json:"gitlab,omitempty" yaml:"gitlab,omitempty"`
	GoogleConnector *ArgocdGoogleConnector `json:"google,omitempty" yaml:"google,omitempty"`
}

type ArgocdOIDC struct {
	Name                   string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Issuer                 string                 `json:"issuer,omitempty" yaml:"issuer,omitempty"`
	SecretName             string                 `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	ClientIDKey            string                 `json:"clientIDKey,omitempty" yaml:"clientIDKey,omitempty"`
	ClientSecretKey        string                 `json:"clientSecretKey,omitempty" yaml:"clientSecret,omitempty"`
	RequestedScopes        []string               `json:"requestedScopes,omitempty" yaml:"requestedScopes,omitempty"`
	RequestedIDTokenClaims map[string]ArgocdClaim `json:"requestedIDTokenClaims,omitempty" yaml:"requestedIDTokenClaims,omitempty"`
}

type ArgocdClaim struct {
	Essential bool     `json:"essential,omitempty" yaml:"essential,omitempty"`
	Values    []string `json:"values,omitempty" yaml:"values,omitempty"`
}

type ArgocdGithubConnector struct {
	ID     string              `json:"id,omitempty" yaml:"id,omitempty"`
	Name   string              `json:"name,omitempty" yaml:"name,omitempty"`
	Config *ArgocdGithubConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

type ArgocdGithubConfig struct {
	SecretName      string             `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	ClientIDKey     string             `json:"clientIDKey,omitempty" yaml:"clientIDKey,omitempty"`
	ClientSecretKey string             `json:"clientSecretKey,omitempty" yaml:"clientSecretKey,omitempty"`
	Orgs            []*ArgocdGithubOrg `json:"orgs,omitempty" yaml:"orgs,omitempty"`
	LoadAllGroups   bool               `json:"loadAllGroups,omitempty" yaml:"loadAllGroups,omitempty"`
	TeamNameField   string             `json:"teamNameField,omitempty" yaml:"teamNameField,omitempty"`
	UseLoginAsID    bool               `json:"useLoginAsID,omitempty" yaml:"useLoginAsID,omitempty"`
}

type ArgocdGithubOrg struct {
	Name  string   `json:"name,omitempty" yaml:"name,omitempty"`
	Teams []string `json:"teams,omitempty" yaml:"teams,omitempty"`
}

type ArgocdGitlabConnector struct {
	ID     string              `json:"id,omitempty" yaml:"id,omitempty"`
	Name   string              `json:"name,omitempty" yaml:"name,omitempty"`
	Config *ArgocdGitlabConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

type ArgocdGitlabConfig struct {
	SecretName      string   `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	ClientIDKey     string   `json:"clientIDKey,omitempty" yaml:"clientIDKey,omitempty"`
	ClientSecretKey string   `json:"clientSecretKey,omitempty" yaml:"clientSecretKey,omitempty"`
	BaseURL         string   `json:"baseURL,omitempty" yaml:"baseURL,omitempty"`
	Groups          []string `json:"groups,omitempty" yaml:"groups,omitempty"`
	UseLoginAsID    bool     `json:"useLoginAsID,omitempty" yaml:"useLoginAsID,omitempty"`
}

type ArgocdGoogleConnector struct {
	ID     string              `json:"id,omitempty" yaml:"id,omitempty"`
	Name   string              `json:"name,omitempty" yaml:"name,omitempty"`
	Config *ArgocdGoogleConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

type ArgocdGoogleConfig struct {
	SecretName             string   `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	ClientIDKey            string   `json:"clientIDKey,omitempty" yaml:"clientIDKey,omitempty"`
	ClientSecretKey        string   `json:"clientSecretKey,omitempty" yaml:"clientSecretKey,omitempty"`
	HostedDomains          []string `json:"hostedDomains,omitempty" yaml:"hostedDomains,omitempty"`
	Groups                 []string `json:"groups,omitempty" yaml:"groups,omitempty"`
	ServiceAccountJSONKey  string   `json:"serviceAccountJSONKey,omitempty" yaml:"serviceAccountJSONKey,omitempty"`
	ServiceAccountFilePath string   `json:"serviceAccountFilePath,omitempty" yaml:"serviceAccountFilePath,omitempty"`
	AdminEmail             string   `json:"adminEmail,omitempty" yaml:"adminEmail,omitempty"`
}
