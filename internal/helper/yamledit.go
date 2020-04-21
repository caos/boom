package helper

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
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
	if FileExists(path) {
		if _, err := f.WriteString("\n---\n"); err != nil {
			return err
		}
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
	exists := false
	if FileExists(path) {
		exists = true
	}

	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	if exists {
		if _, err := f.WriteString("\n---\n"); err != nil {
			return err
		}
	}

	if _, err := f.WriteString(str); err != nil {
		return err
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
		if part == "" {
			continue
		}
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

func ReplacePointInclusiveFiller(filePath, point, content, filler string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	text := string(data)

	os.Remove(filePath)

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.Contains(line, point) {
			line := lines[i]
			lines[i] = strings.Replace(line, point, content, 1)
		}
		if strings.Contains(line, filler) {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}
	output := strings.Join(lines, "\n")

	if err := AddStringObjectToYaml(filePath, output); err != nil {
		return err
	}
	return nil
}
