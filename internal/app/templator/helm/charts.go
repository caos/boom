package helm

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/caos/boom/internal/app/bundle/application"
	"github.com/caos/boom/internal/app/bundle/application/chart"
	"github.com/caos/boom/internal/app/bundle/bundles"
)

func (h *Helm) FetchAllCharts() error {
	allApps := bundles.GetAll()

	logFields := map[string]interface{}{
		"logID": "HELM-Ay5T4nn9kKWTSWU",
	}
	h.logger.WithFields(logFields).Info("Init Helm")
	h.doHelmCommand("init --client-only >& /dev/null")

	h.logger.WithFields(logFields).Info("Fetching all charts")
	for _, appName := range allApps {
		app := application.New(nil, appName)
		temp, ok := app.(Templator)
		if ok {
			chart := temp.GetChartInfo()
			if err := h.fetchChart(chart); err != nil {
				return err
			}
		} else {
			logFields := map[string]interface{}{
				"logID":       "HELM-Ay5T4nn9kKWTSWU",
				"application": appName.String(),
			}
			h.logger.WithFields(logFields).Info("Not helm templated")
		}
	}
	return nil
}

func (h *Helm) fetchChart(chart *chart.Chart) error {
	chartHome := filepath.Join(h.templatorDirectoryPath, chartsFolder)
	chartHomeAbs, err := filepath.Abs(chartHome)
	if err != nil {
		return err
	}

	var indexname string
	if chart.Index != nil {
		if err := h.addIndex(chart.Index); err != nil {
			return err
		}
		indexname = chart.Index.Name
	} else {
		indexname = "stable"
	}

	if err := h.doHelmCommand("repo update >& /dev/null"); err != nil {
		return err
	}

	version := strings.Join([]string{"--version=", chart.Version}, "")
	untardir := strings.Join([]string{"--untardir=", chartHomeAbs}, "")
	chartStr := strings.Join([]string{indexname, chart.Name}, "/")
	command := strings.Join([]string{"fetch --untar", version, untardir, chartStr, ">& /dev/null"}, " ")

	logFields := map[string]interface{}{
		"application": chart.Name,
		"version":     chart.Version,
		"logID":       "HELM-HkLTnAhAJnAyPq8",
	}
	h.logger.WithFields(logFields).Info("Fetching chart")
	return h.doHelmCommand(command)
}

func (h *Helm) addIndex(index *chart.Index) error {

	url := strings.Join([]string{"https://", index.URL}, "")
	command := strings.Join([]string{"repo add", index.Name, url, ">& /dev/null"}, " ")

	logFields := map[string]interface{}{
		"index": index.Name,
		"url":   index.URL,
		"logID": "HELM-RsDBpIgtiJkgQVs",
	}
	h.logger.WithFields(logFields).Info("Adding index")
	return h.doHelmCommand(command)
}

func (h *Helm) doHelmCommand(command string) error {
	helmHomeFolderPath := filepath.Join(h.templatorDirectoryPath, helmHomeFolder)
	helmHomeFolderPathAbs, err := filepath.Abs(helmHomeFolderPath)
	if err != nil {
		return err
	}
	helm := strings.Join([]string{"helm", "--home", helmHomeFolderPathAbs, command}, " ")

	cmd := exec.Command("/bin/sh", "-c", helm)

	return cmd.Run()
}
