package app

import (
	"os"
	"strings"

	"github.com/caos/toolsop/internal/git"
	"github.com/caos/toolsop/internal/toolset"
	"k8s.io/apimachinery/pkg/runtime"
)

type App struct {
	Toolsets           *toolset.Toolsets
	ToolsDirectoryPath string
	CrdDirectoryPath   string
	ToolsGit           *git.Git
	GitCrds            []GitCrd
	Crds               map[string]Crd
}

func New(toolsDirectoryPath, crdDirectoryPath, toolsetsPath, toolsUrl, toolsSecret string) (*App, error) {

	app := &App{
		ToolsDirectoryPath: toolsDirectoryPath,
		CrdDirectoryPath:   crdDirectoryPath,
	}

	g, err := git.New(toolsDirectoryPath, toolsUrl, toolsSecret)
	if err != nil {
		return nil, err
	}
	app.ToolsGit = g

	app.Crds = make(map[string]Crd, 0)
	app.GitCrds = make([]GitCrd, 0)

	err = app.ReloadCurrentToolsets(app.ToolsDirectoryPath, toolsetsPath)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) ReloadCurrentToolsets(toolsDirectoryPath string, toolsetsPath string) error {
	toolsetsFilePath := strings.Join([]string{toolsDirectoryPath, toolsetsPath}, "/")
	toolsets, err := toolset.NewToolsetsFromYaml(toolsetsFilePath)
	if err != nil {
		return err
	}
	a.Toolsets = toolsets
	return nil
}

func (a *App) CleanUp() error {
	for _, g := range a.GitCrds {
		err := g.CleanUp()
		if err != nil {
			return err
		}
	}

	return os.RemoveAll(a.ToolsDirectoryPath)
}

func (a *App) AddGitCrd(version, url, secretPath, crdPath string) error {
	c, err := NewGitCrd(version, a.CrdDirectoryPath, url, secretPath, crdPath, a.ToolsDirectoryPath, a.Toolsets)
	if err != nil {
		return err
	}
	a.GitCrds = append(a.GitCrds, c)
	return nil
}

func (a *App) ReconcileGitCrds(version string) error {
	for _, crdGit := range a.GitCrds {
		err := crdGit.Reconcile(a.ToolsDirectoryPath, a.Toolsets)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) ReconcileCrd(version, namespacedName string, getToolset func(obj runtime.Object) error) error {
	crd, ok := a.Crds[namespacedName]
	if !ok {
		newCrd, err := NewCrd(version, getToolset, a.ToolsDirectoryPath, a.Toolsets)
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

func (a *App) GetCrdDirectoryPath() string {
	return a.CrdDirectoryPath
}

func (a *App) GetToolsDirectoryPath() string {
	return a.CrdDirectoryPath
}
