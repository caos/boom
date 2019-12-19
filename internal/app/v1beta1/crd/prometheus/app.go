package prometheus

import (
	"path/filepath"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	"github.com/caos/boom/internal/app/v1beta1/crd/defaults"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheus/servicemonitor"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheusoperator"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/template"
)

var (
	applicationName = "prometheus"
)

type Config struct {
	Prefix          string
	Namespace       string
	MonitorLabels   map[string]string
	ServiceMonitors []*servicemonitor.Config
	ReplicaCount    int
	StorageSpec     *ConfigStorageSpec
}

type ConfigStorageSpec struct {
	StorageClass string
	AccessModes  []string
	Storage      string
}

type Prometheus struct {
	ApplicationDirectoryPath string
	logger                   logging.Logger
	config                   *Config
}

func New(logger logging.Logger, toolsDirectoryPath string) *Prometheus {

	p := &Prometheus{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
		logger:                   logger,
	}

	return p
}

func (p *Prometheus) Reconcile(overlay, specNamespace string, helm *template.Helm, config *Config) error {

	logFields := map[string]interface{}{
		"application": applicationName,
		"logID":       "CRD-G9XDSdcIChwHP4w",
	}
	p.logger.WithFields(logFields).Info("Reconciling")

	resultFilePath := defaults.GetResultFilePath(overlay, p.ApplicationDirectoryPath, applicationName)
	prefix := defaults.GetPrefix(overlay, applicationName, config.Prefix)
	namespace := defaults.GetNamespace(overlay, applicationName, specNamespace, config.Namespace)
	values, err := specToValues(helm.GetImageTags(applicationName), config)
	if err != nil {
		return err
	}

	writeValues := func(path string) error {
		if err := errors.Wrapf(helper.StructToYaml(values, path), "Failed to write values file overlay %s application %s", overlay, applicationName); err != nil {
			return err
		}
		return nil
	}

	if err := helm.PrepareTemplate(applicationName, prefix, namespace, writeValues); err != nil {
		return err
	}

	if config != nil {

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
		if err := errors.Wrapf(helper.Run(p.logger, kubectlCmd.Build()), "Failed to apply file %s", resultFilePath); err != nil {
			return err
		}

		p.config = config
	} else if config == nil && p.config != nil {
		kubectlCmd := kubectl.New("delete").AddParameter("-f", resultFilePath)
		if err := errors.Wrapf(helper.Run(p.logger, kubectlCmd.Build()), "Failed to apply file %s", resultFilePath); err != nil {
			return err
		}

		p.config = nil
	}

	return nil
}

func specToValues(imageTags map[string]string, config *Config) (*Values, error) {
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
			SelfMonitor: false,
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

	if config.StorageSpec != nil {
		storageSpec := &StorageSpec{
			VolumeClaimTemplate: &VolumeClaimTemplate{
				Spec: &VolumeClaimTemplateSpec{
					StorageClassName: config.StorageSpec.StorageClass,
					AccessModes:      config.StorageSpec.AccessModes,
					Resources: &Resources{
						Requests: &Request{
							Storage: config.StorageSpec.Storage,
						},
					},
				},
			},
		}

		promValues.PrometheusSpec.StorageSpec = storageSpec
	}

	if config.MonitorLabels != nil {
		promValues.PrometheusSpec.ServiceMonitorSelector = &MonitorSelector{
			MatchLabels: config.MonitorLabels,
		}
	}

	if config.ServiceMonitors != nil {
		additionalServiceMonitors := make([]*servicemonitor.Values, 0)
		for _, specServiceMonitor := range config.ServiceMonitors {
			valuesServiceMonitor := servicemonitor.SpecToValues(specServiceMonitor)
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
