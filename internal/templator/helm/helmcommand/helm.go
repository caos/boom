package helmcommand

import (
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	helmHomeFolder string = "helm"
	chartsFolder   string = "charts"
)

func Init(basePath string) error {
	return doHelmCommand(basePath, "init --client-only >& /dev/null")
}

func FetchChart(basePath, name, version, index string) error {

	chartHome := filepath.Join(basePath, chartsFolder)
	chartHomeAbs, err := filepath.Abs(chartHome)
	if err != nil {
		return err
	}

	versionParam := strings.Join([]string{"--version=", version}, "")
	untardirParam := strings.Join([]string{"--untardir=", chartHomeAbs}, "")
	chartStr := strings.Join([]string{index, name}, "/")
	command := strings.Join([]string{"fetch --untar", versionParam, untardirParam, chartStr, ">& /dev/null"}, " ")

	return doHelmCommand(basePath, command)
}

func AddIndex(basePath, indexName, indexURL string) error {

	url := strings.Join([]string{"https://", indexURL}, "")
	command := strings.Join([]string{"repo add", indexName, url, ">& /dev/null"}, " ")

	return doHelmCommand(basePath, command)
}

func RepoUpdate(basePath string) error {
	return doHelmCommand(basePath, "repo update >& /dev/null")
}

func Template(basePath, chartName, releaseName, releaseNamespace, valuesFilePath string) ([]byte, error) {
	var releaseNameParam, releaseNamespaceParam, valuesParam string
	if releaseName != "" {
		releaseNameParam = strings.Join([]string{"--name", releaseName}, " ")
	}
	if releaseNamespace != "" {
		releaseNamespaceParam = strings.Join([]string{"--namespace", releaseNamespace}, " ")
	}
	if valuesFilePath != "" {
		valuesParam = strings.Join([]string{"--values", valuesFilePath}, " ")
	}

	chartHome := filepath.Join(basePath, chartsFolder)
	chartHomeAbs, err := filepath.Abs(chartHome)
	if err != nil {
		return nil, err
	}
	chartStr := strings.Join([]string{chartHomeAbs, chartName}, "/")

	command := addIfNotEmpty("template", releaseNameParam)
	command = addIfNotEmpty(command, releaseNamespaceParam)
	command = addIfNotEmpty(command, valuesParam)
	command = addIfNotEmpty(command, chartStr)

	return doHelmCommandOutput(basePath, command)
}
func addIfNotEmpty(one, two string) string {
	if two != "" {
		return strings.Join([]string{one, two}, " ")
	}
	return one
}

func doHelmCommand(basePath, command string) error {
	helmHomeFolderPath := filepath.Join(basePath, helmHomeFolder)
	helmHomeFolderPathAbs, err := filepath.Abs(helmHomeFolderPath)
	if err != nil {
		return err
	}
	helm := strings.Join([]string{"helm", "--home", helmHomeFolderPathAbs, command}, " ")

	cmd := exec.Command("/bin/sh", "-c", helm)

	return cmd.Run()
}

func doHelmCommandOutput(basePath, command string) ([]byte, error) {
	helmHomeFolderPath := filepath.Join(basePath, helmHomeFolder)
	helmHomeFolderPathAbs, err := filepath.Abs(helmHomeFolderPath)
	if err != nil {
		return nil, err
	}
	helm := strings.Join([]string{"helm", "--home", helmHomeFolderPathAbs, command}, " ")

	cmd := exec.Command("/bin/sh", "-c", helm)

	return cmd.Output()
}
