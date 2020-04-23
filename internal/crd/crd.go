package crd

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle"
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/crd/config"
	v1beta1config "github.com/caos/boom/internal/crd/v1beta1/config"

	"github.com/caos/boom/internal/crd/v1beta1"

	"github.com/pkg/errors"
)

type Crd interface {
	SetBundle(*bundleconfig.Config)
	GetBundle() *bundle.Bundle
	//ReconcileWithFunc([]*clientgo.Resource, func(instance runtime.Object) error)
	Reconcile([]*clientgo.Resource, *toolsetsv1beta1.Toolset)
	CleanUp()
	GetStatus() error
	SetBackStatus()
}

func New(conf *config.Config) (Crd, error) {
	crdMonitor := conf.Monitor.WithFields(map[string]interface{}{
		"version": conf.Version,
	})

	crdMonitor.Info("New CRD")

	if conf.Version != "v1beta1" {
		return nil, errors.Errorf("Unknown CRD version %s", conf.Version)
	}

	crdConf := &v1beta1config.Config{
		Monitor: crdMonitor,
	}

	crd := v1beta1.New(crdConf)

	return crd, nil
}
