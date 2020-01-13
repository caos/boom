package helper

import (
	"github.com/caos/boom/internal/app/name"
	"github.com/caos/boom/internal/app/templator"
	"github.com/caos/boom/internal/app/templator/helm"
	"github.com/caos/orbiter/logging"
)

func NewTemplator(logger logging.Logger, overlay string, baseDirectoryPath string, templatorName name.Templator) templator.Templator {
	switch templatorName {
	case helm.GetName():
		return helm.New(logger, overlay, baseDirectoryPath)
	}

	return nil
}
