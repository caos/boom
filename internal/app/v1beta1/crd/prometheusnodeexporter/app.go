package prometheusnodeexporter

import (
	"path/filepath"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/v1beta1/crd/defaults"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/template"
)

var (
	applicationName = "prometheus-node-exporter"
)

type PrometheusNodeExporter struct {
	ApplicationDirectoryPath string
	logger                   logging.Logger
	spec                     *toolsetsv1beta1.PrometheusNodeExporter
}

func New(logger logging.Logger, toolsDirectoryPath string) *PrometheusNodeExporter {
	pne := &PrometheusNodeExporter{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
		logger:                   logger,
	}

	return pne
}

func (p *PrometheusNodeExporter) Reconcile(overlay, specNamespace string, helm *template.Helm, spec *toolsetsv1beta1.PrometheusNodeExporter) error {

	logFields := map[string]interface{}{
		"application": applicationName,
		"logID":       "CRD-8Z2ueDkBmkBONgc",
	}
	p.logger.WithFields(logFields).Info("Reconciling")

	resultFilePath := defaults.GetResultFilePath(overlay, p.ApplicationDirectoryPath, applicationName)
	prefix := defaults.GetPrefix(overlay, applicationName, spec.Prefix)
	namespace := defaults.GetNamespace(overlay, applicationName, specNamespace, spec.Namespace)

	values := specToValues(helm.GetImageTags(applicationName), spec)

	writeValues := func(path string) error {
		if err := errors.Wrapf(helper.StructToYaml(values, path), "Failed to write values file overlay %s application %s", overlay, applicationName); err != nil {
			return err
		}
		return nil
	}

	if err := helm.PrepareTemplate(applicationName, prefix, namespace, writeValues); err != nil {
		return err
	}

	if spec.Deploy {
		if err := defaults.PrepareForResultOutput(defaults.GetResultFileDirectory(overlay, p.ApplicationDirectoryPath, applicationName)); err != nil {
			return err
		}

		if err := helm.Template(applicationName, resultFilePath); err != nil {
			return err
		}

		if err := helper.DeleteKindFromYaml(resultFilePath, "Namespace"); err != nil {
			return err
		}

		kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)
		if err := errors.Wrapf(helper.Run(p.logger, kubectlCmd.Build()), "Failed to apply with file %s", resultFilePath); err != nil {
			return err
		}

		p.spec = spec
	} else if !spec.Deploy && p.spec != nil {
		kubectlCmd := kubectl.New("delete").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)
		if err := errors.Wrapf(helper.Run(p.logger, kubectlCmd.Build()), "Failed to delete with file %s", resultFilePath); err != nil {
			return err
		}

		p.spec = nil
	}

	return nil
}

func specToValues(imageTags map[string]string, spec *toolsetsv1beta1.PrometheusNodeExporter) *Values {
	values := &Values{
		FullnameOverride: "node-exporter",
		Image: &Image{
			Repository: "quay.io/prometheus/node-exporter",
			Tag:        imageTags["quay.io/prometheus/node-exporter"],
			PullPolicy: "IfNotPresent",
		},
		Service: &Service{
			Type:        "ClusterIP",
			Port:        9100,
			TargetPort:  9100,
			NodePort:    "",
			Annotations: map[string]string{"prometheus.io/scrape": "true"},
		},
		Prometheus: &Prometheus{
			Monitor: &Monitor{
				Enabled:          false,
				AdditionalLabels: map[string]string{},
				Namespace:        "",
				ScrapeTimeout:    "10s",
			},
		},
		ServiceAccount: &ServiceAccount{
			Create: true,
		},
		SecurityContext: &SecurityContext{
			RunAsNonRoot: true,
			RunAsUser:    65534,
		},
		Rbac: &Rbac{
			Create:     true,
			PspEnabled: true,
		},
		HostNetwork: false,
		Tolerations: []*Toleration{&Toleration{
			Effect:   "NoSchedule",
			Operator: "Exists",
		}},
	}

	return values
}
