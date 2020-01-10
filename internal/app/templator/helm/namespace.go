package helm

import (
	"github.com/caos/boom/internal/helper"
	"github.com/pkg/errors"
)

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
	return errors.Wrap(helper.StructToYaml(n, filePath), "Failed to write namespace to file")
}

