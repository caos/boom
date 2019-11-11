package template

import "github.com/caos/toolsop/internal/helper"

type Namespace struct {
	ApiVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   *Metadata `yaml:"metadata"`
}

func NewNamespace(name string) *Namespace {
	return &Namespace{
		ApiVersion: "v1",
		Kind:       "Namespace",
		Metadata: &Metadata{
			Name: name,
		},
	}
}

func (n *Namespace) writeToYaml(filePath string) error {
	return helper.StructToYaml(n, filePath)
}
