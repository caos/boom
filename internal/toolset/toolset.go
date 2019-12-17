package toolset

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	"github.com/caos/boom/internal/helper"
)

type Toolsets struct {
	Toolsets []*Toolset `yaml:"Toolsets"`
}

type Toolset struct {
	Name         string         `yaml:"name"`
	Applications []*Application `yaml:"applications"`
}

type Application struct {
	Name string           `yaml:"name"`
	File *ApplicationFile `yaml:"file"`
}
type ApplicationFile struct {
	Chart     *Chart            `yaml:"chart"`
	ImageTags map[string]string `yaml:"imageTags"`
	Crds      []string          `yaml:"crds,omitempty"`
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

func NewToolsetsFromYaml(logger logging.Logger, toolsetsDirectoryPath string) (*Toolsets, error) {
	toolsets, err := getToolsets(logger, toolsetsDirectoryPath)
	if err != nil {
		return nil, err
	}

	return &Toolsets{
		Toolsets: toolsets,
	}, nil
}

func getToolsets(logger logging.Logger, toolsetsDirectoryPath string) ([]*Toolset, error) {
	toolsets := make([]*Toolset, 0)
	toolsetFolders, err := ioutil.ReadDir(toolsetsDirectoryPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read toolsets directory")
	}

	for _, toolsetFolder := range toolsetFolders {
		if toolsetFolder.IsDir() {
			toolsetDirectoryPath := filepath.Join(toolsetsDirectoryPath, toolsetFolder.Name())
			applications, err := getApplications(logger, toolsetDirectoryPath)
			if err != nil {
				return nil, err
			}

			toolset := &Toolset{
				Name:         toolsetFolder.Name(),
				Applications: applications,
			}

			toolsets = append(toolsets, toolset)
		}
	}

	return toolsets, nil
}

func getApplications(logger logging.Logger, versionDirectoryPath string) ([]*Application, error) {

	applications := make([]*Application, 0)

	applicationFiles, err := ioutil.ReadDir(versionDirectoryPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read appplications of version directory")
	}
	for _, applicationFile := range applicationFiles {
		if !applicationFile.IsDir() {
			applicationName := strings.TrimSuffix(applicationFile.Name(), ".yaml")
			applicationName = strings.TrimSuffix(applicationName, ".yml")

			applicationFilePath := filepath.Join(versionDirectoryPath, applicationFile.Name())
			var file ApplicationFile
			if err := errors.Wrap(
				helper.YamlToStruct(applicationFilePath, &file),
				"Failed to marshal application yaml to struct"); err != nil {
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
