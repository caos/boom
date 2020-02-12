package bundle

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/bundle/bundles"
	"github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator"
	helperTemp "github.com/caos/boom/internal/templator/helper"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

type Bundle struct {
	baseDirectoryPath string
	Applications      map[name.Application]application.Application
	Templator         templator.Templator
	logger            logging.Logger
	status            error
}

func New(conf *config.Config) *Bundle {
	apps := make(map[name.Application]application.Application, 0)
	templator := helperTemp.NewTemplator(conf.Logger, conf.CrdName, conf.BaseDirectoryPath, conf.Templator)

	b := &Bundle{
		baseDirectoryPath: conf.BaseDirectoryPath,
		logger:            conf.Logger,
		Templator:         templator,
		Applications:      apps,
		status:            nil,
	}
	return b
}

func (b *Bundle) GetStatus() error {
	return b.status
}

func (b *Bundle) CleanUp() *Bundle {
	if b.GetStatus() != nil {
		return b
	}

	b.status = b.Templator.CleanUp().GetStatus()
	return b
}

func (b *Bundle) GetApplications() map[name.Application]application.Application {
	return b.Applications
}

func (b *Bundle) AddApplicationsByBundleName(name name.Bundle) error {
	if err := b.GetStatus(); err != nil {
		return err
	}

	names := bundles.Get(name)
	if names == nil {
		return errors.Errorf("No bundle known with name %s", name)
	}

	bnew := b
	for _, name := range names {
		bnew = bnew.AddApplicationByName(name)
		if err := bnew.GetStatus(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bundle) AddApplicationByName(appName name.Application) *Bundle {
	if b.GetStatus() != nil {
		return b
	}

	app := application.New(b.logger, appName)
	return b.AddApplication(app)
}

func (b *Bundle) AddApplication(app application.Application) *Bundle {
	if b.GetStatus() != nil {
		return b
	}

	if _, found := b.Applications[app.GetName()]; found {
		b.status = errors.New("Application already in bundle")
		return b
	}

	b.Applications[app.GetName()] = app
	return b
}

func (b *Bundle) Reconcile(spec *v1beta1.ToolsetSpec) *Bundle {
	applicationCount := 0
	// go through list of application until every application is reconciled
	// and this orderNumber by orderNumber (default is 0)
	for orderNumber := 0; applicationCount < len(b.Applications); orderNumber++ {
		for appName := range b.Applications {
			//if application has the same orderNumber as currently iterating the reconcile the application
			if application.GetOrderNumber(appName) == orderNumber {
				if err := b.ReconcileApplication(appName, spec).GetStatus(); err != nil {
					b.status = err
					return b
				}
				applicationCount++
			}
		}
	}

	return b
}

func (b *Bundle) ReconcileApplication(appName name.Application, spec *v1beta1.ToolsetSpec) *Bundle {
	if b.status != nil {
		return b
	}

	logFields := map[string]interface{}{
		"application": appName,
		"logID":       "CRD-rGkpjHLZtVAWumr",
	}

	app, found := b.Applications[appName]
	if !found {
		b.status = errors.New("Application not found")
		b.logger.WithFields(logFields).Error(b.status)
		return b
	}
	b.logger.WithFields(logFields).Info("Reconciling")

	deploy := app.Deploy(spec)
	var command string
	if deploy {
		command = "apply"
	} else if !deploy && app.Changed(spec) && !app.Initial() {
		command = "delete"
	}

	resultFunc := func(resultFilePath, namespace string) error {
		kubectlCmd := kubectl.New(command).AddParameter("-f", resultFilePath).AddParameter("-n", namespace)
		return errors.Wrapf(helper.Run(b.logger, kubectlCmd.Build()), "Failed to apply with file %s", resultFilePath)
	}

	if command == "" {
		resultFunc = func(resultFilePath, namespace string) error { return nil }
	}

	b.status = b.Templator.Template(app, spec, resultFunc).GetStatus()
	if b.status != nil {
		return b
	}

	app.SetAppliedSpec(spec)
	return b
}
