package app

import (
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/crd"
	"github.com/caos/boom/internal/current"
	"github.com/caos/boom/internal/gitcrd"
	gitcrdconfig "github.com/caos/boom/internal/gitcrd/config"

	"github.com/caos/boom/internal/bundle/bundles"
	"github.com/caos/boom/internal/templator/helm"
	"github.com/caos/orbiter/mntr"
)

type App struct {
	ToolsDirectoryPath string
	GitCrds            []gitcrd.GitCrd
	Crds               map[string]crd.Crd
	monitor            mntr.Monitor
}

func New(monitor mntr.Monitor, toolsDirectoryPath, dashboardsDirectoryPath string) (*App, error) {

	app := &App{
		ToolsDirectoryPath: toolsDirectoryPath,
		monitor:            monitor,
	}

	app.Crds = make(map[string]crd.Crd, 0)
	app.GitCrds = make([]gitcrd.GitCrd, 0)

	return app, nil
}

func (a *App) CleanUp() error {

	a.monitor.Info("Cleanup")

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

func (a *App) AddGitCrd(gitCrdConf *gitcrdconfig.Config) error {
	c, err := gitcrd.New(gitCrdConf)
	if err != nil {
		return err
	}

	bundleConf := &bundleconfig.Config{
		BundleName:        bundles.Caos,
		BaseDirectoryPath: a.ToolsDirectoryPath,
		Templator:         helm.GetName(),
	}

	c.SetBundle(bundleConf)
	if err := c.GetStatus(); err != nil {
		return err
	}

	a.GitCrds = append(a.GitCrds, c)
	return nil
}

func (a *App) getCurrent(monitor mntr.Monitor) ([]*clientgo.Resource, error) {

	resourceInfoList, err := clientgo.GetGroupVersionsResources([]string{})
	if err != nil {
		monitor.Error(err)
		return nil, err
	}

	return current.Get(a.monitor, resourceInfoList), nil
}

func (a *App) ReconcileGitCrds() error {
	monitor := a.monitor.WithFields(map[string]interface{}{
		"action": "reconciling",
	})
	monitor.Info("Started reconciling of GitCRDs")

	for _, crdGit := range a.GitCrds {
		crdGit.SetBackStatus()

		currentResourceList, err := a.getCurrent(monitor)
		if err != nil {
			return err
		}

		crdGit.Reconcile(currentResourceList)
		if err := crdGit.GetStatus(); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) WriteBackCurrentState() error {

	monitor := a.monitor.WithFields(map[string]interface{}{
		"action": "current",
	})
	monitor.Info("Started writeback of currentstate of GitCRDs")

	for _, crdGit := range a.GitCrds {
		crdGit.SetBackStatus()

		currentResourceList, err := a.getCurrent(monitor)
		if err != nil {
			return err
		}

		crdGit.WriteBackCurrentState(currentResourceList)
		if err := crdGit.GetStatus(); err != nil {
			return err
		}
	}
	return nil
}

// func (a *App) ReconcileCrd(version, namespacedName string, getToolsetCRD func(instance runtime.Object) error) error {
// 	a.monitor.WithFields(map[string]interface{}{
// 		"name": namespacedName,
// 	}).Info("Started reconciling of CRD")

// 	var err error
// 	managedcrd, ok := a.Crds[namespacedName]
// 	if !ok {
// 		crdConf := &crdconfig.Config{
// 			Monitor:  a.monitor,
// 			Version: v1beta1.GetVersion(),
// 		}

// 		managedcrd, err = crd.New(crdConf)
// 		if err != nil {
// 			return err
// 		}

// 		bundleConf := &bundleconfig.Config{
// 			Monitor:            a.monitor,
// 			CrdName:           namespacedName,
// 			BundleName:        bundles.Caos,
// 			BaseDirectoryPath: a.ToolsDirectoryPath,
// 			Templator:         helm.GetName(),
// 		}
// 		managedcrd.SetBundle(bundleConf)

// 		if err := managedcrd.GetStatus(); err != nil {
// 			return err
// 		}

// 		a.Crds[namespacedName] = managedcrd
// 	}

// 	managedcrd.ReconcileWithFunc(getToolsetCRD)
// 	return managedcrd.GetStatus()
// }
