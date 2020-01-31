package helper

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func readFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while reading yaml %s", path)
	}

	return data, nil
}

func YamlToStruct(path string, struc interface{}) error {
	data, err := readFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, struc)
	if err != nil {
		return errors.Wrapf(err, "Error while unmarshaling yaml to struct", path)
	}
	return nil
}

func YamlToString(path string) (string, error) {
	data, err := readFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func GetVersionFromYaml(filePath string) (string, error) {
	resource := &Resource{}
	if err := YamlToStruct(filePath, resource); err != nil {
		return "", err
	}

	if resource.ApiVersion == "" {
		return "", errors.New("No attribute apiVersion in yaml")
	}

	parts := strings.Split(resource.ApiVersion, "/")

	return parts[1], nil
}

func GetApiGroupFromYaml(filePath string) (string, error) {
	resource := &Resource{}
	if err := YamlToStruct(filePath, resource); err != nil {
		return "", err
	}

	if resource.ApiVersion == "" {
		return "", errors.New("No attribute apiVersion in yaml")
	}

	parts := strings.Split(resource.ApiVersion, "/")

	return parts[0], nil
}
