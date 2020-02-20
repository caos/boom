package fetch

import (
	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/bundle/bundles"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/boom/internal/templator/helm/helmcommand"
	"github.com/caos/orbiter/logging"
)

type ChartKey struct {
	Name    string
	Version string
}

type ChartInfo struct {
	Name      string
	Version   string
	IndexName string
}

func All(logger logging.Logger, basePath string) error {
	allApps := bundles.GetAll()

	logger.Info("Init Helm")

	// helm init to create a HELMHOME
	if err := helmcommand.Init(basePath); err != nil {
		return err
	}

	//indexes in a map so that no doublicates exist
	indexes := make(map[*ChartKey]*chart.Index, 0)
	charts := make([]*ChartInfo, 0)
	logger.Info("Preparing lists of indexes and charts")

	for _, appName := range allApps {
		app := application.New(nil, appName)
		temp, ok := app.(application.HelmApplication)
		// if application doenst implement helm interface then no charts are defined
		if ok {
			// get chartinfo from application
			chart := temp.GetChartInfo()

			// when no index defined then it the helm stable repository
			var indexName string
			if chart.Index != nil {
				indexName = chart.Index.Name
				indexes[&ChartKey{Name: chart.Name, Version: chart.Version}] = chart.Index
			} else {
				indexName = "stable"
			}

			// only add chart if chart is not used by another application, no doublicates
			var found bool
			found = false
			for _, checkChart := range charts {
				if checkChart.Name == chart.Name && checkChart.Version == chart.Version && checkChart.IndexName == indexName {
					found = true
				}
			}
			if !found {
				charts = append(charts, &ChartInfo{Name: chart.Name, Version: chart.Version, IndexName: indexName})
			}
		} else {
			logFields := map[string]interface{}{
				"application": appName.String(),
			}
			logger.WithFields(logFields).Info("Not helm templated")
		}
	}

	logger.Info("Adding all indexes")
	// add all indexes in a map so that no dublicates exist
	for _, v := range indexes {
		if err := addIndex(logger, basePath, v); err != nil {
			return err
		}
	}

	logger.Info("Repo update")
	if err := helmcommand.RepoUpdate(basePath); err != nil {
		return err
	}

	logger.Info("Fetching all charts")
	for _, chart := range charts {
		if err := fetch(logger, basePath, chart); err != nil {
			return err
		}

	}
	return nil
}

func fetch(logger logging.Logger, basePath string, chart *ChartInfo) error {
	logFields := map[string]interface{}{
		"application": chart.Name,
		"version":     chart.Version,
	}
	logFields["indexname"] = chart.IndexName

	logger.WithFields(logFields).Info("Fetching chart")
	return helmcommand.FetchChart(&helmcommand.FetchConfig{
		TempFolderPath: basePath,
		ChartName: chart.Name, 
		ChartVersion: chart.Version, 
		IndexName: chart.IndexName,
	)}
}

func addIndex(logger logging.Logger, basePath string, index *chart.Index) error {
	logFields := map[string]interface{}{
		"index": index.Name,
		"url":   index.URL,
	}
	logger.WithFields(logFields).Info("Adding index")
	return helmcommand.AddIndex(&helmcommand.IndexConfig{
		TempFolderPath: basePath,
		IndexName:      index.Name,
		IndexURL:       index.URL,
	})
}
