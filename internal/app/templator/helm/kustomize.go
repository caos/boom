package helm

import (
	"github.com/pkg/errors"

	"github.com/caos/boom/internal/helper"
)

type Kustomization struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Resources  []string `yaml:"resources,omitempty"`
	Generators []string `yaml:"generators,omitempty"`
}

func generateKustomization(kustomizationFilePath string, resources []string, generators []string) error {
	kustomization := &Kustomization{
		APIVersion: "kustomize.config.k8s.io/v1beta1",
		Kind:       "Kustomization",
		Resources:  resources,
		Generators: generators,
	}

	return errors.Wrap(helper.StructToYaml(kustomization, kustomizationFilePath), "Failed to write kustomize to file")
}

