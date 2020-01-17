package logging

type Parse struct {
	Type string `yaml:"type"`
}
type Parser struct {
	RemoveKeyNameField bool   `yaml:"remove_key_name_field"`
	ReserveData        bool   `yaml:"reserve_data"`
	Parse              *Parse `yaml:"parse"`
}
type TagNormaliser struct {
	Format string
}

type Filter struct {
	Parser        *Parser        `yaml:"parser,omitempty"`
	TagNormaliser *TagNormaliser `yaml:"tag_normaliser,omitempty"`
}

type FlowSpec struct {
	Filters    []*Filter         `yaml:"filters"`
	Selectors  map[string]string `yaml:"selectors,omitempty"`
	OutputRefs []string          `yaml:"outputRefs"`
}

type Flow struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   *Metadata `yaml:"metadata"`
	Spec       *FlowSpec `yaml:"spec"`
}

func NewFlow(name, namespace string, selectLabels map[string]string, outputs []string) *Flow {
	return &Flow{
		APIVersion: "logging.banzaicloud.io/v1beta1",
		Kind:       "Flow",
		Metadata: &Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Spec: &FlowSpec{
			Filters: []*Filter{
				&Filter{
					Parser: &Parser{
						RemoveKeyNameField: true,
						Parse: &Parse{
							Type: "nginx",
						},
					},
				},
				&Filter{
					TagNormaliser: &TagNormaliser{
						Format: "${namespace_name}.${pod_name}.${container_name}",
					},
				},
			},
			Selectors:  selectLabels,
			OutputRefs: outputs,
		},
	}
}
