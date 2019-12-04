package app

import (
	"github.com/caos/orbiter/logging"
	"github.com/caos/boom/internal/app/v1beta1/crd/ambassador"
	"github.com/caos/boom/internal/app/v1beta1/crd/certmanager"
	"github.com/caos/boom/internal/app/v1beta1/crd/grafana"
	"github.com/caos/boom/internal/app/v1beta1/crd/loggingoperator"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheus"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheusnodeexporter"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheusoperator"
	"github.com/caos/boom/internal/template"
	"github.com/caos/boom/internal/toolset"
	"k8s.io/apimachinery/pkg/runtime"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

type Crd struct {
	helm   *template.Helm
	oldCrd *toolsetsv1beta1.Toolset
	logger logging.Logger
}

func (c *Crd) CleanUp() error {
	return c.helm.CleanUp()
}

func GetVersion() string {
	return "v1beta1"
}

func NewWithFunc(logger logging.Logger, getToolset func(obj runtime.Object) error, toolsDirectoryPath string, toolsets *toolset.Toolsets) (*Crd, error) {
	var toolsetCRD *toolsetsv1beta1.Toolset
	if err := getToolset(toolsetCRD); err != nil {
		return nil, err
	}
	return New(logger, toolsetCRD, toolsDirectoryPath, toolsets)
}

func New(logger logging.Logger, toolsetCRD *toolsetsv1beta1.Toolset, toolsDirectoryPath string, toolsets *toolset.Toolsets) (*Crd, error) {
	crd := &Crd{logger: logger}

	if err := crd.GenerateTemplateComponents(toolsDirectoryPath, toolsets, toolsetCRD); err != nil {
		return nil, err
	}

	if err := crd.ReconcileApplications(toolsetCRD.Name, toolsDirectoryPath, toolsetCRD.Spec); err != nil {
		return nil, err
	}
	crd.oldCrd = toolsetCRD

	return crd, nil
}

func (c *Crd) ReconcileWithFunc(getToolset func(obj runtime.Object) error, toolsDirectoryPath string, toolsets *toolset.Toolsets) error {
	var new *toolsetsv1beta1.Toolset
	if err := getToolset(new); err != nil {
		return err
	}

	return c.Reconcile(new, toolsDirectoryPath, toolsets)
}

func (c *Crd) Reconcile(new *toolsetsv1beta1.Toolset, toolsDirectoryPath string, toolsets *toolset.Toolsets) error {
	fetcherGen := c.NewFetcherGeneration(new)
	if fetcherGen {
		c.logger.WithFields(map[string]interface{}{
			"logID": "CRD-6e7csH4wkujsRYE",
		}).Info("Generate template components")
		if err := c.GenerateTemplateComponents(toolsDirectoryPath, toolsets, new); err != nil {
			return err
		}
	}

	template := c.NewTemplate(new)
	if template {
		c.logger.WithFields(map[string]interface{}{
			"logID": "CRD-6e7csH4wkujsRYE",
			"CRD":   new.Name,
		}).Info("Reconcile applications")
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

	helm, err := template.NewHelm(c.logger, toolsDirectoryPath, toolsets, toolsetCRD.Spec.Name, toolsetCRD.Spec.Version, toolsetCRD.Name)
	if err != nil {
		return err
	}
	c.helm = helm
	return nil
}

func (c *Crd) NewFetcherGeneration(new *toolsetsv1beta1.Toolset) bool {
	if new.Spec.Name != c.oldCrd.Spec.Name || new.Spec.Version != c.oldCrd.Spec.Version {
		return true
	}
	return false
}

func (c *Crd) NewTemplate(new *toolsetsv1beta1.Toolset) bool {
	fetcher := c.NewFetcherGeneration(new)

	if fetcher ||
		new.Spec.LoggingOperator != c.oldCrd.Spec.LoggingOperator ||
		new.Spec.PrometheusOperator != c.oldCrd.Spec.PrometheusOperator ||
		new.Spec.PrometheusNodeExporter != c.oldCrd.Spec.PrometheusNodeExporter ||
		new.Spec.Grafana != c.oldCrd.Spec.Grafana {
		return true
	}
	return false
}

func (c *Crd) ReconcileApplications(overlay, toolsDirectoryPath string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) error {
	lo := loggingoperator.New(c.logger, toolsDirectoryPath)
	if err := lo.Reconcile(overlay, c.helm, toolsetCRDSpec.LoggingOperator); err != nil {
		return err
	}

	po := prometheusoperator.New(c.logger, toolsDirectoryPath)
	if err := po.Reconcile(overlay, c.helm, toolsetCRDSpec.PrometheusOperator); err != nil {
		return err
	}

	pne := prometheusnodeexporter.New(c.logger, toolsDirectoryPath)
	if err := pne.Reconcile(overlay, c.helm, toolsetCRDSpec.PrometheusNodeExporter); err != nil {
		return err
	}

	g := grafana.New(c.logger, toolsDirectoryPath)
	if err := g.Reconcile(overlay, c.helm, toolsetCRDSpec.Grafana); err != nil {
		return err
	}

	p := prometheus.New(c.logger, toolsDirectoryPath)
	if err := p.Reconcile(overlay, c.helm, toolsetCRDSpec.Prometheus); err != nil {
		return err
	}

	cert := certmanager.New(c.logger, toolsDirectoryPath)
	if err := cert.Reconcile(overlay, c.helm, toolsetCRDSpec.CertManager); err != nil {
		return err
	}

	ambassador := ambassador.New(c.logger, toolsDirectoryPath)
	if err := ambassador.Reconcile(overlay, c.helm, toolsetCRDSpec.Ambassador); err != nil {
		return err
	}

	return nil
}
