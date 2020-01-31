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
	bundle := bundle.New(conf)

	c.status = bundle.AddApplicationsByBundleName(conf.BundleName)
	if c.status != nil {
		return
	}

	c.bundle = bundle
	c.status = nil
	return
}

func (c *Crd) GetBundle() *bundle.Bundle {
	return c.bundle
}

func (c *Crd) ReconcileWithFunc(getToolsetCRD func(instance runtime.Object) error) {
	if c.GetStatus() != nil {
		return
	}

	logFields := map[string]interface{}{
		"logID": "CRD-6e7csH4wkujsRYE",
	}

	if getToolsetCRD == nil {
		c.status = errors.New("ToolsetCRDFunc is nil")
		c.logger.WithFields(logFields).Error(c.status)
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
		"logID": "CRD-6e7csH4wkujsRYE",
	}

	if toolsetCRD == nil {
		c.status = errors.New("ToolsetCRD is nil")
		c.logger.WithFields(logFields).Error(c.status)
		return
	}
	logFields["CRD"] = toolsetCRD.Name

	c.logger.WithFields(logFields).Info("Reconcile applications")

	c.reconcileApplications(toolsetCRD.Name, toolsetCRD.Spec)
}

func (c *Crd) reconcileApplications(overlay string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) *Crd {
	if c.GetStatus() != nil {
		return c
	}

	logFields := map[string]interface{}{
		"logID": "CRD-sOwmSVvpAQfO1PS",
		"CRD":   overlay,
	}

	if c.bundle == nil {
		c.status = errors.New("No bundle for crd")
		c.logger.WithFields(logFields).Error(c.status)
		return c
	}

	c.status = c.bundle.Reconcile(toolsetCRDSpec).GetStatus()

	return c
}
