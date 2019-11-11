package app

import (
	"os"
	"strings"

	"github.com/caos/toolsop/internal/app/loggingoperator"
	"github.com/caos/toolsop/internal/app/prometheusnodeexporter"
	"github.com/caos/toolsop/internal/app/prometheusoperator"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"

	appcrd "github.com/caos/toolsop/internal/app/crd"
	"github.com/caos/toolsop/internal/git"
	"github.com/caos/toolsop/internal/template"
	"github.com/caos/toolsop/internal/toolset"
)

type App struct {
	Toolsets           *toolset.Toolsets
	ToolsDirectoryPath string
	ToolsGit           *git.Git
	CrdGit             []*appcrd.Crd
	Helms              map[string]*template.Helm
}

func New(toolsDirectoryPath string, toolsetsPath string, toolsUrl string, toolsSecret string) (*App, error) {

	app := &App{
		ToolsDirectoryPath: toolsDirectoryPath,
	}

	g, err := git.New(toolsDirectoryPath, toolsUrl, toolsSecret)
	if err != nil {
		return nil, err
	}
	app.ToolsGit = g

	toolsetsFilePath := strings.Join([]string{toolsDirectoryPath, toolsetsPath}, "/")
	toolsets, err := toolset.NewToolsetsFromYaml(toolsetsFilePath)
	if err != nil {
		return nil, err
	}
	app.Toolsets = toolsets

	app.Helms = make(map[string]*template.Helm, 0)

	return app, nil
}

func (a *App) CleanUp() error {
	for _, g := range a.CrdGit {
		err := g.CleanUp()
		if err != nil {
			return err
		}
	}

	return os.RemoveAll(a.ToolsDirectoryPath)
}

func (a *App) GenerateTemplateComponents(name, crdName, crdVersion string) error {
	_, ok := a.Helms[name]
	if !ok {
		helm, err := template.NewHelm(a.ToolsDirectoryPath, a.Toolsets, crdName, crdVersion, name)
		if err != nil {
			return err
		}
		a.Helms[name] = helm
	}
	return nil
}

func (a *App) Reconcile(name string, crd *toolsetsv1beta1.ToolsetSpec) error {
	helm := a.Helms[name]

	lo := loggingoperator.New(a.ToolsDirectoryPath)
	if err := lo.Reconcile(name, helm, crd.LoggingOperator); err != nil {
		return err
	}

	po := prometheusoperator.New(a.ToolsDirectoryPath)
	if err := po.Reconcile(name, helm, crd.PrometheusOperator); err != nil {
		return err
	}

	pne := prometheusnodeexporter.New(a.ToolsDirectoryPath)
	if err := pne.Reconcile(name, helm, crd.PrometheusNodeExporter); err != nil {
		return err
	}

	return nil
}

func (a *App) AddSupervisedCrd(directoryPath, url, secretPath, crdPath string) error {
	c, err := appcrd.New(directoryPath, url, secretPath, crdPath, a.GenerateTemplateComponents, a.Reconcile)
	if err != nil {
		return err
	}
	a.CrdGit = append(a.CrdGit, c)

	c.Apply()
	return nil
}

func (a *App) MaintainSupervisedCrd() error {
	for _, crdGit := range a.CrdGit {
		err := crdGit.Maintain()
		if err != nil {
			return err
		}
	}
	return nil
}
