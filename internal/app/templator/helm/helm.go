package helm

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/bundle/application/chart"
	"github.com/caos/boom/internal/app/name"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

const (
	templatorName  name.Templator = "helm"
	helmHomeFolder string         = "helm"
	chartsFolder   string         = "charts"
)

func GetName() name.Templator {
	return templatorName
}

type Templator interface {
	SpecToHelmValues(spec *v1beta1.ToolsetSpec) interface{}
	GetChartInfo() *chart.Chart
}

type TemplatorPreSteps interface {
	HelmPreApplySteps(spec *v1beta1.ToolsetSpec)
}

type Helm struct {
	overlay                string
	releaseName            string
	releaseNamespace       string
	applications           map[name.Application]Templator
	status                 error
	logger                 logging.Logger
	templatorDirectoryPath string
}

var (
	kustomizationFileName = "kustomization.yaml"
	templatorFileName     = "templator.yaml"
	valuesFileName        = "values.yaml"
	namespaceFileName     = "namespace.yaml"
)

func New(logger logging.Logger, overlay, templatorDirectoryPath string) *Helm {
	return &Helm{
		overlay:                overlay,
		logger:                 logger,
		releaseName:            "",            // releaseName,
		releaseNamespace:       "caos-system", //releaseNamespace,
		applications:           make(map[name.Application]Templator, 0),
		templatorDirectoryPath: templatorDirectoryPath,
	}
}

func (h *Helm) CleanUp() error {
	return os.RemoveAll(h.templatorDirectoryPath)
}

func (h *Helm) getResultsFileDirectory(appName name.Application) string {
	return filepath.Join(h.templatorDirectoryPath, appName.String(), h.overlay, "results")
}

func (h *Helm) GetResultsFilePath(appName name.Application) string {
	return filepath.Join(h.getResultsFileDirectory(appName), "results.yaml")
}

func (h *Helm) GetStatus() error {
	return h.status
}

func (h *Helm) AddApplication(appName name.Application, templatorInterface interface{}) {
	templator := templatorInterface.(Templator)

	newApplications := make(map[name.Application]Templator, 0)
	for k, v := range h.applications {
		newApplications[k] = v
	}
	newApplications[appName] = templator
	h.applications = newApplications
}

func (h *Helm) PrepareTemplate(appName name.Application, spec *v1beta1.ToolsetSpec) {

	logFields := map[string]interface{}{
		"application": appName,
		"overlay":     h.overlay,
	}

	logFields["logID"] = "HELM-YF1XEbmiazmdCmN"
	h.logger.WithFields(logFields).Info("Generating templator")
	h.status = h.generateTemplator(appName, spec)

	logFields["logID"] = "HELM-zqYnVAzGXHqBbhu"
	h.logger.WithFields(logFields).Info("Deleting old results")
	h.status = h.deleteResults(appName)
}

func (h *Helm) deleteResults(appName name.Application) error {
	resultsFileDirectory := h.getResultsFileDirectory(appName)
	if err := os.RemoveAll(resultsFileDirectory); err != nil {
		return err
	}
	if err := os.MkdirAll(resultsFileDirectory, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (h *Helm) PreApplySteps(appName name.Application, spec *v1beta1.ToolsetSpec) {
	app := h.applications[appName]

	pre, ok := app.(TemplatorPreSteps)
	if ok {
		pre.HelmPreApplySteps(spec)
	}
}

func (h *Helm) Template(appName name.Application) {

	var base, result string
	resultfilepath := h.GetResultsFilePath(appName)

	base, h.status = filepath.Abs(h.templatorDirectoryPath)
	if h.status != nil {
		return
	}
	result, h.status = filepath.Abs(resultfilepath)
	if h.status != nil {
		return
	}

	cdCommand := strings.Join([]string{"cd", base}, " ")
	startCommand := strings.Join([]string{"./start.sh", appName.String(), h.overlay}, " ")
	startCommand = strings.Join([]string{startCommand, ">", result}, " ")
	command := strings.Join([]string{cdCommand, startCommand}, " && ")

	cmd := exec.Command("/bin/sh", "-c", command)
	h.status = errors.Wrapf(helper.Run(h.logger, *cmd), "Failed on templating overlay %s application %s", h.overlay, appName)
	if h.status != nil {
		return
	}

	helper.DeleteKindFromYaml(result, "Namespace")
}

func (h *Helm) generateTemplator(appName name.Application, spec *v1beta1.ToolsetSpec) error {
	logFields := map[string]interface{}{
		"application": appName,
		"overlay":     h.overlay,
	}

	app := h.applications[appName]
	chartInfo := app.GetChartInfo()
	templatorFileDirectoryPath := filepath.Join(h.templatorDirectoryPath, appName.String(), h.overlay)
	templatorFilePath := filepath.Join(templatorFileDirectoryPath, templatorFileName)
	namespaceFilePath := filepath.Join(templatorFileDirectoryPath, namespaceFileName)
	kustomizationFilePath := filepath.Join(templatorFileDirectoryPath, kustomizationFileName)

	logFields["logID"] = "HELM-4cR49WgtuN1757N"
	h.logger.WithFields(logFields).Info("Deleting old files for templator")
	_ = os.RemoveAll(templatorFileDirectoryPath)
	logFields["logID"] = "HELM-aoqK04FN3QxbOQx"
	h.logger.WithFields(logFields).Info("Creating folder for templator")
	_ = os.MkdirAll(templatorFileDirectoryPath, os.ModePerm)

	// values file
	valuesFilePath := filepath.Join(templatorFileDirectoryPath, valuesFileName)
	values := app.SpecToHelmValues(spec)
	if err := helper.StructToYaml(values, valuesFilePath); err != nil {
		return err
	}

	// templator with valuesfilename
	logFields["logID"] = "HELM-bRlJLZvwmDxrNIN"
	h.logger.WithFields(logFields).Info("Generating templator")
	templator := NewTemplatorFile(appName.String(), chartInfo.Name, h.releaseName, h.releaseNamespace)
	err := templator.writeToYaml(templatorFilePath)
	if err != nil {
		return err
	}
	//namespace file
	logFields["logID"] = "HELM-jG5n3lf5TJQGdLc"
	h.logger.WithFields(logFields).Info("Generating namespace")
	namespace := NewNamespace(h.releaseNamespace)
	err = namespace.writeToYaml(namespaceFilePath)
	if err != nil {
		return err
	}
	//kustomization with templatorfilename
	recources := []string{namespaceFileName}
	generators := []string{templatorFileName}
	logFields["logID"] = "HELM-3o5gZY6roaWisQ7"
	h.logger.WithFields(logFields).Info("Generating templator kustomize")
	return generateKustomization(kustomizationFilePath, recources, generators)
}
