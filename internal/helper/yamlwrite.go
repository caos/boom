package helper

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func StructToYaml(struc interface{}, path string) error {

	data, err := yaml.Marshal(struc)
	if err != nil {
		return errors.Wrapf(err, "Error while marshaling struct to byte")
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return errors.Wrapf(err, "Error while writing struct to yaml %s", path)
	}
	return nil
}
