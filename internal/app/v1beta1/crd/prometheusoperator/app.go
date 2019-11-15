package prometheusoperator

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
	applicationName      = "prometheus-operator"
	resultsDirectoryName = "results"
	resultsFilename      = "results.yaml"
	defaultNamespace     = "monitoring"
)

type PrometheusOperator struct {
	ApplicationDirectoryPath string
}

func New(toolsDirectoryPath string) *PrometheusOperator {
	lo := &PrometheusOperator{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
	}

	return lo
}

func (p *PrometheusOperator) Reconcile(overlay string, helm *template.Helm, spec *toolsetsv1beta1.PrometheusOperator) error {

	logging.Log("CRD-2JlElqA8Zqu7wcw").Infof("Reconciling application %s", applicationName)
	resultsFileDirectory := filepath.Join(p.ApplicationDirectoryPath, resultsDirectoryName, overlay)
	_ = os.RemoveAll(resultsFileDirectory)
	_ = os.MkdirAll(resultsFileDirectory, os.ModePerm)
	resultFilePath := filepath.Join(resultsFileDirectory, resultsFilename)

	values := specToValues(helm.GetImageTags(applicationName), spec)

	writeValues := func(path string) error {
		if err := helper.StructToYaml(values, path); err != nil {
			logging.Log("CRD-l4CrOkKeQcIFehd").Debugf("Failed to write values file overlay %s application %s", overlay, applicationName)
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

	if spec.Deploy {
		kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath)
		if err := kubectlCmd.Run(); err != nil {
			logging.Log("CRD-hOW8nm0IlUQdpSj").OnError(err).Debugf("Failed to apply file %s", resultFilePath)
			return err
		}
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
