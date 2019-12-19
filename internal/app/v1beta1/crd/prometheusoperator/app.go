package prometheusoperator

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
	applicationName = "prometheus-operator"
)

type PrometheusOperator struct {
	ApplicationDirectoryPath string
	logger                   logging.Logger
	spec                     *toolsetsv1beta1.PrometheusOperator
}

func New(logger logging.Logger, toolsDirectoryPath string) *PrometheusOperator {
	lo := &PrometheusOperator{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
		logger:                   logger,
	}

	return lo
}

func (p *PrometheusOperator) Reconcile(overlay, specNamespace string, helm *template.Helm, spec *toolsetsv1beta1.PrometheusOperator) error {

	logFields := map[string]interface{}{
		"application": applicationName,
		"logID":       "CRD-2JlElqA8Zqu7wcw",
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

		if err := helper.DeletePartOfYaml(resultFilePath, "kind: Namespace"); err != nil {
			return err
		}

		kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath)
		if err := errors.Wrapf(helper.Run(p.logger, kubectlCmd.Build()), "Failed to apply with file %s", resultFilePath); err != nil {
			return err
		}

		p.spec = spec
	} else if !spec.Deploy && p.spec != nil {
		kubectlCmd := kubectl.New("delete").AddParameter("-f", resultFilePath)
		if err := errors.Wrapf(helper.Run(p.logger, kubectlCmd.Build()), "Failed to delete with file %s", resultFilePath); err != nil {
			return err
		}

		p.spec = nil
	}

	return nil
}

func specToValues(imageTags map[string]string, spec *toolsetsv1beta1.PrometheusOperator) *Values {

	values := &Values{
		DefaultRules: &DefaultRules{
			Create: true,
			Rules: &Rules{
				Alertmanager:                true,
				Etcd:                        true,
				General:                     true,
				K8S:                         true,
				KubeApiserver:               true,
				KubePrometheusNodeAlerting:  true,
				KubePrometheusNodeRecording: true,
				KubernetesAbsent:            true,
				KubernetesApps:              true,
				KubernetesResources:         true,
				KubernetesStorage:           true,
				KubernetesSystem:            true,
				KubeScheduler:               true,
				Network:                     true,
				Node:                        true,
				Prometheus:                  true,
				PrometheusOperator:          true,
				Time:                        true,
			},
		},
		Global: &Global{
			Rbac: &Rbac{
				Create:     true,
				PspEnabled: true,
			},
		},
		Alertmanager: &DisabledTool{
			Enabled: false,
		},
		Grafana: &DisabledTool{
			Enabled: false,
		},
		KubeAPIServer: &DisabledTool{
			Enabled: false,
		},
		Kubelet: &DisabledTool{
			Enabled: false,
		},
		KubeControllerManager: &DisabledTool{
			Enabled: false,
		},
		CoreDNS: &DisabledTool{
			Enabled: false,
		},
		KubeDNS: &DisabledTool{
			Enabled: false,
		},
		KubeEtcd: &DisabledTool{
			Enabled: false,
		},
		KubeScheduler: &DisabledTool{
			Enabled: false,
		},
		KubeProxy: &DisabledTool{
			Enabled: false,
		},
		KubeStateMetricsScrap: &DisabledTool{
			Enabled: false,
		},
		KubeStateMetrics: &DisabledTool{
			Enabled: false,
		},
		NodeExporter: &DisabledTool{
			Enabled: false,
		},
		PrometheusNodeExporter: &DisabledTool{
			Enabled: false,
		},
		PrometheusOperator: &PrometheusOperatorValues{
			Enabled: true,
			TLSProxy: &TLSProxy{
				Enabled: false,
				Image: &Image{
					Repository: "squareup/ghostunnel",
					Tag:        imageTags["squareup/ghostunnel"],
					PullPolicy: "IfNotPresent",
				},
			},
			AdmissionWebhooks: &AdmissionWebhooks{
				FailurePolicy: "Fail",
				Enabled:       false,
				Patch: &Patch{
					Enabled: false,
					Image: &Image{
						Repository: "jettech/kube-webhook-certgen",
						Tag:        imageTags["jettech/kube-webhook-certgen"],
						PullPolicy: "IfNotPresent",
					},
					PriorityClassName: "",
				},
			},
			DenyNamespaces: []string{},
			ServiceAccount: &ServiceAccount{
				Create: true,
			},
			Service: &Service{
				NodePort:    30080,
				NodePortTLS: 30443,
				Type:        "ClusterIP",
			},
			CreateCustomResource:  true,
			CrdAPIGroup:           "monitoring.coreos.com",
			CleanupCustomResource: false,
			KubeletService: &KubeletService{
				Enabled:   false,
				Namespace: "kube-system",
			},
			ServiceMonitor: &ServiceMonitor{
				Interval:    "",
				SelfMonitor: false,
			},
			SecurityContext: &SecurityContext{
				RunAsNonRoot: true,
				RunAsUser:    65534,
			},
			Image: &Image{
				Repository: "quay.io/coreos/prometheus-operator",
				Tag:        imageTags["quay.io/coreos/prometheus-operator"],
				PullPolicy: "IfNotPresent",
			},
			ConfigmapReloadImage: &Image{
				Repository: "quay.io/coreos/configmap-reload",
				Tag:        imageTags["quay.io/coreos/configmap-reload"],
				PullPolicy: "IfNotPresent",
			},
			PrometheusConfigReloaderImage: &Image{
				Repository: "quay.io/coreos/prometheus-config-reloader",
				Tag:        imageTags["quay.io/coreos/prometheus-config-reloader"],
				PullPolicy: "IfNotPresent",
			},
			ConfigReloaderCPU:    "100m",
			ConfigReloaderMemory: "25Mi",

			HyperkubeImage: &Image{
				Repository: "k8s.gcr.io/hyperkube",
				Tag:        imageTags["k8s.gcr.io/hyperkube"],
				PullPolicy: "IfNotPresent",
			},
		},
		Prometheus: &DisabledTool{
			Enabled: false,
		},
	}
	return values
}
