package template

import (
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/utils/logging"
)

type Templator struct {
	ApiVersion       string    `yaml:"apiVersion"`
	Kind             string    `yaml:"kind"`
	Metadata         *Metadata `yaml:"metadata"`
	ChartName        string    `yaml:"chartName"`
	ChartVersion     string    `yaml:"chartVersion"`
	ReleaseName      string    `yaml:"releaseName"`
	ReleaseNamespace string    `yaml:"releaseNamespace"`
	ValuesFile       string    `yaml:"valuesFile"`
}

type Metadata struct {
	Name string `yaml:"name"`
}

func NewTemplator(name, chartName, chartVersion, releaseName, releaseNamespace string) *Templator {
	return &Templator{
		ApiVersion: "caos.ch/v1",
		Kind:       "Templator",
		Metadata: &Metadata{
			Name: name,
		},
		ChartName:        chartName,
		ChartVersion:     chartVersion,
		ReleaseName:      releaseName,
		ReleaseNamespace: releaseNamespace,
		ValuesFile:       "values.yaml",
	}
}

func (t *Templator) writeToYaml(templatorFilePath string) error {
	err := helper.StructToYaml(t, templatorFilePath)
	logging.Log("KUSTOMIZE-OpcyPaHsFxThLqH").OnError(err).Debugf("Failed to write templator to file path %s", templatorFilePath)
	return err
}
