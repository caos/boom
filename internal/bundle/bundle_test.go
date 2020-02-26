package bundle

import (
	"sync"

	"github.com/caos/boom/api/v1beta1"
	application "github.com/caos/boom/internal/bundle/application/mock"
	"github.com/caos/boom/internal/bundle/bundles"
	"github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/yaml"
	"github.com/caos/orbiter/mntr"
	"github.com/stretchr/testify/assert"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"testing"
)

const (
	baseDirectoryPath       = "../../tools"
	dashboardsDirectoryPath = "../../dashboards"
)

func newMonitor() mntr.Monitor {
	monitor := mntr.Monitor{
		OnInfo:   mntr.LogMessage,
		OnChange: mntr.LogMessage,
		OnError:  mntr.LogError,
	}

	return monitor
}

func NewBundle(templator name.Templator) *Bundle {
	monitor := newMonitor()

	bundleConf := &config.Config{
		Monitor:           monitor,
		CrdName:           "caos_test",
		BaseDirectoryPath: baseDirectoryPath,
		Templator:         templator,
	}

	b := New(bundleConf)
	return b
}

func init() {
	Testmode = true
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
	app.SetDeploy(spec, true).SetGetYaml(spec, "test")
	b.AddApplication(app.Application())

	apps := b.GetApplications()
	assert.Equal(t, 1, len(apps))
}

func TestBundle_AddApplication_AlreadyAdded(t *testing.T) {
	b := NewBundle(yaml.GetName())

	spec := &v1beta1.ToolsetSpec{}
	app := application.NewTestYAMLApplication(t)
	app.SetDeploy(spec, true).SetGetYaml(spec, "test")
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
	app.SetDeploy(spec, true).SetGetYaml(spec, "test")

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
	app.SetDeploy(spec, true).SetGetYaml(spec, "test")

	var wg sync.WaitGroup
	wg.Add(1)
	err := b.ReconcileApplication(app.Application().GetName(), nil, &wg).GetStatus()
	assert.Error(t, err)
}

func TestBundle_Reconcile(t *testing.T) {
	b := NewBundle(yaml.GetName())

	spec := &v1beta1.ToolsetSpec{}
	app := application.NewTestYAMLApplication(t)
	app.SetDeploy(spec, true).SetGetYaml(spec, "test")
	b.AddApplication(app.Application())

	b.Reconcile(spec)
	err := b.GetStatus()
	assert.NoError(t, err)
}

func TestBundle_Reconcile_NoApplications(t *testing.T) {
	b := NewBundle(yaml.GetName())

	spec := &v1beta1.ToolsetSpec{}
	b.Reconcile(spec)
	err := b.GetStatus()
	assert.NoError(t, err)
}
