package app

import (
	"github.com/caos/toolsop/internal/app/crd/grafana"
	"github.com/caos/toolsop/internal/app/crd/loggingoperator"
	"github.com/caos/toolsop/internal/app/crd/prometheus"
	"github.com/caos/toolsop/internal/app/crd/prometheusnodeexporter"
	"github.com/caos/toolsop/internal/app/crd/prometheusoperator"
	"github.com/caos/toolsop/internal/template"
	"github.com/caos/toolsop/internal/toolset"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
)

type Crd struct {
	helm   *template.Helm
	oldCrd *toolsetsv1beta1.Toolset
}

func New(toolsetCRD *toolsetsv1beta1.Toolset, toolsDirectoryPath string, toolsets *toolset.Toolsets) (*Crd, error) {
	crd := &Crd{}

	if err := crd.GenerateTemplateComponents(toolsDirectoryPath, toolsets, toolsetCRD); err != nil {
		return nil, err
	}

	if err := crd.ReconcileApplications(toolsetCRD.Name, toolsDirectoryPath, toolsetCRD.Spec); err != nil {
		return nil, err
	}
	crd.oldCrd = toolsetCRD

	return crd, nil
}

func (c *Crd) Reconcile(new *toolsetsv1beta1.Toolset, toolsDirectoryPath string, toolsets *toolset.Toolsets) error {

	fetcherGen, err := c.NewFetcherGeneration(new)
	if err != nil {
		return nil
	}
	if fetcherGen {
		if err := c.GenerateTemplateComponents(toolsDirectoryPath, toolsets, new); err != nil {
			return err
		}
	}

	template, err := c.NewTemplate(new)
	if err != nil {
		return nil
	}
	if template {
		if err := c.ReconcileApplications(new.Name, toolsDirectoryPath, new.Spec); err != nil {
			return err
		}
	}

	c.oldCrd = new
	return nil
}

func (c *Crd) GenerateTemplateComponents(toolsDirectoryPath string, toolsets *toolset.Toolsets, toolsetCRD *toolsetsv1beta1.Toolset) error {
	if c.helm != nil {
		c.helm.CleanUp()
	}

	helm, err := template.NewHelm(toolsDirectoryPath, toolsets, toolsetCRD.Spec.Name, toolsetCRD.Spec.Version, toolsetCRD.Name)
	if err != nil {
		return err
	}
	c.helm = helm
	return nil
}

func (c *Crd) NewFetcherGeneration(new *toolsetsv1beta1.Toolset) (bool, error) {
	if new.Spec.Name != c.oldCrd.Spec.Name || new.Spec.Version != c.oldCrd.Spec.Version {
		return true, nil
	}
	return false, nil
}

func (c *Crd) NewTemplate(new *toolsetsv1beta1.Toolset) (bool, error) {
	fetcher, err := c.NewFetcherGeneration(new)
	if err != nil {
		return false, err
	}

	if fetcher ||
		new.Spec.LoggingOperator != c.oldCrd.Spec.LoggingOperator ||
		new.Spec.PrometheusOperator != c.oldCrd.Spec.PrometheusOperator ||
		new.Spec.PrometheusNodeExporter != c.oldCrd.Spec.PrometheusNodeExporter ||
		new.Spec.Grafana != c.oldCrd.Spec.Grafana {
		return true, nil
	}
	return false, nil
}

func (c *Crd) ReconcileApplications(overlay, toolsDirectoryPath string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) error {
	lo := loggingoperator.New(toolsDirectoryPath)
	if err := lo.Reconcile(overlay, c.helm, toolsetCRDSpec.LoggingOperator); err != nil {
		return err
	}

	po := prometheusoperator.New(toolsDirectoryPath)
	if err := po.Reconcile(overlay, c.helm, toolsetCRDSpec.PrometheusOperator); err != nil {
		return err
	}

	pne := prometheusnodeexporter.New(toolsDirectoryPath)
	if err := pne.Reconcile(overlay, c.helm, toolsetCRDSpec.PrometheusNodeExporter); err != nil {
		return err
	}

	g := grafana.New(toolsDirectoryPath)
	if err := g.Reconcile(overlay, c.helm, toolsetCRDSpec.Grafana); err != nil {
		return err
	}

	p := prometheus.New(toolsDirectoryPath)
	if err := p.Reconcile(overlay, c.helm, toolsetCRDSpec.Prometheus); err != nil {
		return err
	}

	return nil
}
