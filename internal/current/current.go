package current

import (
	"sort"

	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/orbiter/mntr"
)

type Current struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Resources  []*clientgo.Resource
}

func Get(monitor mntr.Monitor, resourceInfoList []*clientgo.ResourceInfo) *Current {
	globalLabels := labels.GetGlobalLabels()

	resources, err := clientgo.ListResources(monitor, resourceInfoList, globalLabels)
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
