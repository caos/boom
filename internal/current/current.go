package current

import (
	"sort"

	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/orbiter/logging"
)

type Current struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Resources  []*clientgo.Resource
}

func Get(logger logging.Logger) *Current {
	globalLabels := labels.GetGlobalLabels()

	resources, err := clientgo.ListResources(logger, globalLabels)
	if err != nil {
		return nil
	}

	sort.Sort(clientgo.ResourceSorter(resources))

	return &Current{
		APIVersion: "boom.caos.ch/v1beta1",
		Kind:       "currentstate",
		Resources:  resources,
	}
}
