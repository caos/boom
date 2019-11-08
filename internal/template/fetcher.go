package template

import (
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
}

func NewFetcher(name, chartName, chartVersion, indexName, indexUrl string) *Fetcher {
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
	}
}

func (f *Fetcher) writeToYaml(fetcherFilePath string) error {
	return helper.StructToYaml(f, fetcherFilePath)
}
