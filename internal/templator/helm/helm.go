package helm

import (
	"os"
	"path/filepath"

	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/orbiter/mntr"
	"github.com/pkg/errors"
)

const (
	templatorName name.Templator = "helm"
)

func GetName() name.Templator {
	return templatorName
}

func GetPrio() name.Templator {
	return templatorName
}

type Helm struct {
	overlay                string
	monitor                mntr.Monitor
	templatorDirectoryPath string
}

func New(monitor mntr.Monitor, overlay, templatorDirectoryPath string) templator.Templator {
	return &Helm{
		monitor:                monitor,
		templatorDirectoryPath: templatorDirectoryPath,
		overlay:                overlay,
	}
}

func (h *Helm) CleanUp() error {
	return os.RemoveAll(h.templatorDirectoryPath)
}

func (h *Helm) getResultsFileDirectory(appName name.Application, overlay, basePath string) string {
	return filepath.Join(basePath, appName.String(), overlay, "results")
}

func (h *Helm) GetResultsFilePath(appName name.Application, overlay, basePath string) string {
	return filepath.Join(h.getResultsFileDirectory(appName, overlay, basePath), "results.yaml")
}

func checkTemplatorInterface(templatorInterface interface{}) (templator.HelmApplication, error) {
	app, isTemplator := templatorInterface.(templator.HelmApplication)
	if !isTemplator {
		err := errors.Errorf("Helm templating interface not implemented")
		return nil, err
	}

	return app, nil
}
