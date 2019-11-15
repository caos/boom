package template

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/caos/toolsop/internal/toolset"
	"github.com/caos/utils/logging"
)

var (
	fetcherDirectoryName   = "fetchers"
	fetcherFileName        = "fetcher.yaml"
	kustomizationFileName  = "kustomization.yaml"
	templatorDirectoryName = "templators"
	templatorFileName      = "templator.yaml"
	valuesFileName         = "values.yaml"
	namespaceFileName      = "namespace.yaml"
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

	logging.Log("HELM-FuyctETdHmsd7xH").Info("Collecting list of applications from provided toolsets")
	helm.collectApplications(crdName, crdVersion)

	logging.Log("HELM-NVpmSPi56GezX7D").Info("Generating fetchers for necessary helm charts")
	if err := helm.generateFetchers(); err != nil {
		return nil, err
	}

	logging.Log("HELM-AH5kPecXCXAYVGy").Info("Fetching all helm charts to local")
	if err := helm.fetchAllCharts(); err != nil {
		return nil, err
	}

	return helm, nil
}

func (h *Helm) CleanUp() error {
	for name := range h.Applications {
		logging.Log("HELM-cUCPkTW3n2paliN").Infof("Cleanup fetchers application %s overlay %s", name, h.Overlay)
		fetcherDirectoryPath := filepath.Join(h.ToolsDirectoryPath, name, fetcherDirectoryName, h.Overlay)
		if err := os.RemoveAll(fetcherDirectoryPath); err != nil {
			logging.Log("HELM-o9Fj1ljqCjtqKPj").OnError(err).Debugf("Failed cleanup fetchers application %s overlay %s", name, h.Overlay)
			return err
		}

		logging.Log("HELM-YCXTWaGws6NsMNJ").Infof("Cleanup templators application %s overlay %s", name, h.Overlay)
		templatorDirectoryPath := filepath.Join(h.ToolsDirectoryPath, name, templatorDirectoryName, h.Overlay)
		if err := os.RemoveAll(templatorDirectoryPath); err != nil {
			logging.Log("HELM-IOzamCt9i2GFohA").OnError(err).Debugf("Failed cleanup templators application %s overlay %s", name, h.Overlay)
			return err
		}
	}
	return nil
}

