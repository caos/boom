package helper

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func StructToYaml(struc interface{}, path string) error {

	data, err := yaml.Marshal(struc)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

func YamlToStruct(path string, struc interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, struc)
}
