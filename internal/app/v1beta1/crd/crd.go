package app

import (
	"strings"

	"github.com/caos/boom/internal/app/v1beta1/crd/ambassador"
	"github.com/caos/boom/internal/app/v1beta1/crd/certmanager"
	"github.com/caos/boom/internal/app/v1beta1/crd/grafana"
	"github.com/caos/boom/internal/app/v1beta1/crd/kubestatemetrics"
	"github.com/caos/boom/internal/app/v1beta1/crd/loggingoperator"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheus"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheus/servicemonitor"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheusnodeexporter"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheusoperator"
	"github.com/caos/boom/internal/template"
	"github.com/caos/boom/internal/toolset"
	"github.com/caos/orbiter/logging"
	"k8s.io/apimachinery/pkg/runtime"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

type Crd struct {
	helm         *template.Helm
	oldCrd       *toolsetsv1beta1.Toolset
	applications *Applications
	logger       logging.Logger
}

type Applications struct {
	LoggingOperator        *loggingoperator.LoggingOperator
	Ambassador             *ambassador.Ambassador
	Prometheus             *prometheus.Prometheus
	PrometheusOperator     *prometheusoperator.PrometheusOperator
	PrometheusNodeExporter *prometheusnodeexporter.PrometheusNodeExporter
	Grafana                *grafana.Grafana
	CertManager            *certmanager.CertManager
	KubeStateMetrics       *kubestatemetrics.KubeStateMetrics
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

	apps, err := crd.NewApplications(toolsDirectoryPath)
	if err != nil {
		return nil, err
	}
	crd.applications = apps

	helm, err := template.NewHelm(crd.logger, toolsDirectoryPath, toolsets, toolsetCRD.Spec.Name, toolsetCRD.Name)
	if err != nil {
		return nil, err
	}
	crd.helm = helm

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

func (c *Crd) NewTemplate(new *toolsetsv1beta1.Toolset) bool {

	if new.Spec != c.oldCrd.Spec {
		return true
	}
	return false
}

func (c *Crd) ReconcileApplications(overlay, toolsDirectoryPath string, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) error {

	if err := c.applications.LoggingOperator.Reconcile(overlay, toolsetCRDSpec.Namespace, c.helm, toolsetCRDSpec.LoggingOperator); err != nil {
		return err
	}

	if err := c.applications.PrometheusOperator.Reconcile(overlay, toolsetCRDSpec.Namespace, c.helm, toolsetCRDSpec.PrometheusOperator); err != nil {
		return err
	}

	if err := c.applications.PrometheusNodeExporter.Reconcile(overlay, toolsetCRDSpec.Namespace, c.helm, toolsetCRDSpec.PrometheusNodeExporter); err != nil {
		return err
	}

	if err := c.applications.CertManager.Reconcile(overlay, toolsetCRDSpec.Namespace, c.helm, toolsetCRDSpec.CertManager); err != nil {
		return err
	}

	if err := c.applications.Ambassador.Reconcile(overlay, toolsetCRDSpec.Namespace, c.helm, toolsetCRDSpec.Ambassador); err != nil {
		return err
	}

	if err := c.applications.KubeStateMetrics.Reconcile(overlay, toolsetCRDSpec.Namespace, c.helm, toolsetCRDSpec.KubeStateMetrics); err != nil {
		return err
	}

	conf, datasource, err := c.ScrapeMetricsCrdsConfig(toolsetCRDSpec)
	if err != nil {
		return err
	}
	if conf != nil {
		if err := c.applications.Prometheus.Reconcile(overlay, toolsetCRDSpec.Namespace, c.helm, conf); err != nil {
			return err
		}
	}

	toolsetCRDSpec.Grafana.Datasources = append(toolsetCRDSpec.Grafana.Datasources, &toolsetsv1beta1.Datasource{
		Name:   "prometheus",
		Type:   "prometheus",
		Url:    datasource,
		Access: "proxy",
	})

	if err := c.applications.Grafana.Reconcile(overlay, toolsetCRDSpec.Namespace, c.helm, toolsetCRDSpec.Grafana); err != nil {
		return err
	}

	return nil
}

func (c *Crd) NewApplications(toolsDirectoryPath string) (*Applications, error) {
	applications := &Applications{
		LoggingOperator:        loggingoperator.New(c.logger, toolsDirectoryPath),
		PrometheusOperator:     prometheusoperator.New(c.logger, toolsDirectoryPath),
		PrometheusNodeExporter: prometheusnodeexporter.New(c.logger, toolsDirectoryPath),
		CertManager:            certmanager.New(c.logger, toolsDirectoryPath),
		Ambassador:             ambassador.New(c.logger, toolsDirectoryPath),
		Prometheus:             prometheus.New(c.logger, toolsDirectoryPath),
		Grafana:                grafana.New(c.logger, toolsDirectoryPath),
		KubeStateMetrics:       kubestatemetrics.New(c.logger, toolsDirectoryPath),
	}

	return applications, nil
}

func (c *Crd) ScrapeMetricsCrdsConfig(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) (*prometheus.Config, string, error) {

	monitorlabels := make(map[string]string, 0)
	monitorlabels["app.kubernetes.io/managed-by"] = "boom.caos.ch"

	servicemonitors := make([]*servicemonitor.Config, 0)

	if toolsetCRDSpec.Ambassador.Deploy {
		endpoint := &servicemonitor.ConfigEndpoint{
			Port: "ambassador-admin",
			Path: "/metrics",
		}
		labels := map[string]string{"service": "ambassador-admin"}

		smconfig := &servicemonitor.Config{
			Name:                  "ambassador-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	if toolsetCRDSpec.CertManager.Deploy {
		endpoint := &servicemonitor.ConfigEndpoint{
			TargetPort: "9402",
			Path:       "/metrics",
		}
		labels := map[string]string{"app": "cert-manager"}

		smconfig := &servicemonitor.Config{
			Name:                  "cert-manager-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	if toolsetCRDSpec.PrometheusOperator.Deploy {
		endpoint := &servicemonitor.ConfigEndpoint{
			Port: "http",
			Path: "/metrics",
		}
		labels := map[string]string{"app": "prometheus-operator-operator"}

		smconfig := &servicemonitor.Config{
			Name:                  "prometheus-operator-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	if toolsetCRDSpec.PrometheusNodeExporter.Deploy {
		endpoint := &servicemonitor.ConfigEndpoint{
			Port: "metrics",
			Path: "/metrics",
		}
		labels := map[string]string{"app": "prometheus-node-exporter"}

		smconfig := &servicemonitor.Config{
			Name:                  "prometheus-node-exporter-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	if toolsetCRDSpec.KubeStateMetrics.Deploy {
		endpoint := &servicemonitor.ConfigEndpoint{
			Port: "http",
			Path: "/metrics",
		}

		labels := map[string]string{
			"app.kubernetes.io/name": "kube-state-metrics",
		}

		smconfig := &servicemonitor.Config{
			Name:                  "kube-state-metrics-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)
	}

	if len(servicemonitors) > 0 {
		endpoint := &servicemonitor.ConfigEndpoint{
			Port: "web",
			Path: "/metrics",
		}
		labels := map[string]string{"app": "prometheus-operator-prometheus", "release": "caos"}

		smconfig := &servicemonitor.Config{
			Name:                  "prometheus-servicemonitor",
			Endpoints:             []*servicemonitor.ConfigEndpoint{endpoint},
			MonitorMatchingLabels: monitorlabels,
			ServiceMatchingLabels: labels,
		}
		servicemonitors = append(servicemonitors, smconfig)

		prom := &prometheus.Config{
			Prefix:          "caos",
			Namespace:       "caos-system",
			MonitorLabels:   monitorlabels,
			ServiceMonitors: servicemonitors,
		}
		datasource := strings.Join([]string{"http://", prom.Prefix, "-prometheus-operator-prometheus.", prom.Namespace, ":9090"}, "")

		return prom, datasource, nil
	}

	return nil, "", nil
}
