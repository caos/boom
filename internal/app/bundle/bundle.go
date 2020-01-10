package bundle

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/bundle/application"
	"github.com/caos/boom/internal/app/bundle/bundles"
	"github.com/caos/boom/internal/app/name"
	"github.com/caos/boom/internal/app/templator"
	"github.com/caos/boom/internal/app/templator/helm"
	"github.com/caos/orbiter/logging"
)

var namespace = "caos-system"

type Bundle struct {
	baseDirectoryPath       string
	dashboardsDirectoryPath string
	Applications            map[name.Application]application.Application
	Templator               templator.Templator
	logger                  logging.Logger
}

func New(logger logging.Logger, crdName, baseDirectoryPath, dashboardsDirectoryPath string) *Bundle {
	apps := make(map[name.Application]application.Application, 0)
	templator := templator.New(logger, crdName, baseDirectoryPath, helm.GetName())

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
	return b.Templator.CleanUp()
}

func (b *Bundle) addApplication(appName name.Application) *Bundle {
	app := application.New(b.logger, appName)
	b.Templator.AddApplication(appName, app)
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

	b.Templator.PrepareTemplate(appName, spec)

	deploy := application.Deploy(appName, spec)
	// resultsfilepath := b.Templator.GetResultsFilePath(appName)

	if deploy {
		b.Templator.Template(appName)
		// kubectlCmd := kubectl.New("apply").AddParameter("-f", resultsfilepath).AddParameter("-n", namespace)
		// if err := errors.Wrapf(helper.Run(b.logger, kubectlCmd.Build()), "Failed to apply with file %s", resultsfilepath); err != nil {
		// 	return err
		// }

		app.SetAppliedSpec(spec)
	} else if !deploy && app.Changed(spec) {
		// kubectlCmd := kubectl.New("delete").AddParameter("-f", resultsfilepath).AddParameter("-n", namespace)
		// if err := errors.Wrapf(helper.Run(b.logger, kubectlCmd.Build()), "Failed to delete with file %s", resultsfilepath); err != nil {
		// 	return err
		// }

		app.SetAppliedSpec(nil)
	}

	return nil
}
