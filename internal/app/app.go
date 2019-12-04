package app

import (
	"path/filepath"

	"github.com/caos/orbiter/logging"
	"github.com/caos/boom/internal/toolset"
	"k8s.io/apimachinery/pkg/runtime"
)

type App struct {
	Toolsets           *toolset.Toolsets
	ToolsDirectoryPath string
	CrdDirectoryPath   string
	GitCrds            []GitCrd
	Crds               map[string]Crd
	logger             logging.Logger
}

func New(logger logging.Logger, toolsDirectoryPath, crdDirectoryPath, toolsetsPath string) (*App, error) {

	app := &App{
		ToolsDirectoryPath: toolsDirectoryPath,
		CrdDirectoryPath:   crdDirectoryPath,
		logger:             logger,
	}

	app.Crds = make(map[string]Crd, 0)
	app.GitCrds = make([]GitCrd, 0)

	err := app.ReloadCurrentToolsets(app.ToolsDirectoryPath, toolsetsPath)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) ReloadCurrentToolsets(toolsDirectoryPath string, toolsetsPath string) error {
	toolsetsFilePath := filepath.Join(toolsDirectoryPath, toolsetsPath)
	toolsets, err := toolset.NewToolsetsFromYaml(a.logger, toolsetsFilePath)
	if err != nil {
		return err
	}
	a.Toolsets = toolsets
	return nil
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

func (a *App) AddGitCrd(url, secretPath, crdPath string) error {
	c, err := NewGitCrd(a.logger, a.CrdDirectoryPath, url, secretPath, crdPath, a.ToolsDirectoryPath, a.Toolsets)
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
		err := crdGit.Reconcile(a.ToolsDirectoryPath, a.Toolsets)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) ReconcileCrd(version, namespacedName string, getToolset func(obj runtime.Object) error) error {
	a.logger.WithFields(map[string]interface{}{
		"logID": "APP-aZAeIqcAmHzflSB",
		"name":  namespacedName,
	}).Info("Started reconciling of CRD")
	crd, ok := a.Crds[namespacedName]
	if !ok {
		newCrd, err := NewCrd(a.logger, version, getToolset, a.ToolsDirectoryPath, a.Toolsets)
		if err != nil {
			return err
		}

		a.Crds[namespacedName] = newCrd
	} else {
		if err := crd.ReconcileWithFunc(getToolset, a.ToolsDirectoryPath, a.Toolsets); err != nil {
			return err
		}
	}
	return nil
}
