package helper

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

func DeleteKindFromYaml(path string, kind string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	text := string(content)
	parts := strings.Split(text, "\n---\n")
	if err := os.Remove(path); err != nil {
		return err
	}

	for _, part := range parts {
		struc := &Resource{}
		if err := yaml.Unmarshal([]byte(part), struc); err != nil {
			return err
		}

		if struc.Kind != kind {
			if err := AddStringObjectToYaml(path, part); err != nil {
				return err
			}
		}
	}

	return nil
}

func DeleteFirstResourceFromYaml(path, apiVersion, kind, name string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	text := string(content)
	parts := strings.Split(text, "\n---\n")
	if err := os.Remove(path); err != nil {
		return err
	}

	found := false
	for _, part := range parts {
		struc := &Resource{}
		if err := yaml.Unmarshal([]byte(part), struc); err != nil {
			return err
		}

		if found || !(struc.ApiVersion == apiVersion && struc.Kind == kind && struc.Metadata.Name == name) {
			if err := AddStringObjectToYaml(path, part); err != nil {
				return err
			}
		} else {
			found = true
		}
	}

	return nil
}
