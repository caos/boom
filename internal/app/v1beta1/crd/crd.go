package app

import (
	"github.com/caos/toolsop/internal/app/v1beta1/crd/ambassador"
	"github.com/caos/toolsop/internal/app/v1beta1/crd/certmanager"
	"github.com/caos/toolsop/internal/app/v1beta1/crd/grafana"
	"github.com/caos/toolsop/internal/app/v1beta1/crd/loggingoperator"
	"github.com/caos/toolsop/internal/app/v1beta1/crd/prometheus"
	"github.com/caos/toolsop/internal/app/v1beta1/crd/prometheusnodeexporter"
	"github.com/caos/toolsop/internal/app/v1beta1/crd/prometheusoperator"
	"github.com/caos/toolsop/internal/template"
	"github.com/caos/toolsop/internal/toolset"
	"github.com/caos/utils/logging"
	"k8s.io/apimachinery/pkg/runtime"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
)

type Crd struct {
	helm   *template.Helm
	oldCrd *toolsetsv1beta1.Toolset
}

func (c *Crd) CleanUp() error {
	return c.helm.CleanUp()
}

func GetVersion() string {
	return "v1beta1"
}

func NewWithFunc(getToolset func(obj runtime.Object) error, toolsDirectoryPath string, toolsets *toolset.Toolsets) (*Crd, error) {
	var toolsetCRD *toolsetsv1beta1.Toolset
	if err := getToolset(toolsetCRD); err != nil {
		return nil, err
	}
	return New(toolsetCRD, toolsDirectoryPath, toolsets)
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
		logging.Log("CRD-6e7csH4wkujsRYE").Infof("Generate template components")
		if err := c.GenerateTemplateComponents(toolsDirectoryPath, toolsets, new); err != nil {
			return err
		}
	}

	template := c.NewTemplate(new)
	if template {
		logging.Log("CRD-6e7csH4wkujsRYE").Infof("Reconcile applications for CRD %s", new.Name)
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

	cert := certmanager.New(toolsDirectoryPath)
	if err := cert.Reconcile(overlay, c.helm, toolsetCRDSpec.CertManager); err != nil {
		return err
	}

	ambassador := ambassador.New(toolsDirectoryPath)
	if err := ambassador.Reconcile(overlay, c.helm, toolsetCRDSpec.Ambassador); err != nil {
		return err
	}

	return nil
}
