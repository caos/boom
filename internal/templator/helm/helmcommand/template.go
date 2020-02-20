package helmcommand

import (
	"path/filepath"
	"strings"

	"github.com/caos/boom/internal/helper"
)

type TemplateConfig struct {
	TempFolderPath   string
	ChartName        string
	ReleaseName      string
	ReleaseNamespace string
	ValuesFilePath   string
}

func Template(conf *TemplateConfig) ([]byte, error) {
	var releaseNameParam, releaseNamespaceParam, valuesParam string
	if conf.ReleaseName != "" {
		releaseNameParam = strings.Join([]string{"--name", conf.ReleaseName}, " ")
	}
	if conf.ReleaseNamespace != "" {
		releaseNamespaceParam = strings.Join([]string{"--namespace", conf.ReleaseNamespace}, " ")
	}
	if conf.ValuesFilePath != "" {
		valuesParam = strings.Join([]string{"--values", conf.ValuesFilePath}, " ")
	}

	chartHomeAbs, err := helper.GetAbsPath(conf.TempFolderPath, chartsFolder)
	if err != nil {
		return nil, err
	}
	chartStr := filepath.Join(chartHomeAbs, conf.ChartName)

	command := addIfNotEmpty("template", releaseNameParam)
	command = addIfNotEmpty(command, releaseNamespaceParam)
	command = addIfNotEmpty(command, valuesParam)
	command = addIfNotEmpty(command, chartStr)

	return doHelmCommandOutput(conf.TempFolderPath, command)
}
