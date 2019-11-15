package toolset

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/caos/toolsop/internal/helper"
)

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
	Name string           `yaml:"name"`
	File *ApplicationFile `yaml:"file"`
}
type ApplicationFile struct {
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

func NewToolsetsFromYaml(toolsetsDirectoryPath string) (*Toolsets, error) {
	toolsets, err := getToolsets(toolsetsDirectoryPath)
	if err != nil {
		return nil, err
	}

	return &Toolsets{
		Toolsets: toolsets,
	}, nil
}

func getToolsets(toolsetsDirectoryPath string) ([]*Toolset, error) {
	toolsets := make([]*Toolset, 0)
	toolsetFolders, err := ioutil.ReadDir(toolsetsDirectoryPath)
	if err != nil {
		return nil, err
	}

	for _, toolsetFolder := range toolsetFolders {
		if toolsetFolder.IsDir() {
			toolsetDirectoryPath := filepath.Join(toolsetsDirectoryPath, toolsetFolder.Name())
			versions, err := getVersions(toolsetDirectoryPath)
			if err != nil {
				return nil, err
			}

			toolset := &Toolset{
				Name:     toolsetFolder.Name(),
				Versions: versions,
			}

			toolsets = append(toolsets, toolset)
		}
	}

	return toolsets, nil
}

func getVersions(toolsetDirectoryPath string) ([]*Version, error) {
	versions := make([]*Version, 0)

	versionFolders, err := ioutil.ReadDir(toolsetDirectoryPath)
	if err != nil {
		return nil, err
	}
	for _, versionFolder := range versionFolders {
		if versionFolder.IsDir() {
			versionDirectoryPath := filepath.Join(toolsetDirectoryPath, versionFolder.Name())
			applications, err := getApplications(versionDirectoryPath)
			if err != nil {
				return nil, err
			}

			version := &Version{
				Version:      versionFolder.Name(),
				Applications: applications,
			}

			versions = append(versions, version)
		}
	}

	return versions, nil
}

func getApplications(versionDirectoryPath string) ([]*Application, error) {
	applications := make([]*Application, 0)

	applicationFiles, err := ioutil.ReadDir(versionDirectoryPath)
	if err != nil {
		return nil, err
	}
	for _, applicationFile := range applicationFiles {
		if !applicationFile.IsDir() {
			applicationName := strings.TrimSuffix(applicationFile.Name(), ".yaml")
			applicationName = strings.TrimSuffix(applicationName, ".yml")

			applicationFilePath := filepath.Join(versionDirectoryPath, applicationFile.Name())
			var file ApplicationFile
			if err := helper.YamlToStruct(applicationFilePath, &file); err != nil {
				return nil, err
			}

			application := &Application{
				Name: applicationName,
				File: &file,
			}
			applications = append(applications, application)
		}
	}
	return applications, nil
}
