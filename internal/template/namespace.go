package template

import (
	"github.com/caos/orbiter/logging"
	"github.com/caos/toolsop/internal/helper"
	"github.com/pkg/errors"
)

type Namespace struct {
	ApiVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   *Metadata      `yaml:"metadata"`
	logger     logging.Logger `yaml:"-"`
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
	return errors.Wrap(helper.StructToYaml(n, filePath), "Failed to write namespace to file")
}
