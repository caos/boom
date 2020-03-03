package helper

import (
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/boom/internal/templator/helm"
	"github.com/caos/boom/internal/templator/yaml"
	"github.com/caos/orbiter/mntr"
)

func NewTemplator(monitor mntr.Monitor, overlay string, baseDirectoryPath string, templatorName name.Templator) templator.Templator {
	switch templatorName {
	case helm.GetName():
		return helm.New(monitor, overlay, baseDirectoryPath)
	case yaml.GetName():
		return yaml.New(monitor, overlay, baseDirectoryPath)
	}

	return nil
}
