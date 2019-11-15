package template

import (
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/utils/logging"
)

type Kustomization struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Resources  []string `yaml:"resources,omitempty"`
	Generators []string `yaml:"generators,omitempty"`
}

func generateKustomization(kustomizationFilePath string, resources []string, generators []string) error {
	kustomization := &Kustomization{
		ApiVersion: "kustomize.config.k8s.io/v1beta1",
		Kind:       "Kustomization",
		Resources:  resources,
		Generators: generators,
	}

	err := helper.StructToYaml(kustomization, kustomizationFilePath)
	logging.Log("KUSTOMIZE-QCiaUk3u7mwOhLe").OnError(err).Debug("Failed to write kustomize to file")
	return err
}
