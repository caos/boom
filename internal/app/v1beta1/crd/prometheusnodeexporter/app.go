package prometheusnodeexporter

import (
	"os"
	"path/filepath"
	"strings"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/kubectl"
	"github.com/caos/toolsop/internal/template"
	"github.com/caos/utils/logging"
)

var (
	applicationName      = "prometheus-node-exporter"
	resultsDirectoryName = "results"
	resultsFileName      = "results.yaml"
	defaultNamespace     = "monitoring"
)

type PrometheusNodeExporter struct {
	ApplicationDirectoryPath string
}

func New(toolsDirectoryPath string) *PrometheusNodeExporter {
	pne := &PrometheusNodeExporter{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
	}

	return pne
}

func (p *PrometheusNodeExporter) Reconcile(overlay string, helm *template.Helm, spec *toolsetsv1beta1.PrometheusNodeExporter) error {

	logging.Log("CRD-uXoxEPz8UVURI6g").Infof("Reconciling application %s", applicationName)
	resultsFileDirectory := filepath.Join(p.ApplicationDirectoryPath, resultsDirectoryName, overlay)
	_ = os.RemoveAll(resultsFileDirectory)
	_ = os.MkdirAll(resultsFileDirectory, os.ModePerm)
	resultFilePath := filepath.Join(resultsFileDirectory, resultsFileName)

	values := specToValues(helm.GetImageTags(applicationName), spec)

	writeValues := func(path string) error {
		if err := helper.StructToYaml(values, path); err != nil {
			logging.Log("CRD-zkqhLXoLJpLhUE9").Debugf("Failed to write values file overlay %s application %s", overlay, applicationName)
			return err
		}
		return nil
	}

	prefix := spec.Prefix
	if prefix == "" {
		prefix = overlay
	}
	namespace := spec.Namespace
	if namespace == "" {
		namespace = strings.Join([]string{overlay, defaultNamespace}, "-")
	}

	if err := helm.Template(applicationName, prefix, namespace, resultFilePath, writeValues); err != nil {
		return err
	}

	kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)
	if spec.Deploy {
		if err := kubectlCmd.Run(); err != nil {
			logging.Log("CRD-BcRGwbZs6siXam0").OnError(err).Debugf("Failed to apply file %s", resultFilePath)
			return err
		}
	}

	return nil
}

func specToValues(imageTags map[string]string, spec *toolsetsv1beta1.PrometheusNodeExporter) *Values {
	values := &Values{
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
		HostNetwork: true,
		Tolerations: []*Toleration{&Toleration{
			Effect:   "NoSchedule",
			Operator: "Exists",
		}},
	}

	if spec.Monitor != nil {
		values.Prometheus = &Prometheus{
			Monitor: &Monitor{
				Enabled:   spec.Monitor.Enabled,
				Namespace: spec.Monitor.Namespace,
			},
		}
	}

	return values
}