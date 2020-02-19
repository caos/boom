package current

import (
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
	currentLogger := logger.WithFields(map[string]interface{}{
		"logID": "CURRENT-fEIFck9R1KRY5e8",
	})

	globalLabels := labels.GetGlobalLabels()

	resources, err := clientgo.ListResources(currentLogger, globalLabels)
	if err != nil {
		return nil
	}

	return &Current{
		APIVersion: "boom.caos.ch/v1beta1",
		Kind:       "currentstate",
		Resources:  resources,
	}
}
