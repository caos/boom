package crd

import (
	"github.com/caos/orbiter/logging"
	"k8s.io/apimachinery/pkg/runtime"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/bundle"
)

type Crd struct {
	bundle *bundle.Bundle
	logger logging.Logger
	oldCrd *toolsetsv1beta1.Toolset
}

func (c *Crd) CleanUp() error {
	return c.bundle.CleanUp()
}

func GetVersion() string {
	return "v1beta1"
}

func New(logger logging.Logger, name string, toolsDirectoryPath, dashboardsDirectoryPath string) (*Crd, error) {

	bundle := bundle.New(logger, name, toolsDirectoryPath, dashboardsDirectoryPath)

	crd := &Crd{
		logger: logger,
		bundle: bundle,
	}

	return crd, nil
}

func (c *Crd) NewTemplate(new *toolsetsv1beta1.Toolset) bool {

	if c.oldCrd == nil || new.Spec != c.oldCrd.Spec {
		return true
	}
	return false
}

func (c *Crd) ReconcileWithFunc(getToolsetCRD func(instance runtime.Object) error) error {

	var toolsetCRD *toolsetsv1beta1.Toolset
	if err := getToolsetCRD(toolsetCRD); err != nil {
		return err
	}

	c.Reconcile(toolsetCRD)
	return nil
}

func (c *Crd) Reconcile(toolsetCRD *toolsetsv1beta1.Toolset) error {
	template := c.NewTemplate(toolsetCRD)
	if template {
		c.logger.WithFields(map[string]interface{}{
			"logID": "CRD-6e7csH4wkujsRYE",
			"CRD":   toolsetCRD.Name,
		}).Info("Reconcile applications")

		if err := c.reconcileApplications(toolsetCRD.Name, toolsetCRD.Spec); err != nil {
			return err
		}
	}

	c.oldCrd = toolsetCRD
	return nil
}

func (c *Crd) reconcileApplications(overlay string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) error {

	if err := c.bundle.Reconcile(toolsetCRDSpec); err != nil {
		return err
	}

	return nil
}
