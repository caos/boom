package v1beta1

type CertManager struct {
	Deploy    bool   `json:"deploy,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	// ClusterIssuer []*ClusterIssuer `json:"clusterIssuers,omitempty"`
}

type ClusterIssuer struct {
	Name string `json:"name,omitempty"`
	Acme *Acme  `json:"acme,omitempty"`
}

type Acme struct {
	Email               string     `json:"email,omitempty"`
	Server              string     `json:"server,omitempty"`
	PrivateKeySecretRef *SecretRef `json:"privateKeySecretRef,omitempty"`
	Http01              *Http01    `json:"http01,omitempty"`
}

type Http01 struct {
}

type SecretRef struct {
	Name string `json:"name,omitempty"`
}
