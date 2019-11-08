package template

import (
	"os"
	"os/exec"
	"strings"

	"github.com/caos/toolsop/internal/toolset"
)

var (
	fetcherDirectoryName   = "fetchers"
	fetcherFileName        = "fetcher.yaml"
	kustomizationFileName  = "kustomization.yaml"
	templatorDirectoryName = "templators"
	templatorFileName      = "templator.yaml"
	valuesFileName         = "values.yaml"
)

type Helm struct {
	ToolsDirectoryPath string
	Toolsets           *toolset.Toolsets
	Applications       map[string]*Application
	Overlay            string
}

type Application struct {
	ChartName    string
	ChartVersion string
	IndexName    string
	IndexUrl     string
	ImageTags    map[string]string
}

func NewHelm(toolsDirectoryPath string, toolsets *toolset.Toolsets, crdName, crdVersion, overlay string) (*Helm, error) {
	applications := make(map[string]*Application, 0)
	helm := &Helm{
		ToolsDirectoryPath: toolsDirectoryPath,
		Toolsets:           toolsets,
		Overlay:            overlay,
		Applications:       applications,
	}

	helm.collectApplications(crdName, crdVersion)

	if err := helm.generateFetchers(); err != nil {
		return nil, err
	}

	if err := helm.fetchAllCharts(); err != nil {
		return nil, err
	}

	return helm, nil
}

func (h *Helm) collectApplications(crdName, crdVersion string) {
	for _, toolset := range h.Toolsets.Toolsets {
		if toolset.Name == crdName {
			for _, version := range toolset.Versions {
				if version.Version == crdVersion {
					for _, application := range version.Applications {
						app := &Application{
							ChartName:    application.Chart.Name,
							ChartVersion: application.Chart.Version,
							IndexName:    application.Chart.Index.Name,
							IndexUrl:     application.Chart.Index.URL,
							ImageTags:    application.ImageTags,
						}
						h.Applications[application.Name] = app
					}
				}
			}
		}
	}
}

func (h *Helm) GetDefaultValuesPath(appName string) string {
	application := h.Applications[appName]
	dir := strings.Join([]string{application.ChartName, application.ChartVersion}, "-")
	defaultValuesFilePath := strings.Join([]string{h.ToolsDirectoryPath, "charts", dir, application.ChartName, "values.yaml"}, "/")
	return defaultValuesFilePath
}

func (h *Helm) GetImageTags(appName string) map[string]string {
	application := h.Applications[appName]
	return application.ImageTags
}

func (h *Helm) generateFetchers() error {
	for name, application := range h.Applications {
		fetcherDirectoryPath := strings.Join([]string{h.ToolsDirectoryPath, name, fetcherDirectoryName, h.Overlay}, "/")
		_ = os.MkdirAll(fetcherDirectoryPath, os.ModePerm)

		fetcherFilePath := strings.Join([]string{fetcherDirectoryPath, fetcherFileName}, "/")
		fetcher := NewFetcher(name, application.ChartName, application.ChartVersion, application.IndexName, application.IndexUrl)
		if err := fetcher.writeToYaml(fetcherFilePath); err != nil {
			return nil
		}

		kustomizationFilePath := strings.Join([]string{fetcherDirectoryPath, kustomizationFileName}, "/")
		filePaths := []string{fetcherFileName}
		if err := generateKustomization(kustomizationFilePath, filePaths); err != nil {
			return err
		}
	}
	return nil
}

func (h *Helm) fetchAllCharts() error {
	for name, _ := range h.Applications {
		cdCommand := strings.Join([]string{"cd", h.ToolsDirectoryPath}, " ")
		overlay := strings.Join([]string{fetcherDirectoryName, h.Overlay}, "/")
		startCommand := strings.Join([]string{"./start.sh", name, overlay}, " ")
		command := strings.Join([]string{cdCommand, startCommand}, " && ")

		cmd := exec.Command("/bin/sh", "-c", command)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (h *Helm) Template(appName, releaseName, releaseNamespace, resultfilepath string, writeValues func(path string) error) error {

	if err := h.generateTemplators(appName, releaseName, releaseNamespace, writeValues); err != nil {
		return nil
	}

	cdCommand := strings.Join([]string{"cd", h.ToolsDirectoryPath}, " ")
	overlay := strings.Join([]string{templatorDirectoryName, h.Overlay}, "/")
	startCommand := strings.Join([]string{"./start.sh", appName, overlay}, " ")
	startCommand = strings.Join([]string{startCommand, ">>", resultfilepath}, " ")
	command := strings.Join([]string{cdCommand, startCommand}, " && ")

	cmd := exec.Command("/bin/sh", "-c", command)
	return cmd.Run()
}

func (h *Helm) generateTemplators(appName, releaseName, releaseNamespace string, writeValues func(path string) error) error {
	templatorDirectoryPath := strings.Join([]string{h.ToolsDirectoryPath, appName, templatorDirectoryName, h.Overlay}, "/")
	_ = os.MkdirAll(templatorDirectoryPath, os.ModePerm)

	// values file
	valuesFilePath := strings.Join([]string{templatorDirectoryPath, valuesFileName}, "/")
	if err := writeValues(valuesFilePath); err != nil {
		return err
	}

	// templator with valuesfilename
	templatorFilePath := strings.Join([]string{templatorDirectoryPath, templatorFileName}, "/")
	app := h.Applications[appName]
	templator := NewTemplator(appName, app.ChartName, app.ChartVersion, releaseName, releaseNamespace)
	err := templator.writeToYaml(templatorFilePath)
	if err != nil {
		return err
	}

	//kustomization with templatorfilename
	kustomizationFilePath := strings.Join([]string{templatorDirectoryPath, kustomizationFileName}, "/")
	filePaths := []string{templatorFileName}
	err = generateKustomization(kustomizationFilePath, filePaths)
	if err != nil {
		return err
	}

	return nil
}
