package template

import (
	"github.com/caos/orbiter/logging"
	"github.com/caos/toolsop/internal/helper"
	"github.com/pkg/errors"
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
	logger           logging.Logger
}

type Metadata struct {
	Name string `yaml:"name"`
}

func NewTemplator(logger logging.Logger, name, chartName, chartVersion, releaseName, releaseNamespace string) *Templator {
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
		logger:           logger,
	}
}

func (t *Templator) writeToYaml(templatorFilePath string) error {
	err := helper.StructToYaml(t, templatorFilePath)
	if err != nil {
		t.logger.WithFields(map[string]interface{}{"logID": "KUSTOMIZE-OpcyPaHsFxThLqH"}).Error(errors.Wrapf(err, "Failed to write templator to file path %s", templatorFilePath))
	}
	return err
}
