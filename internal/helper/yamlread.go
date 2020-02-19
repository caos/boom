package helper

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
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
		return errors.Wrapf(err, "Error while unmarshaling yaml %s to struct", path)
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
