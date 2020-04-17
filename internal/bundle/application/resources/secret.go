package resources

type SecretConfig struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Data      map[string]string
}

type Secret struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   *Metadata         `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
	Type       string            `yaml:"type"`
}

func NewSecret(conf *SecretConfig) *Secret {
	return &Secret{
		APIVersion: "v1",
		Kind:       "Secret",
		Metadata: &Metadata{
			Name:      conf.Name,
			Namespace: conf.Namespace,
			Labels:    conf.Labels,
		},
		Data: conf.Data,
	}
}
