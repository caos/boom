package fetch

import (
	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/bundle/bundles"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/boom/internal/templator/helm/helmcommand"
	"github.com/caos/orbiter/logging"
)

func All(logger logging.Logger, basePath string) error {
	allApps := bundles.GetAll()

	logFields := map[string]interface{}{
		"logID": "HELM-Ay5T4nn9kKWTSWU",
	}
	logger.WithFields(logFields).Info("Init Helm")
	if err := helmcommand.Init(basePath); err != nil {
		return err
	}

	logger.WithFields(logFields).Info("Fetching all charts")
	for _, appName := range allApps {
		app := application.New(nil, appName)
		temp, ok := app.(application.HelmApplication)
		if ok {
			chart := temp.GetChartInfo()
			if err := fetch(logger, basePath, chart); err != nil {
				return err
			}
		} else {
			logFields := map[string]interface{}{
				"logID":       "HELM-Ay5T4nn9kKWTSWU",
				"application": appName.String(),
			}
			logger.WithFields(logFields).Info("Not helm templated")
		}
	}
	return nil

}

func fetch(logger logging.Logger, basePath string, chart *chart.Chart) error {

	var indexname string
	if chart.Index != nil {
		if err := addIndex(logger, basePath, chart.Index); err != nil {
			return err
		}
		indexname = chart.Index.Name
	} else {
		indexname = "stable"
	}

	if err := helmcommand.RepoUpdate(basePath); err != nil {
		return err
	}

	logFields := map[string]interface{}{
		"application": chart.Name,
		"version":     chart.Version,
		"logID":       "HELM-HkLTnAhAJnAyPq8",
	}
	logger.WithFields(logFields).Info("Fetching chart")
	return helmcommand.FetchChart(basePath, chart.Name, chart.Version, indexname)
}

func addIndex(logger logging.Logger, basePath string, index *chart.Index) error {
	logFields := map[string]interface{}{
		"index": index.Name,
		"url":   index.URL,
		"logID": "HELM-RsDBpIgtiJkgQVs",
	}
	logger.WithFields(logFields).Info("Adding index")
	return helmcommand.AddIndex(basePath, index.Name, index.URL)
}