func (h *Helm) collectApplications(crdName, crdVersion string) {
	for _, toolset := range h.Toolsets.Toolsets {
		if toolset.Name == crdName {
			for _, version := range toolset.Versions {
				if version.Version == crdVersion {
					for _, application := range version.Applications {
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

func (h *Helm) generateFetchers() error {
	h.CleanUp()

	// as the helm template is per crd instanced, all necessary applications are known
	for name, application := range h.Applications {
		fetcherDirectoryPath := filepath.Join(h.ToolsDirectoryPath, name, fetcherDirectoryName, h.Overlay)
		logging.Log("HELM-Imjm7vrslYIIzIa").Infof("Create fetcher directory for appplication %s and overlay %s", name, h.Overlay)
		_ = os.MkdirAll(fetcherDirectoryPath, os.ModePerm)

		fetcherFilePath := filepath.Join(fetcherDirectoryPath, fetcherFileName)
		logging.Log("HELM-8vzi64f69T7E4y2").Infof("Generate fetcher for appplication %s and overlay %s", name, h.Overlay)
		fetcher := NewFetcher(name, application.ChartName, application.ChartVersion, application.IndexName, application.IndexUrl)
		if err := fetcher.writeToYaml(fetcherFilePath); err != nil {
			return nil
		}

		kustomizationFilePath := filepath.Join(fetcherDirectoryPath, kustomizationFileName)
		filePaths := []string{fetcherFileName}
		logging.Log("HELM-F1SpEvzx9DFTksX").Infof("Generate fetcher-kustomize for appplication %s and overlay %s", name, h.Overlay)
		if err := generateKustomization(kustomizationFilePath, []string{}, filePaths); err != nil {
			return err
		}
	}
	return nil
}

func (h *Helm) fetchAllCharts() error {
	for name, _ := range h.Applications {
		logging.Log("HELM-QyuO17EOfqoEDP8").Infof("Fetch chart for application %s", name)
		cdCommand := strings.Join([]string{"cd", h.ToolsDirectoryPath}, " ")
		overlay := strings.Join([]string{fetcherDirectoryName, h.Overlay}, "/")
		startCommand := strings.Join([]string{"./start.sh", name, overlay}, " ")
		command := strings.Join([]string{cdCommand, startCommand}, " && ")

		cmd := exec.Command("/bin/sh", "-c", command)
		if err := cmd.Run(); err != nil {
			logging.Log("HELM-QyuO17EOfqoEDP8").OnError(err).Debugf("Failed to fetch chart for application %s", name)
			return err
		}
	}
	return nil
}

func (h *Helm) Template(appName, releaseName, releaseNamespace, resultfilepath string, writeValues func(path string) error) error {

	logging.Log("HELM-YF1XEbmiazmdCmN").Infof("Generating templator for overlay %s application %s", h.Overlay, appName)
	if err := h.generateTemplator(appName, releaseName, releaseNamespace, writeValues); err != nil {
		return nil
	}

	base, err := filepath.Abs(h.ToolsDirectoryPath)
	if err != nil {
		return err
	}
	result, err := filepath.Abs(resultfilepath)
	if err != nil {
		return err
	}

	cdCommand := strings.Join([]string{"cd", base}, " ")
	overlay := strings.Join([]string{templatorDirectoryName, h.Overlay}, "/")
	startCommand := strings.Join([]string{"./start.sh", appName, overlay}, " ")
	startCommand = strings.Join([]string{startCommand, ">", result}, " ")
	command := strings.Join([]string{cdCommand, startCommand}, " && ")

	cmd := exec.Command("/bin/sh", "-c", command)
	err = cmd.Run()
	logging.Log("HELM-mzF3DUV1zAi4vom").OnError(err).Debugf("Failed on templating overlay %s application %s", h.Overlay, appName)
	return err
}

func (h *Helm) generateTemplator(appName, releaseName, releaseNamespace string, writeValues func(path string) error) error {
	templatorDirectoryPath := filepath.Join(h.ToolsDirectoryPath, appName, templatorDirectoryName, h.Overlay)
	logging.Log("HELM-4cR49WgtuN1757N").Infof("Delete old files for templator overlay %s application %s", h.Overlay, appName)
	_ = os.RemoveAll(templatorDirectoryPath)
	logging.Log("HELM-aoqK04FN3QxbOQx").Infof("Create folder for templator overlay %s application %s", h.Overlay, appName)
	_ = os.MkdirAll(templatorDirectoryPath, os.ModePerm)

	// values file
	valuesFilePath := filepath.Join(templatorDirectoryPath, valuesFileName)
	if err := writeValues(valuesFilePath); err != nil {
		return err
	}

	// templator with valuesfilename
	templatorFilePath := filepath.Join(templatorDirectoryPath, templatorFileName)
	app := h.Applications[appName]
	logging.Log("HELM-bRlJLZvwmDxrNIN").Infof("Generate templator overlay %s application %s", h.Overlay, appName)
	templator := NewTemplator(appName, app.ChartName, app.ChartVersion, releaseName, releaseNamespace)
	err := templator.writeToYaml(templatorFilePath)
	if err != nil {
		return err
	}

	namespaceFilePath := filepath.Join(templatorDirectoryPath, namespaceFileName)
	logging.Log("HELM-jG5n3lf5TJQGdLc").Infof("Generate namespace overlay %s application %s", h.Overlay, appName)
	namespace := NewNamespace(releaseNamespace)
	err = namespace.writeToYaml(namespaceFilePath)
	if err != nil {
		return err
	}
	//kustomization with templatorfilename
	kustomizationFilePath := filepath.Join(templatorDirectoryPath, kustomizationFileName)
	recources := []string{namespaceFileName}
	generators := []string{templatorFileName}
	logging.Log("HELM-3o5gZY6roaWisQ7").Infof("Generate templator kustomize overlay %s application %s", h.Overlay, appName)
	err = generateKustomization(kustomizationFilePath, recources, generators)
	if err != nil {
		return err
	}

	return nil
}
