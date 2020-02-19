package bundle

import (
	"os"
	"sync"

	"github.com/caos/boom/api/v1beta1"
	application "github.com/caos/boom/internal/bundle/application/mock"
	"github.com/caos/boom/internal/bundle/bundles"
	"github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/yaml"
	logcontext "github.com/caos/orbiter/logging/context"
	"github.com/caos/orbiter/logging/kubebuilder"
	"github.com/caos/orbiter/logging/stdlib"
	"github.com/stretchr/testify/assert"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"

	"testing"
)

const (
	baseDirectoryPath       = "../../tools"
	dashboardsDirectoryPath = "../../dashboards"
)

func NewBundle(templator name.Templator) *Bundle {
	logger := logcontext.Add(stdlib.New(os.Stdout))
	ctrl.SetLogger(kubebuilder.New(logger))

	bundleConf := &config.Config{
		Logger:            logger,
		CrdName:           "caos_test",
		BaseDirectoryPath: baseDirectoryPath,
		Templator:         templator,
	}

	b := New(bundleConf)
	return b
}

func TestBundle_EmptyApplicationList(t *testing.T) {
	b := NewBundle(yaml.GetName())
	eqApps := b.GetApplications()
	assert.Zero(t, len(eqApps))
}

func TestBundle_AddApplicationsByBundleName(t *testing.T) {
	b := NewBundle(yaml.GetName())

	//Add basisset
	err := b.AddApplicationsByBundleName(bundles.Caos)
	assert.NoError(t, err)
	apps := bundles.GetCaos()

	eqApps := b.GetApplications()
	assert.Equal(t, len(eqApps), len(apps))
	for eqApp := range eqApps {
		assert.Contains(t, apps, eqApp)
	}
}

func TestBundle_AddApplicationsByBundleName_nonexistent(t *testing.T) {
	b := NewBundle(yaml.GetName())
	var nonexistent name.Bundle
	nonexistent = "nonexistent"
	err := b.AddApplicationsByBundleName(nonexistent)
	assert.Error(t, err)
	eqApps := b.GetApplications()
	assert.Equal(t, 0, len(eqApps))
}
func TestBundle_AddApplication(t *testing.T) {
	b := NewBundle(yaml.GetName())

	spec := &v1beta1.ToolsetSpec{}
	app := application.NewTestYAMLApplication(t)
	app.AllowSetAppliedSpec(spec).SetChanged(spec, true).SetDeploy(spec, true).SetInitial(true).SetGetYaml("test")
	b.AddApplication(app.Application())

	apps := b.GetApplications()
	assert.Equal(t, 1, len(apps))
}

func TestBundle_AddApplication_AlreadyAdded(t *testing.T) {
	b := NewBundle(yaml.GetName())

	spec := &v1beta1.ToolsetSpec{}
	app := application.NewTestYAMLApplication(t)
	app.AllowSetAppliedSpec(spec).SetChanged(spec, true).SetDeploy(spec, true).SetInitial(true).SetGetYaml("test")
	err := b.AddApplication(app.Application()).GetStatus()
	assert.NoError(t, err)

	apps := b.GetApplications()
	assert.Equal(t, 1, len(apps))

	err2 := b.AddApplication(app.Application()).GetStatus()
	assert.Error(t, err2)

}

func TestBundle_ReconcileApplication(t *testing.T) {
	b := NewBundle(yaml.GetName())

	spec := &v1beta1.ToolsetSpec{}
	app := application.NewTestYAMLApplication(t)
	app.AllowSetAppliedSpec(spec).SetChanged(spec, true).SetDeploy(spec, true).SetInitial(true).SetGetYaml("test")
	b.AddApplication(app.Application())

	var wg sync.WaitGroup
	wg.Add(1)
	err := b.ReconcileApplication(app.Application().GetName(), spec, &wg).GetStatus()
	assert.NoError(t, err)
}

func TestBundle_ReconcileApplication_nonexistent(t *testing.T) {
	b := NewBundle(yaml.GetName())

	spec := &v1beta1.ToolsetSpec{}
	app := application.NewTestYAMLApplication(t)
	app.AllowSetAppliedSpec(spec).SetChanged(spec, true).SetDeploy(spec, true).SetInitial(true).SetGetYaml("test")

	var wg sync.WaitGroup
	wg.Add(1)
	err := b.ReconcileApplication(app.Application().GetName(), nil, &wg).GetStatus()
	assert.Error(t, err)
}

func TestBundle_Reconcile(t *testing.T) {
	b := NewBundle(yaml.GetName())

	spec := &v1beta1.ToolsetSpec{}
	app := application.NewTestYAMLApplication(t)
	app.AllowSetAppliedSpec(spec).SetChanged(spec, true).SetDeploy(spec, true).SetInitial(true).SetGetYaml("test")
	b.AddApplication(app.Application())

	err := b.Reconcile(spec).GetStatus()
	assert.NoError(t, err)
}

func TestBundle_Reconcile_NoApplications(t *testing.T) {
	b := NewBundle(yaml.GetName())

	spec := &v1beta1.ToolsetSpec{}

	err := b.Reconcile(spec).GetStatus()
	assert.NoError(t, err)
}
