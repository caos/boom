package templator

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/name"
	"github.com/caos/boom/internal/app/templator/helm"
	"github.com/caos/orbiter/logging"
)

type Templator interface {
	AddApplication(name.Application, interface{})
	PrepareTemplate(name.Application, *v1beta1.ToolsetSpec)
	Template(name.Application)
	GetResultsFilePath(appName name.Application) string
	CleanUp() error
}

func New(logger logging.Logger, overlay string, baseDirectoryPath string, templator name.Templator) Templator {
	switch templator {
	case helm.GetName():
		return helm.New(logger, overlay, baseDirectoryPath)
	}

	return nil
}
