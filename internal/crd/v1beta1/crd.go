package v1beta1

import (
	"errors"

	"github.com/caos/orbiter/logging"
	"k8s.io/apimachinery/pkg/runtime"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle"
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/crd/v1beta1/config"

	"github.com/caos/boom/internal/name"
)

const (
	version name.Version = "v1beta1"
)

type Crd struct {
	bundle *bundle.Bundle
	logger logging.Logger
	status error
}

func (c *Crd) GetStatus() error {
	return c.status
}

func (c *Crd) CleanUp() {
	if c.GetStatus() != nil {
		return
	}

	c.status = c.bundle.CleanUp().GetStatus()
}

func GetVersion() name.Version {
	return version
}

func New(conf *config.Config) *Crd {
	crdLogger := conf.Logger.WithFields(map[string]interface{}{
		"version": "v1beta1",
	})

	return &Crd{
		logger: crdLogger,
		status: nil,
	}
}

func (c *Crd) SetBundle(conf *bundleconfig.Config) {
	if c.GetStatus() != nil {
		return
	}
	bundle := bundle.New(conf)

	c.status = bundle.AddApplicationsByBundleName(conf.BundleName)
	if c.status != nil {
		return
	}

	c.bundle = bundle
}

func (c *Crd) GetBundle() *bundle.Bundle {
	return c.bundle
}

func (c *Crd) ReconcileWithFunc(getToolsetCRD func(instance runtime.Object) error) {
	if c.GetStatus() != nil {
		return
	}

	if getToolsetCRD == nil {
		c.status = errors.New("ToolsetCRDFunc is nil")
		c.logger.Error(c.status)
		return
	}

	var toolsetCRD *toolsetsv1beta1.Toolset
	if err := getToolsetCRD(toolsetCRD); err != nil {
		c.status = err
		return
	}

	c.Reconcile(toolsetCRD)
}

func (c *Crd) Reconcile(toolsetCRD *toolsetsv1beta1.Toolset) {
	if c.GetStatus() != nil {
		return
	}
	logFields := map[string]interface{}{
		"CRD":    toolsetCRD.Name,
		"action": "reconciling",
	}
	logger := c.logger.WithFields(logFields)

	if toolsetCRD == nil {
		c.status = errors.New("ToolsetCRD is nil")
		logger.Error(c.status)
		return
	}

	if c.bundle == nil {
		c.status = errors.New("No bundle for crd")
		logger.Error(c.status)
		return
	}

	c.status = c.bundle.Reconcile(toolsetCRD.Spec).GetStatus()
}
