package template

import (
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/utils/logging"
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
	err := helper.StructToYaml(n, filePath)
	logging.Log("KUSTOMIZE-M0FA6gOrtD32Pgf").OnError(err).Debug("Failed to write namespace to file")
	return err
}
