package defaults

import (
	"os"
	"path/filepath"
	"strings"
)

var Defaults = initDefaults()
var defaultValue = "default"

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
	def, ok := Defaults[application]
	if ok && def.defaultNamespace != "" {
		return strings.Join([]string{overlay, Defaults[application].defaultNamespace}, "-")
	}

	//default namespace in combination with crd name
	if Defaults[defaultValue].defaultNamespace != "" {
		return strings.Join([]string{overlay, Defaults[defaultValue].defaultNamespace}, "-")
	}

	//crd name
	return overlay
}

func GetPrefix(overlay, application, prefix string) string {

	if prefix != "" {
		return prefix
	}

	def, ok := Defaults[application]
	if ok && def.defaultPrefix != "" {
		return strings.Join([]string{overlay, Defaults[application].defaultPrefix}, "-")
	}

	if Defaults[defaultValue].defaultPrefix != "" {
		return strings.Join([]string{overlay, Defaults[defaultValue].defaultPrefix}, "-")
	}

	return overlay
}

func GetResultFileDirectory(overlay, path, application string) string {
	def, ok := Defaults[application]
	if ok {
		filepath.Join(path, def.resultsDirectoryName, overlay)
	}

	return filepath.Join(path, Defaults[defaultValue].resultsDirectoryName, overlay)
}

func GetResultFilePath(overlay, path, application string) string {
	def, ok := Defaults[application]
	if ok {
		filepath.Join(GetResultFileDirectory(overlay, path, application), def.resultsFileName)
	}

	return filepath.Join(GetResultFileDirectory(overlay, path, application), Defaults[defaultValue].resultsFileName)
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

	defaults[defaultValue] = NewValue("results", "results.yaml", "system", "")
	// defaults["ambassador"] = NewValue("results", "results.yaml", "system", "")
	// defaults["cert-manager"] = NewValue("results", "results.yaml", "system", "")
	// defaults["grafana"] = NewValue("results", "results.yaml", "system", "")
	// defaults["logging-operator"] = NewValue("results", "results.yaml", "system", "")
	// defaults["prometheus"] = NewValue("results", "results.yaml", "system", "")
	// defaults["prometheus-node-exporter"] = NewValue("results", "results.yaml", "system", "")
	// defaults["prometheus-operator"] = NewValue("results", "results.yaml", "system", "")

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
