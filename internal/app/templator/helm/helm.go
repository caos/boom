package helm

import (
	"os"
	"path/filepath"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/bundle/application/chart"
	"github.com/caos/boom/internal/app/name"
	"github.com/caos/boom/internal/app/templator"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

const (
	templatorName name.Templator = "helm"
)

func GetName() name.Templator {
	return templatorName
}

type Templator interface {
	GetName() name.Application
	SpecToHelmValues(spec *v1beta1.ToolsetSpec) interface{}
	GetChartInfo() *chart.Chart
	GetNamespace() string
}

type Helm struct {
	overlay                string
	status                 error
	logger                 logging.Logger
	templatorDirectoryPath string
}

func New(logger logging.Logger, overlay, templatorDirectoryPath string) templator.Templator {
	return &Helm{
		logger:                 logger,
		templatorDirectoryPath: templatorDirectoryPath,
		overlay:                overlay,
	}
}

func (h *Helm) CleanUp() templator.Templator {
	if h.status != nil {
		return h
	}

	h.status = os.RemoveAll(h.templatorDirectoryPath)
	return h
}

func (h *Helm) getResultsFileDirectory(appName name.Application, overlay, basePath string) string {
	return filepath.Join(basePath, appName.String(), overlay, "results")
}

func (h *Helm) GetResultsFilePath(appName name.Application, overlay, basePath string) string {
	return filepath.Join(h.getResultsFileDirectory(appName, overlay, basePath), "results.yaml")
}

func (h *Helm) GetStatus() error {
	return h.status
}

func checkTemplatorInterface(templatorInterface interface{}) (Templator, error) {
	templator, isTemplator := templatorInterface.(Templator)
	if !isTemplator {
		logFields := map[string]interface{}{
			"application": templator.GetName().String(),
		}
		logFields["logID"] = "HELM-gHHLU2a49osYzgl"
		err := errors.Errorf("Helm templating interface not implemented")
		return nil, err
	}

	return templator, nil
}
