package app

import (
	"github.com/caos/orbiter/logging"
	"k8s.io/apimachinery/pkg/runtime"
)

type App struct {
	ToolsDirectoryPath      string
	CrdDirectoryPath        string
	DashboardsDirectoryPath string
	GitCrds                 []GitCrd
	Crds                    map[string]Crd
	logger                  logging.Logger
}

func New(logger logging.Logger, toolsDirectoryPath, crdDirectoryPath, dashboardsDirectoryPath string) (*App, error) {

	app := &App{
		ToolsDirectoryPath:      toolsDirectoryPath,
		CrdDirectoryPath:        crdDirectoryPath,
		DashboardsDirectoryPath: dashboardsDirectoryPath,
		logger:                  logger,
	}

	app.Crds = make(map[string]Crd, 0)
	app.GitCrds = make([]GitCrd, 0)

	return app, nil
}

func (a *App) CleanUp() error {

	a.logger.WithFields(map[string]interface{}{
		"logID": "APP-GiK5XPA5PzwQtjR",
	}).Info("Cleanup")

	for _, g := range a.GitCrds {
		err := g.CleanUp()
		if err != nil {
			return err
		}
	}

	for _, c := range a.Crds {
		err := c.CleanUp()
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) AddGitCrd(url string, privateKey []byte, crdPath string) error {
	c, err := NewGitCrd(a.logger, a.CrdDirectoryPath, url, privateKey, crdPath, a.ToolsDirectoryPath, a.DashboardsDirectoryPath)
	if err != nil {
		return err
	}
	a.GitCrds = append(a.GitCrds, c)
	return nil
}

func (a *App) ReconcileGitCrds() error {
	for _, crdGit := range a.GitCrds {
		a.logger.WithFields(map[string]interface{}{
			"logID": "APP-aZAeIqcAmHzflSB",
		}).Info("Started reconciling of GitCRDs")

		err := crdGit.Reconcile()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) ReconcileCrd(version, namespacedName string, getToolsetCRD func(instance runtime.Object) error) error {
	a.logger.WithFields(map[string]interface{}{
		"logID": "APP-aZAeIqcAmHzflSB",
		"name":  namespacedName,
	}).Info("Started reconciling of CRD")

	var err error
	crd, ok := a.Crds[namespacedName]
	if !ok {
		crd, err = NewCrd(a.logger, version, getToolsetCRD, a.ToolsDirectoryPath, a.DashboardsDirectoryPath)
		if err != nil {
			return err
		}

		a.Crds[namespacedName] = crd
		return nil
	}

	return crd.ReconcileWithFunc(getToolsetCRD)
}
