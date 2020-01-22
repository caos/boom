package app

import (
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/crd"
	crdconfig "github.com/caos/boom/internal/crd/config"
	"github.com/caos/boom/internal/crd/v1beta1"
	"github.com/caos/boom/internal/gitcrd"
	gitcrdconfig "github.com/caos/boom/internal/gitcrd/config"

	"github.com/caos/boom/internal/bundle/bundles"
	"github.com/caos/boom/internal/templator/helm"
	"github.com/caos/orbiter/logging"
	"k8s.io/apimachinery/pkg/runtime"
)

type App struct {
	ToolsDirectoryPath string
	CrdDirectoryPath   string
	GitCrds            []gitcrd.GitCrd
	Crds               map[string]crd.Crd
	logger             logging.Logger
}

func New(logger logging.Logger, toolsDirectoryPath, crdDirectoryPath, dashboardsDirectoryPath string) (*App, error) {

	app := &App{
		ToolsDirectoryPath: toolsDirectoryPath,
		CrdDirectoryPath:   crdDirectoryPath,
		logger:             logger,
	}

	app.Crds = make(map[string]crd.Crd, 0)
	app.GitCrds = make([]gitcrd.GitCrd, 0)

	return app, nil
}

func (a *App) CleanUp() error {

	a.logger.WithFields(map[string]interface{}{
		"logID": "APP-GiK5XPA5PzwQtjR",
	}).Info("Cleanup")

	for _, g := range a.GitCrds {
		g.CleanUp()

		if err := g.GetStatus(); err != nil {
			return err
		}
	}

	for _, c := range a.Crds {
		c.CleanUp()

		if err := c.GetStatus(); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) AddGitCrd(url string, privateKey []byte, crdPath string) error {

	gitcrdConf := &gitcrdconfig.Config{
		Logger:           a.logger,
		CrdDirectoryPath: a.CrdDirectoryPath,
		CrdUrl:           url,
		PrivateKey:       privateKey,
		CrdPath:          crdPath,
	}

	c, err := gitcrd.New(gitcrdConf)
	if err != nil {
		return err
	}

	toolsetCRD, err := c.GetCrdContent()
	if err != nil {
		return err
	}

	bundleConf := &bundleconfig.Config{
		Logger:            a.logger,
		CrdName:           toolsetCRD.Name,
		BundleName:        bundles.Caos,
		BaseDirectoryPath: a.ToolsDirectoryPath,
		Templator:         helm.GetName(),
	}

	c.SetBundle(bundleConf)
	a.GitCrds = append(a.GitCrds, c)
	return nil
}

func (a *App) ReconcileGitCrds() error {
	for _, crdGit := range a.GitCrds {
		a.logger.WithFields(map[string]interface{}{
			"logID": "APP-aZAeIqcAmHzflSB",
		}).Info("Started reconciling of GitCRDs")

		crdGit.Reconcile()
		err := crdGit.GetStatus()
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
	managedcrd, ok := a.Crds[namespacedName]
	if !ok {
		crdConf := &crdconfig.Config{
			Logger:  a.logger,
			Version: v1beta1.GetVersion(),
		}

		managedcrd, err = crd.New(crdConf)
		if err != nil {
			return err
		}

		bundleConf := &bundleconfig.Config{
			Logger:            a.logger,
			CrdName:           namespacedName,
			BundleName:        bundles.Caos,
			BaseDirectoryPath: a.ToolsDirectoryPath,
			Templator:         helm.GetName(),
		}
		managedcrd.SetBundle(bundleConf)

		if err := managedcrd.GetStatus(); err != nil {
			return err
		}

		a.Crds[namespacedName] = managedcrd
	}

	managedcrd.ReconcileWithFunc(getToolsetCRD)
	return managedcrd.GetStatus()
}
