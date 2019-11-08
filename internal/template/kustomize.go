package template

import "github.com/caos/toolsop/internal/helper"

type Kustomization struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Generators []string `yaml:"generators"`
}

func generateKustomization(kustomizationFilePath string, filePaths []string) error {
	kustomization := &Kustomization{
		ApiVersion: "kustomize.config.k8s.io/v1beta1",
		Kind:       "Kustomization",
		Generators: filePaths,
	}

	return helper.StructToYaml(kustomization, kustomizationFilePath)
}
