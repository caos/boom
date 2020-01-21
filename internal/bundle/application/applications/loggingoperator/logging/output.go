package logging

type ConfigOutput struct {
	Name      string
	Namespace string
	URL       string
}

type Buffer struct {
	Timekey       string `yaml:"timekey"`
	TimekeyWait   string `yaml:"timekey_wait"`
	TimekeyUseUtc bool   `yaml:"timekey_use_utc"`
}
type Loki struct {
	URL                       string  `yaml:"url"`
	ConfigureKubernetesLabels bool    `yaml:"configure_kubernetes_labels"`
	Buffer                    *Buffer `yaml:"buffer"`
}

type OutputSpec struct {
	Loki *Loki `yaml:"loki"`
}

type Output struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   *Metadata   `yaml:"metadata"`
	Spec       *OutputSpec `yaml:"spec"`
}

func NewOutput(conf *ConfigOutput) *Output {
	return &Output{
		APIVersion: "logging.banzaicloud.io/v1beta1",
		Kind:       "Output",
		Metadata: &Metadata{
			Name:      conf.Name,
			Namespace: conf.Namespace,
		},
		Spec: &OutputSpec{
			Loki: &Loki{
				URL:                       conf.URL,
				ConfigureKubernetesLabels: true,
				Buffer: &Buffer{
					Timekey:       "1m",
					TimekeyWait:   "30s",
					TimekeyUseUtc: true,
				},
			},
		},
	}
}
