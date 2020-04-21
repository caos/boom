package helper

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
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

	text := string(data)
	parts := strings.Split(text, "\n---\n")
	for _, part := range parts {
		if part == "" {
			continue
		}
		err = yaml.Unmarshal([]byte(part), struc)
		if err != nil {
			return errors.Wrapf(err, "Error while unmarshaling yaml %s to struct", path)
		}
		return nil
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
