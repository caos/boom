package defaults

import (
	"os"
	"path/filepath"
	"strings"
)

var Defaults = initDefaults()

type value struct {
	resultsDirectoryName string
	resultsFileName      string
	defaultNamespace     string
	defaultPrefix        string
}

func GetNamespace(overlay, application, specNamespace, namespace string) string {
	//applicationspecific namespace
	if namespace != "" {
		return namespace
	}

	//crdspec defined namespace
	if specNamespace != "" {
		return specNamespace
	}

	//default applicationspecific namespace in combination with crd name
	if Defaults[application].defaultNamespace != "" {
		return strings.Join([]string{overlay, Defaults[application].defaultNamespace}, "-")
	}

	//crd name
	return overlay
}

func GetPrefix(overlay, application, prefix string) string {

	if prefix != "" {
		return prefix
	}

	if Defaults[application].defaultPrefix != "" {
		return strings.Join([]string{overlay, Defaults[application].defaultPrefix}, "-")
	}

	return overlay
}

func GetResultFileDirectory(overlay, path, application string) string {
	return filepath.Join(path, Defaults[application].resultsDirectoryName, overlay)
}

func GetResultFilePath(overlay, path, application string) string {
	return filepath.Join(GetResultFileDirectory(overlay, path, application), Defaults[application].resultsFileName)
}

func PrepareForResultOutput(resultsFileDirectory string) error {
	if err := os.RemoveAll(resultsFileDirectory); err != nil {
		return err
	}
	if err := os.MkdirAll(resultsFileDirectory, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func initDefaults() map[string]*value {
	defaults := make(map[string]*value, 0)

	defaults["ambassador"] = NewValue("results", "results.yaml", "system", "")
	defaults["cert-manager"] = NewValue("results", "results.yaml", "system", "")
	defaults["grafana"] = NewValue("results", "results.yaml", "system", "")
	defaults["logging-operator"] = NewValue("results", "results.yaml", "system", "")
	defaults["prometheus"] = NewValue("results", "results.yaml", "system", "")
	defaults["prometheus-node-exporter"] = NewValue("results", "results.yaml", "system", "")
	defaults["prometheus-operator"] = NewValue("results", "results.yaml", "system", "")

	return defaults
}

func NewValue(dir, file, namespace, prefix string) *value {
	return &value{
		resultsDirectoryName: dir,
		resultsFileName:      file,
		defaultNamespace:     namespace,
		defaultPrefix:        prefix,
	}
}
