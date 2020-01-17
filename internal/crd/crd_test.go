package crd

import (
	"os"
	"testing"

	"github.com/caos/boom/api/v1beta1"
	application "github.com/caos/boom/internal/bundle/application/mock"
	"github.com/caos/boom/internal/bundle/bundles"
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/crd/config"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/yaml"
	logcontext "github.com/caos/orbiter/logging/context"
	"github.com/caos/orbiter/logging/kubebuilder"
	"github.com/caos/orbiter/logging/stdlib"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	fullToolset *v1beta1.Toolset = &v1beta1.Toolset{
		ObjectMeta: v1.ObjectMeta{
			Name:      "caos_test",
			Namespace: "caos-system",
		},
		Spec: &v1beta1.ToolsetSpec{
			KubeVersion: "1.17.0",
			Ambassador: &v1beta1.Ambassador{
				Deploy: true,
			},
			Argocd: &v1beta1.Argocd{
				Deploy: true,
			},
			KubeStateMetrics: &v1beta1.KubeStateMetrics{
				Deploy: true,
			},
			PrometheusOperator: &v1beta1.PrometheusOperator{
				Deploy: true,
			},
			PrometheusNodeExporter: &v1beta1.PrometheusNodeExporter{
				Deploy: true,
			},
			Grafana: &v1beta1.Grafana{
				Deploy: true,
			},
		},
	}
	changedToolset *v1beta1.Toolset = &v1beta1.Toolset{
		ObjectMeta: v1.ObjectMeta{
			Name:      "caos_test",
			Namespace: "caos-system",
		},
		Spec: &v1beta1.ToolsetSpec{
			KubeVersion: "1.17.0",
			Ambassador: &v1beta1.Ambassador{
				Deploy: false,
			},
			Argocd: &v1beta1.Argocd{
				Deploy: true,
			},
			KubeStateMetrics: &v1beta1.KubeStateMetrics{
				Deploy: true,
			},
			PrometheusOperator: &v1beta1.PrometheusOperator{
				Deploy: true,
			},
			PrometheusNodeExporter: &v1beta1.PrometheusNodeExporter{
				Deploy: true,
			},
			Grafana: &v1beta1.Grafana{
				Deploy: true,
			},
		},
	}
)

func newCrd() (Crd, error) {
	logger := logcontext.Add(stdlib.New(os.Stdout))
	ctrl.SetLogger(kubebuilder.New(logger))
	conf := &config.Config{
		Logger:  logger,
		Version: "v1beta1",
	}

	return New(conf)
}

func setBundle(c Crd, bundle name.Bundle) {
	logger := logcontext.Add(stdlib.New(os.Stdout))
	ctrl.SetLogger(kubebuilder.New(logger))
	bundleConfig := &bundleconfig.Config{
		Logger:                  logger,
		CrdName:                 "caos_test",
		BundleName:              bundle,
		BaseDirectoryPath:       "../../../tools",
		DashboardsDirectoryPath: "../../../dashboards",
		Templator:               yaml.GetName(),
	}

	c.SetBundle(bundleConfig)
}

func TestNew(t *testing.T) {
	crd, err := newCrd()
	assert.NoError(t, err)
	assert.NotNil(t, crd)
}

func TestNew_noexistendbundle(t *testing.T) {
	var nonexistent name.Bundle
	nonexistent = "nonexistent"
	crd, err := newCrd()
	assert.NoError(t, err)
	setBundle(crd, nonexistent)
	assert.Error(t, crd.GetStatus())
	assert.NotNil(t, crd)
}

func TestCrd_Reconcile_initial(t *testing.T) {
	crd, err := newCrd()
	setBundle(crd, bundles.Empty)
	bundle := crd.GetBundle()

	app := application.NewTestYAMLApplication(t)
	app.AllowSetAppliedSpec(fullToolset.Spec).SetChanged(fullToolset.Spec, true).SetDeploy(fullToolset.Spec, true).SetInitial(true).SetGetYaml("test")
	bundle.AddApplication(app.Application())
	assert.NoError(t, err)
	assert.NotNil(t, crd)

	// when crd is nil
	crd.Reconcile(fullToolset)
	err = crd.GetStatus()
	assert.NoError(t, err)
}

