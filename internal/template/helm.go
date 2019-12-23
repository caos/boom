package template

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/toolset"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

var (
	kustomizationFileName = "kustomization.yaml"
	templatorFileName     = "templator.yaml"
	valuesFileName        = "values.yaml"
	namespaceFileName     = "namespace.yaml"
)

type Helm struct {
	ToolsDirectoryPath string
	Toolsets           *toolset.Toolsets
	Applications       map[string]*Application
	Overlay            string
	logger             logging.Logger
}

type Application struct {
	ChartName    string
	ChartVersion string
	IndexName    string
	IndexUrl     string
	ImageTags    map[string]string
}

func NewHelm(logger logging.Logger, toolsDirectoryPath string, toolsets *toolset.Toolsets, crdName, overlay string) (*Helm, error) {
	applications := make(map[string]*Application, 0)
	helm := &Helm{
		ToolsDirectoryPath: toolsDirectoryPath,
		Toolsets:           toolsets,
		Overlay:            overlay,
		Applications:       applications,
		logger:             logger,
	}

	logger.WithFields(map[string]interface{}{
		"logID": "HELM-FuyctETdHmsd7xH",
	}).Info("Collecting list of applications from provided toolsets")
	helm.collectApplications(crdName)

	return helm, nil
}

func (h *Helm) CleanUp() error {
	logFields := map[string]interface{}{
		"overlay": h.Overlay,
	}

	for name := range h.Applications {
		logFields["application"] = name
		logFields["logID"] = "HELM-YCXTWaGws6NsMNJ"
		h.logger.WithFields(logFields).Info("Cleaning up templators")
		templatorDirectoryPath := filepath.Join(h.ToolsDirectoryPath, name, h.Overlay)
		if err := errors.Wrapf(os.RemoveAll(templatorDirectoryPath), "Failed cleanup templators application %s overlay %s", name, h.Overlay); err != nil {
			return err
		}
	}
	return nil
}

func (h *Helm) collectApplications(crdName string) {
	for _, toolset := range h.Toolsets.Toolsets {
		if toolset.Name == crdName {
			for _, application := range toolset.Applications {
				app := &Application{
					ChartName:    application.File.Chart.Name,
					ChartVersion: application.File.Chart.Version,
					IndexName:    application.File.Chart.Index.Name,
					IndexUrl:     application.File.Chart.Index.URL,
					ImageTags:    application.File.ImageTags,
				}
				h.Applications[application.Name] = app
			}
		}
	}
}

func (h *Helm) GetDefaultValuesPath(appName string) string {
	application := h.Applications[appName]
	dir := strings.Join([]string{application.ChartName, application.ChartVersion}, "-")
	defaultValuesFilePath := filepath.Join(h.ToolsDirectoryPath, "charts", dir, application.ChartName, "values.yaml")
	return defaultValuesFilePath
}

func (h *Helm) GetImageTags(appName string) map[string]string {
	application := h.Applications[appName]
	return application.ImageTags
}

func (h *Helm) PrepareTemplate(appName, releaseName, releaseNamespace string, writeValues func(path string) error) error {
	logFields := map[string]interface{}{
		"application": appName,
		"overlay":     h.Overlay,
	}

	logFields["logID"] = "HELM-YF1XEbmiazmdCmN"
	h.logger.WithFields(logFields).Info("Generating templator")
	if err := h.generateTemplator(appName, releaseName, releaseNamespace, writeValues); err != nil {
		return err
	}
	return nil
}

func (h *Helm) Template(appName, resultfilepath string) error {

	base, err := filepath.Abs(h.ToolsDirectoryPath)
	if err != nil {
		return err
	}
	result, err := filepath.Abs(resultfilepath)
	if err != nil {
		return err
	}

	cdCommand := strings.Join([]string{"cd", base}, " ")
	startCommand := strings.Join([]string{"./start.sh", appName, h.Overlay}, " ")
	startCommand = strings.Join([]string{startCommand, ">", result}, " ")
	command := strings.Join([]string{cdCommand, startCommand}, " && ")

	cmd := exec.Command("/bin/sh", "-c", command)
	return errors.Wrapf(helper.Run(h.logger, *cmd), "Failed on templating overlay %s application %s", h.Overlay, appName)
}

func (h *Helm) generateTemplator(appName, releaseName, releaseNamespace string, writeValues func(path string) error) error {

	logFields := map[string]interface{}{
		"application": appName,
		"overlay":     h.Overlay,
	}

	logFields["logID"] = "HELM-4cR49WgtuN1757N"
	h.logger.WithFields(logFields).Info("Deleting old files for templator")
	templatorDirectoryPath := filepath.Join(h.ToolsDirectoryPath, appName, h.Overlay)
	_ = os.RemoveAll(templatorDirectoryPath)
	logFields["logID"] = "HELM-aoqK04FN3QxbOQx"
	h.logger.WithFields(logFields).Info("Creating folder for templator")
	_ = os.MkdirAll(templatorDirectoryPath, os.ModePerm)

	// values file
	valuesFilePath := filepath.Join(templatorDirectoryPath, valuesFileName)
	if err := writeValues(valuesFilePath); err != nil {
		return err
	}

	// templator with valuesfilename
	templatorFilePath := filepath.Join(templatorDirectoryPath, templatorFileName)
	app := h.Applications[appName]
	logFields["logID"] = "HELM-bRlJLZvwmDxrNIN"
	h.logger.WithFields(logFields).Info("Generating templator")
	templator := NewTemplator(appName, app.ChartName, releaseName, releaseNamespace)
	err := templator.writeToYaml(templatorFilePath)
	if err != nil {
		return err
	}

	namespaceFilePath := filepath.Join(templatorDirectoryPath, namespaceFileName)
	logFields["logID"] = "HELM-jG5n3lf5TJQGdLc"
	h.logger.WithFields(logFields).Info("Generating namespace")
	namespace := NewNamespace(releaseNamespace)
	err = namespace.writeToYaml(namespaceFilePath)
	if err != nil {
		return err
	}
	//kustomization with templatorfilename
	kustomizationFilePath := filepath.Join(templatorDirectoryPath, kustomizationFileName)
	recources := []string{namespaceFileName}
	generators := []string{templatorFileName}
	logFields["logID"] = "HELM-3o5gZY6roaWisQ7"
	h.logger.WithFields(logFields).Info("Generating templator kustomize")
	err = generateKustomization(kustomizationFilePath, recources, generators)
	if err != nil {
		return err
	}

	return nil
}
