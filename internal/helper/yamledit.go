package helper

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Name string `yaml:"name"`
}
type Resource struct {
	Kind       string    `yaml:"kind"`
	ApiVersion string    `yaml:"apiVersion"`
	Metadata   *Metadata `yaml:"metadata"`
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

	if _, err := f.WriteString(str); err != nil {
		return err
	}
	return nil
}

func AddStringObjectToYaml(path string, str string) error {
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

func AddYamlToYaml(filePath, addFilePath string) error {

	f, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	addContent, err := ioutil.ReadFile(addFilePath)
	if err != nil {
		return err
	}
	addText := string(addContent)

	if _, err := f.WriteString("\n---\n"); err != nil {
		return err
	}
	if _, err := f.WriteString(addText); err != nil {
		return err
	}
	return nil
}

func DeleteKindFromYaml(path string, kind string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	text := string(content)
	parts := strings.Split(text, "\n---\n")

	os.Remove(path)
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

func AddStringBeforePointForKindAndName(filePath, kind, name, point, addContent string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	text := string(content)
	parts := strings.Split(text, "\n---\n")

	os.Remove(filePath)
	for _, part := range parts {
		struc := &Resource{}
		if err := yaml.Unmarshal([]byte(part), struc); err != nil {
			return err
		}
		output := part
		if struc.Kind == kind && struc.Metadata.Name == name {
			lines := strings.Split(part, "\n")
			for i, line := range lines {
				if strings.Contains(line, point) {
					lines[i] = strings.Join([]string{addContent, line}, "")
				}
			}
			output = strings.Join(lines, "\n")
		}

		if err := AddStringObjectToYaml(filePath, output); err != nil {
			return err
		}
	}
	return nil
}
