package crd

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/caos/boom/internal/bundle"
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/crd/config"
	v1beta1config "github.com/caos/boom/internal/crd/v1beta1/config"

	"github.com/caos/boom/internal/crd/v1beta1"

	"github.com/pkg/errors"
)

type Crd interface {
	SetBundle(*bundleconfig.Config)
	GetBundle() *bundle.Bundle
	ReconcileWithFunc(getToolsetCRD func(instance runtime.Object) error)
	Reconcile(toolsetCRD *toolsetsv1beta1.Toolset)
	CleanUp()
	GetStatus() error
}

func New(conf *config.Config) (Crd, error) {
	crdLogger := conf.Logger.WithFields(map[string]interface{}{
		"version": conf.Version,
	})

	crdLogger.Info("New CRD")

	if conf.Version != "v1beta1" {
		return nil, errors.Errorf("Unknown CRD version %s", conf.Version)
	}

	crdConf := &v1beta1config.Config{
		Logger: conf.Logger,
	}

	crd := v1beta1.New(crdConf)

	return crd, nil
}
