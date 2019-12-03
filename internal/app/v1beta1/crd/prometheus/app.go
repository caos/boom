package prometheus

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
	"github.com/caos/toolsop/internal/app/v1beta1/crd/prometheus/servicemonitor"
	"github.com/caos/toolsop/internal/app/v1beta1/crd/prometheusoperator"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/kubectl"
	"github.com/caos/toolsop/internal/template"
)

var (
	applicationName      = "prometheus"
	resultsDirectoryName = "results"
	resultsFileName      = "results.yaml"
	defaultNamespace     = "monitoring"
)

type Prometheus struct {
	ApplicationDirectoryPath string
	logger                   logging.Logger
}

func New(logger logging.Logger, toolsDirectoryPath string) *Prometheus {

	p := &Prometheus{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
		logger:                   logger,
	}

	return p
}

func (p *Prometheus) Reconcile(overlay string, helm *template.Helm, spec *toolsetsv1beta1.Prometheus) error {

	logFields := map[string]interface{}{
		"application": applicationName,
	}
	logFields["logID"] = "CRD-G9XDSdcIChwHP4w"
	p.logger.WithFields(logFields).Info("Reconciling")

	resultsFileDirectory := filepath.Join(p.ApplicationDirectoryPath, resultsDirectoryName, overlay)
	_ = os.RemoveAll(resultsFileDirectory)
	_ = os.MkdirAll(resultsFileDirectory, os.ModePerm)
	resultFilePath := filepath.Join(resultsFileDirectory, resultsFileName)

	values, err := specToValues(helm.GetImageTags(applicationName), spec)
	if err != nil {
		return err
	}

	writeValues := func(path string) error {
		if err := errors.Wrapf(helper.StructToYaml(values, path), "Failed to write values file overlay %s application %s", overlay, applicationName); err != nil {
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
		if err := errors.Wrapf(helper.Run(p.logger, kubectlCmd.Build()), "Failed to apply file %s", resultFilePath); err != nil {
			return err
		}
	}

	return nil
}

func specToValues(imageTags map[string]string, spec *toolsetsv1beta1.Prometheus) (*Values, error) {
	promValues := &PrometheusValues{
		Enabled: true,
		ServiceAccount: &ServiceAccount{
			Create: true,
		},
		Service: &Service{
			Port:       9090,
			TargetPort: 9090,
			NodePort:   30090,
			Type:       "ClusterIP",
		},
		ServicePerReplica: &ServicePerReplica{
			Enabled:    false,
			Port:       9090,
			TargetPort: 9090,
			NodePort:   30091,
		},
		PodDisruptionBudget: &PodDisruptionBudget{
			Enabled:      false,
			MinAvailable: 1,
		},
		Ingress: &Ingress{
			Enabled: false,
		},
		IngressPerReplica: &IngressPerReplica{
			Enabled: false,
		},
		PodSecurityPolicy: &PodSecurityPolicy{},
		ServiceMonitor: &ServiceMonitor{
			SelfMonitor: true,
		},
		PrometheusSpec: &PrometheusSpec{
			Image: &Image{
				Repository: "quay.io/prometheus/prometheus",
				Tag:        imageTags["quay.io/prometheus/prometheus"],
			},
			RuleSelectorNilUsesHelmValues:           true,
			ServiceMonitorSelectorNilUsesHelmValues: true,
			PodMonitorSelectorNilUsesHelmValues:     true,
			Retention:                               "10d",
			Replicas:                                1,
			LogLevel:                                "info",
			LogFormat:                               "logfmt",
			RoutePrefix:                             "/",
			PodAntiAffinityTopologyKey:              "kubernetes.io/hostname",
			SecurityContext: &SecurityContext{
				RunAsNonRoot: true,
				RunAsUser:    1000,
				FsGroup:      2000,
			},
		},
	}

	if spec.MonitorLabels != nil {
		promValues.PrometheusSpec.ServiceMonitorSelector = &MonitorSelector{
			MatchLabels: spec.MonitorLabels,
		}
		additionalServiceMonitors := make([]*servicemonitor.Values, 0)
		for _, specServiceMonitor := range spec.ServiceMonitors {
			valuesServiceMonitor := servicemonitor.SpecToValues(spec.MonitorLabels, specServiceMonitor)
			additionalServiceMonitors = append(additionalServiceMonitors, valuesServiceMonitor)
		}

		promValues.AdditionalServiceMonitors = additionalServiceMonitors
	}

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
		PrometheusOperator: &prometheusoperator.PrometheusOperatorValues{
			Enabled: false,
			TLSProxy: &prometheusoperator.TLSProxy{
				Enabled: false,
				Image: &prometheusoperator.Image{
					Repository: "squareup/ghostunnel",
					Tag:        imageTags["squareup/ghostunnel"],
					PullPolicy: "IfNotPresent",
				},
			},
			AdmissionWebhooks: &prometheusoperator.AdmissionWebhooks{
				FailurePolicy: "Fail",
				Enabled:       false,
				Patch: &prometheusoperator.Patch{
					Enabled: false,
					Image: &prometheusoperator.Image{
						Repository: "jettech/kube-webhook-certgen",
						Tag:        imageTags["jettech/kube-webhook-certgen"],
						PullPolicy: "IfNotPresent",
					},
					PriorityClassName: "",
				},
			},
			ServiceAccount: &prometheusoperator.ServiceAccount{
				Create: false,
			},
			ServiceMonitor: &prometheusoperator.ServiceMonitor{
				Interval:    "",
				SelfMonitor: false,
			},
			CreateCustomResource: true,
		},
		Prometheus: promValues,
	}

	return values, nil
}
