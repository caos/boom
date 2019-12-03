package template

import (
	"github.com/caos/orbiter/logging"
	"github.com/caos/toolsop/internal/helper"
	"github.com/pkg/errors"
)

type Namespace struct {
	ApiVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   *Metadata `yaml:"metadata"`
	logger     logging.Logger
}

func NewNamespace(logger logging.Logger, name string) *Namespace {
	return &Namespace{
		ApiVersion: "v1",
		Kind:       "Namespace",
		Metadata: &Metadata{
			Name: name,
		},
		logger: logger,
	}
}

func (n *Namespace) writeToYaml(filePath string) error {
	err := helper.StructToYaml(n, filePath)
	if err != nil {
		n.logger.WithFields(map[string]interface{}{"logID": "KUSTOMIZE-M0FA6gOrtD32Pgf"}).Error(errors.Wrap(err, "Failed to write namespace to file"))
	}
	return err
}
