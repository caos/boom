package helm

import (
	"github.com/caos/boom/internal/helper"
	"github.com/pkg/errors"
)

type TemplatorFile struct {
	APIVersion       string    `yaml:"apiVersion"`
	Kind             string    `yaml:"kind"`
	Metadata         *Metadata `yaml:"metadata"`
	ChartName        string    `yaml:"chartName,omitempty"`
	ReleaseName      string    `yaml:"releaseName,omitempty"`
	ReleaseNamespace string    `yaml:"releaseNamespace,omitempty"`
	ValuesFile       string    `yaml:"valuesFile,omitempty"`
}

type Metadata struct {
	Name string `yaml:"name,omitempty"`
}

func NewTemplatorFile(name, chartName, releaseName, releaseNamespace string) *TemplatorFile {
	return &TemplatorFile{
		APIVersion: "caos.ch/v1",
		Kind:       "Templator",
		Metadata: &Metadata{
			Name: name,
		},
		ChartName:        chartName,
		ReleaseName:      releaseName,
		ReleaseNamespace: releaseNamespace,
		ValuesFile:       "values.yaml",
	}
}

func (t *TemplatorFile) writeToYaml(templatorFilePath string) error {
	return errors.Wrapf(helper.StructToYaml(t, templatorFilePath), "Failed to write templator to file path %s", templatorFilePath)
}