func TestCrd_Reconcile_changed(t *testing.T) {
	crd, err := newCrd()
	setBundle(crd, bundles.Empty)
	bundle := crd.GetBundle()

	app := application.NewTestYAMLApplication(t)
	app.AllowSetAppliedSpec(fullToolset.Spec).SetChanged(fullToolset.Spec, true).SetDeploy(fullToolset.Spec, true).SetInitial(true).SetGetYaml("test")
	bundle.AddApplication(app.Application())
	assert.NoError(t, err)
	assert.NotNil(t, crd)

	// when crd is nil
	crd.Reconcile(fullToolset)
	err = crd.GetStatus()
	assert.NoError(t, err)

	//changed crd
	app.AllowSetAppliedSpec(changedToolset.Spec).SetChanged(changedToolset.Spec, true).SetDeploy(changedToolset.Spec, true).SetInitial(false).SetGetYaml("test2")
	crd.Reconcile(changedToolset)
	err = crd.GetStatus()
	assert.NoError(t, err)
}

func TestCrd_Reconcile_changedDelete(t *testing.T) {
	crd, err := newCrd()
	setBundle(crd, bundles.Empty)
	bundle := crd.GetBundle()

	app := application.NewTestYAMLApplication(t)
	app.AllowSetAppliedSpec(fullToolset.Spec).SetChanged(fullToolset.Spec, true).SetDeploy(fullToolset.Spec, true).SetInitial(true).SetGetYaml("test")
	bundle.AddApplication(app.Application())
	assert.NoError(t, err)
	assert.NotNil(t, crd)

	// when crd is nil
	crd.Reconcile(fullToolset)
	err = crd.GetStatus()
	assert.NoError(t, err)

	//changed crd
	app.AllowSetAppliedSpec(changedToolset.Spec).SetChanged(changedToolset.Spec, true).SetDeploy(changedToolset.Spec, false).SetInitial(false).SetGetYaml("test2")
	crd.Reconcile(changedToolset)
	err = crd.GetStatus()
	assert.NoError(t, err)
}

func TestCrd_Reconcile_initialNotDeployed(t *testing.T) {
	crd, err := newCrd()
	setBundle(crd, bundles.Empty)
	bundle := crd.GetBundle()

	app := application.NewTestYAMLApplication(t)
	app.AllowSetAppliedSpec(fullToolset.Spec).SetChanged(fullToolset.Spec, true).SetDeploy(fullToolset.Spec, false).SetInitial(true).SetGetYaml("test")
	bundle.AddApplication(app.Application())
	assert.NoError(t, err)
	assert.NotNil(t, crd)

	// when crd is nil
	crd.Reconcile(fullToolset)
	err = crd.GetStatus()
	assert.NoError(t, err)

	//changed crd
	app.AllowSetAppliedSpec(changedToolset.Spec).SetChanged(changedToolset.Spec, true).SetDeploy(changedToolset.Spec, false).SetInitial(false).SetGetYaml("test2")
	crd.Reconcile(changedToolset)
	err = crd.GetStatus()
	assert.NoError(t, err)
}

func TestCrd_ReconcileWithFunc(t *testing.T) {
	assert.True(t, true)

	//TODO: correct function to read crd

	// crd, err := newCrd()
	// setBundle(crd, bundles.Empty)
	// bundle := crd.GetBundle()

	// app := application.NewTestYAMLApplication(t)
	// app.AllowSetAppliedSpec(fullToolset.Spec).SetChanged(fullToolset.Spec, true).SetDeploy(fullToolset.Spec, false).SetInitial(true).SetGetYaml("test")
	// bundle.AddApplication(app.Application())
	// assert.NoError(t, err)
	// assert.NotNil(t, crd)

	// getToolsetFunc := func(obj runtime.Object) error {
	// 	obj = fullToolset
	// 	return nil
	// }

	// // when crd is nil
	// crd.ReconcileWithFunc(getToolsetFunc)
	// err = crd.GetStatus()
	// assert.NoError(t, err)
}
