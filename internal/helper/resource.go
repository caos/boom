package helper

import (
	"strings"
)

type Resource struct {
	Kind       string `yaml:"kind"`
	ApiVersion string `yaml:"apiVersion"`
}

func GetVersionFromYaml(filePath string) (string, error) {
	resource := &Resource{}
	if err := YamlToStruct(filePath, resource); err != nil {
		return "", err
	}
	parts := strings.Split(resource.ApiVersion, "/")

	return parts[1], nil
}
