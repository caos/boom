package helper

import (
	"io/ioutil"
	"os"
	"strings"

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

func AddStructToYaml(path string, struc interface{}) error {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	data, err := yaml.Marshal(struc)
	if err != nil {
		return err
	}

	if _, err := f.WriteString("\n---\n"); err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

func AddStringToYaml(path string, str string) error {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.WriteString("\n---\n"); err != nil {
		return err
	}
	if _, err := f.WriteString(str); err != nil {
		return err
	}
	return nil
}

func DeletePartOfYaml(path string, str string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	text := string(content)
	parts := strings.Split(text, "\n---\n")

	os.Remove(path)
	for _, part := range parts {
		if !strings.Contains(part, str) {
			if err := AddStringToYaml(path, part); err != nil {
				return err
			}
		}
	}

	return nil
}
