package template

import "github.com/caos/toolsop/internal/helper"

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
		ChartName:        name,
		ChartVersion:     chartVersion,
		ReleaseName:      releaseName,
		ReleaseNamespace: releaseNamespace,
		ValuesFile:       "values.yaml",
	}
}

func (t *Templator) writeToYaml(templatorFilePath string) error {
	return helper.StructToYaml(t, templatorFilePath)
}
