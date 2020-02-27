package bundle

import (
	"sync"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/bundle/bundles"
	"github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/boom/internal/templator/helm"
	helperTemp "github.com/caos/boom/internal/templator/helper"
	"github.com/caos/boom/internal/templator/yaml"
	"github.com/caos/orbiter/mntr"
	"github.com/pkg/errors"
)

var (
	Testmode bool = false
)

type Bundle struct {
	baseDirectoryPath string
	crdName           string
	Applications      map[name.Application]application.Application
	HelmTemplator     templator.Templator
	YamlTemplator     templator.Templator
	monitor           mntr.Monitor
	status            error
}

func New(conf *config.Config) *Bundle {
	apps := make(map[name.Application]application.Application, 0)
	helmTemplator := helperTemp.NewTemplator(conf.Monitor, conf.CrdName, conf.BaseDirectoryPath, helm.GetName())
	yamlTemplator := helperTemp.NewTemplator(conf.Monitor, conf.CrdName, conf.BaseDirectoryPath, yaml.GetName())

	b := &Bundle{
		crdName:           conf.CrdName,
		baseDirectoryPath: conf.BaseDirectoryPath,
		monitor:           conf.Monitor,
		HelmTemplator:     helmTemplator,
		YamlTemplator:     yamlTemplator,
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

	b.status = b.HelmTemplator.CleanUp().GetStatus()
	if b.GetStatus() != nil {
		return b
	}
	b.status = b.YamlTemplator.CleanUp().GetStatus()
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

	app := application.New(b.monitor, appName)
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
	// and this orderNumber by orderNumber (default is 1)
	for orderNumber := 0; applicationCount < len(b.Applications); orderNumber++ {
		var wg sync.WaitGroup
		for appName := range b.Applications {
			//if application has the same orderNumber as currently iterating the reconcile the application
			if application.GetOrderNumber(appName) == orderNumber {
				wg.Add(1)
				b.ReconcileApplication(appName, spec, &wg)
				if err := b.GetStatus(); err != nil {
					b.status = err
					return b
				}
				applicationCount++
			}
		}
		wg.Wait()
	}

	return b
}

func (b *Bundle) ReconcileApplication(appName name.Application, spec *v1beta1.ToolsetSpec, wg *sync.WaitGroup) *Bundle {
	defer wg.Done()

	if b.status != nil {
		return b
	}

	logFields := map[string]interface{}{
		"application": appName,
		"action":      "reconciling",
	}
	monitor := b.monitor.WithFields(logFields)

	resourceInfoList, err := clientgo.GetGroupVersionsResources([]string{})
	if err != nil {
		b.status = err
		monitor.Error(b.status)
		return b
	}

	app, found := b.Applications[appName]
	if !found {
		b.status = errors.New("Application not found")
		monitor.Error(b.status)
		return b
	}
	monitor.Info("Start")

	deploy := app.Deploy(spec)

	var resultFunc func(string, string) error
	if Testmode {
		resultFunc = func(resultFilePath, namespace string) error {
			return nil
		}
	} else {
		if deploy {
			resultFunc = applyWithCurrentState(monitor, resourceInfoList, app)
		} else {
			resultFunc = deleteWithCurrentState(monitor, resourceInfoList, app)
		}
	}

	_, usedHelm := app.(application.HelmApplication)
	if usedHelm {
		b.status = b.HelmTemplator.Template(app, spec, resultFunc).GetStatus()
		if b.status != nil {
			return b
		}
	}
	_, usedYaml := app.(application.YAMLApplication)
	if usedYaml {
		b.status = b.YamlTemplator.Template(app, spec, resultFunc).GetStatus()
	}

	monitor.Info("Done")
	return b
}
