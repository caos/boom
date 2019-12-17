package service

type Config struct {
	Name       string
	Namespace  string
	Labels     map[string]string
	Port       int
	TargetPort int
	Protocol   string
	Selector   map[string]string
}

type Metadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
}

type Spec struct {
	Selector map[string]string `yaml:"selector"`
	Ports    []*Port           `yaml:"ports"`
}
type Port struct {
	Name       string `yaml:"name"`
	Protocol   string `yaml:"protocol"`
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
}

type Service struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   *Metadata `yaml:"metadata"`
	Spec       *Spec     `yaml:"spec"`
}

func New(conf *Config) *Service {
	return &Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: &Metadata{
			Name:      conf.Name,
			Namespace: conf.Namespace,
			Labels:    conf.Labels,
		},
		Spec: &Spec{
			Selector: conf.Selector,
			Ports: []*Port{&Port{
				Name:       "http",
				Protocol:   conf.Protocol,
				Port:       conf.Port,
				TargetPort: conf.TargetPort,
			}},
		},
	}
}
