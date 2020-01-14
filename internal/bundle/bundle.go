package bundle

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/boom/internal/templator/helm"
	helperTemp "github.com/caos/boom/internal/templator/helper"
	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/bundle/bundles"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

type Bundle struct {
	baseDirectoryPath       string
	dashboardsDirectoryPath string
	Applications            map[name.Application]application.Application
	Templator               templator.Templator
	logger                  logging.Logger
}

func New(logger logging.Logger, crdName, baseDirectoryPath, dashboardsDirectoryPath string) *Bundle {
	apps := make(map[name.Application]application.Application, 0)
	templator := helperTemp.NewTemplator(logger, crdName, baseDirectoryPath, helm.GetName())

	b := &Bundle{
		baseDirectoryPath:       baseDirectoryPath,
		dashboardsDirectoryPath: dashboardsDirectoryPath,
		logger:                  logger,
		Templator:               templator,
		Applications:            apps,
	}

	basis := bundles.GetBasisset()
	for _, app := range basis {
		b.addApplication(app)
	}
	return b
}

func (b *Bundle) CleanUp() error {
	return b.Templator.CleanUp().GetStatus()
}

func (b *Bundle) addApplication(appName name.Application) *Bundle {
	app := application.New(b.logger, appName)
	b.Applications[appName] = app
	return b
}

func (b *Bundle) Reconcile(spec *v1beta1.ToolsetSpec) error {
	for appName := range b.Applications {
		if err := b.ReconcileApplication(appName, spec); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bundle) ReconcileApplication(appName name.Application, spec *v1beta1.ToolsetSpec) error {
	logFields := map[string]interface{}{
		"application": appName,
		"logID":       "CRD-rGkpjHLZtVAWumr",
	}

	app := b.Applications[appName]

	b.logger.WithFields(logFields).Info("Reconciling")

	deploy := application.Deploy(appName, spec)
	var command string
	if deploy {
		command = "apply"
	} else if !deploy && app.Changed(spec) && !app.Initial() {
		command = "delete"
	}

	resultFunc := func(resultFilePath string) error {
		kubectlCmd := kubectl.New(command).AddParameter("-f", resultFilePath).AddParameter("-n", "caos-system")
		return errors.Wrapf(helper.Run(b.logger, kubectlCmd.Build()), "Failed to apply with file %s", resultFilePath)
	}
	if command == "" {
		resultFunc = func(resultFilePath string) error { return nil }
	}

	err := b.Templator.Template(app, spec, resultFunc).GetStatus()
	if err != nil {
		return err
	}

	app.SetAppliedSpec(spec)
	return nil
}
