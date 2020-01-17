package logging

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

func NewOutput(name string, namespace string, url string) *Output {
	return &Output{
		APIVersion: "logging.banzaicloud.io/v1beta1",
		Kind:       "Output",
		Metadata: &Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Spec: &OutputSpec{
			Loki: &Loki{
				URL:                       url,
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
