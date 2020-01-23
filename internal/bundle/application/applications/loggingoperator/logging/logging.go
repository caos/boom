package logging

type Config struct {
	Name             string
	Namespace        string
	ControlNamespace string
}
type Requests struct {
	Storage string `yaml:"storage,omitempty"`
}
type Resources struct {
	Requests *Requests `yaml:"requests,omitempty"`
}
type FluentdPvcSpec struct {
	AccessModes      []string   `yaml:"accessModes,omitempty"`
	Resources        *Resources `yaml:"resources,omitempty"`
	StorageClassName string     `yaml:"storageClassName,omitempty"`
}
type Fluentd struct {
	Metrics        *Metrics        `yaml:"metrics,omitempty"`
	FluentdPvcSpec *FluentdPvcSpec `yaml:"fluentdPvcSpec,omitempty"`
}
type Metrics struct {
	Port int `yaml:"port"`
}
type Image struct {
	PullPolicy string `yaml:"pullPolicy"`
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

type Fluentbit struct {
	Metrics *Metrics `yaml:"metrics,omitempty"`
	Image   *Image   `yaml:"image,omitempty"`
}
type Spec struct {
	Fluentd          *Fluentd   `yaml:"fluentd"`
	Fluentbit        *Fluentbit `yaml:"fluentbit"`
	ControlNamespace string     `yaml:"controlNamespace"`
}
type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}
type Logging struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   *Metadata `yaml:"metadata"`
	Spec       *Spec     `yaml:"spec"`
}

func New(conf *Config) *Logging {
	return &Logging{
		APIVersion: "logging.banzaicloud.io/v1beta1",
		Kind:       "Logging",
		Metadata: &Metadata{
			Name:      conf.Name,
			Namespace: conf.Namespace,
		},
		Spec: &Spec{
			ControlNamespace: conf.ControlNamespace,
			Fluentd: &Fluentd{
				Metrics: &Metrics{
					Port: 8080,
				},
			},
			Fluentbit: &Fluentbit{
				Metrics: &Metrics{
					Port: 8080,
				},
				Image: &Image{
					Repository: "fluent/fluent-bit",
					PullPolicy: "IfNotPresent",
					Tag:        "1.3.6-debug",
				},
			},
		},
	}
}
