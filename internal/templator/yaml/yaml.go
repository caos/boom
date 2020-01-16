package yaml

import (
	"path/filepath"

	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

const (
	templatorName name.Templator = "yaml"
)

func GetName() name.Templator {
	return templatorName
}

type YAML struct {
	logger                 logging.Logger
	status                 error
	overlay                string
	templatorDirectoryPath string
}

func New(logger logging.Logger, overlay, templatorDirectoryPath string) *YAML {
	return &YAML{
		logger:                 logger,
		overlay:                overlay,
		templatorDirectoryPath: templatorDirectoryPath,
	}
}

func (y *YAML) GetStatus() error {
	return y.status
}

func (y *YAML) getResultsFileDirectory(appName name.Application, overlay, basePath string) string {
	return filepath.Join(basePath, appName.String(), overlay, "results")
}

func (y *YAML) GetResultsFilePath(appName name.Application, overlay, basePath string) string {
	return filepath.Join(y.getResultsFileDirectory(appName, overlay, basePath), "results.yaml")
}

func (y *YAML) CleanUp() templator.Templator {
	return y
}

func checkTemplatorInterface(templatorInterface interface{}) (templator.YamlApplication, error) {
	app, isTemplator := templatorInterface.(templator.YamlApplication)
	if !isTemplator {
		err := errors.Errorf("YAML templating interface not implemented")
		return nil, err
	}

	return app, nil
}
