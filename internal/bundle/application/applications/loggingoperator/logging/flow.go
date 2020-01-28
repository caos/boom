package logging

type FlowConfig struct {
	Name         string
	Namespace    string
	SelectLabels map[string]string
	Outputs      []string
	ParserType   string
}

type Parse struct {
	Type string `yaml:"type"`
}
type Parser struct {
	RemoveKeyNameField bool   `yaml:"remove_key_name_field"`
	ReserveData        bool   `yaml:"reserve_data"`
	Parse              *Parse `yaml:"parse"`
}

type Filter struct {
	Parser        *Parser           `yaml:"parser,omitempty"`
	TagNormaliser map[string]string `yaml:"tag_normaliser,omitempty"`
}

type FlowSpec struct {
	Filters    []*Filter         `yaml:"filters,omitempty"`
	Selectors  map[string]string `yaml:"selectors,omitempty"`
	OutputRefs []string          `yaml:"outputRefs"`
}

type Flow struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   *Metadata `yaml:"metadata"`
	Spec       *FlowSpec `yaml:"spec"`
}

func NewFlow(conf *FlowConfig) *Flow {
	return &Flow{
		APIVersion: "logging.banzaicloud.io/v1beta1",
		Kind:       "Flow",
		Metadata: &Metadata{
			Name:      conf.Name,
			Namespace: conf.Namespace,
		},
		Spec: &FlowSpec{
			Filters: []*Filter{
				&Filter{
					Parser: &Parser{
						RemoveKeyNameField: true,
						ReserveData:        true,
						Parse: &Parse{
							Type: conf.ParserType,
						},
					},
				},
				&Filter{
					TagNormaliser: map[string]string{
						"metadata":      "${namespace}.${container}.${pod}",
						"metadata_name": "${namespace_name}.${container_name}.${pod_name}",
					},
				},
			},
			Selectors:  conf.SelectLabels,
			OutputRefs: conf.Outputs,
		},
	}
}
