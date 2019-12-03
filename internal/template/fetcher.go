package template

import (
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	"github.com/caos/toolsop/internal/helper"
)

type Fetcher struct {
	ApiVersion   string    `yaml:"apiVersion"`
	Kind         string    `yaml:"kind"`
	Metadata     *Metadata `yaml:"metadata"`
	ChartName    string    `yaml:"chartName"`
	ChartVersion string    `yaml:"chartVersion"`
	IndexName    string    `yaml:"indexName,omitempty"`
	IndexUrl     string    `yaml:"indexUrl,omitempty"`
	logger       logging.Logger
}

func NewFetcher(logger logging.Logger, name, chartName, chartVersion, indexName, indexUrl string) *Fetcher {
	return &Fetcher{
		ApiVersion: "caos.ch/v1",
		Kind:       "Fetcher",
		Metadata: &Metadata{
			Name: name,
		},
		ChartName:    chartName,
		ChartVersion: chartVersion,
		IndexName:    indexName,
		IndexUrl:     indexUrl,
		logger:       logger,
	}
}

func (f *Fetcher) writeToYaml(fetcherFilePath string) error {
	err := errors.Wrapf(helper.StructToYaml(f, fetcherFilePath), "Failed to write fetcher to yaml, path: %s", fetcherFilePath)
	if err != nil {
		f.logger.WithFields(map[string]interface{}{
			"logID": "FETCHER-qIBSyOB0CD37u0P",
		}).Error(err)
	}
	return err
}
