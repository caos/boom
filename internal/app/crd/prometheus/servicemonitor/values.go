package servicemonitor

type Selector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

type NamespaceSelector struct {
	Any        bool     `yaml:"any"`
	MatchNames []string `yaml:"matchNames"`
}

type Endpoint struct {
	Port            string     `yaml:"port,omitempty"`
	TargetPort      string     `yaml:"targetPort,omitempty"`
	BearerTokenFile string     `yaml:"bearerTokenFile,omitempty"`
	Interval        string     `yaml:"interval,omitempty"`
	Path            string     `yaml:"path,omitempty"`
	Scheme          string     `yaml:"scheme,omitempty"`
	TLSConfig       *TLSConfig `yaml:"tlsConfig,omitempty"`
}

type TLSConfig struct {
	CaFile             string `yaml:"caFile"`
	CertFile           string `yaml:"certFile"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
	KeyFile            string `yaml:"keyFile"`
	ServerName         string `yaml:"serverName"`
}

type Values struct {
	Name              string             `yaml:"name"`
	AdditionalLabels  map[string]string  `yaml:"additionalLabels"`
	JobLabel          string             `yaml:"jobLabel"`
	TargetLabels      string             `yaml:"targetLabels"`
	Selector          *Selector          `yaml:"selector"`
	NamespaceSelector *NamespaceSelector `yaml:"namespaceSelector"`
	Endpoints         []*Endpoint        `yaml:"endpoints"`
}
