package toolset

import "github.com/caos/toolsop/internal/helper"

type Toolsets struct {
	Toolsets []*Toolset `yaml:"Toolsets"`
}

type Toolset struct {
	Name     string     `yaml:"name"`
	Versions []*Version `yaml:"versions"`
}

type Version struct {
	Version      string         `yaml:"version"`
	Applications []*Application `yaml:"applications"`
}

type Application struct {
	Name      string            `yaml:"name"`
	Chart     *Chart            `yaml:"chart"`
	ImageTags map[string]string `yaml:"imageTags"`
}

type Chart struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Index   Index  `yaml:"index"`
}

type Index struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func NewToolsetsFromYaml(path string) (*Toolsets, error) {
	toolsets := &Toolsets{}

	err := helper.YamlToStruct(path, toolsets)
	if err != nil {
		return nil, err
	}
	return toolsets, nil
}
