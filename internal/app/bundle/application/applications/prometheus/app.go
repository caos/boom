package prometheus

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	"github.com/caos/boom/api/v1beta1"
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/app/bundle/application/applications/prometheusoperator"
	"github.com/caos/boom/internal/app/name"
)

const (
	applicationName name.Application = "prometheus"
)

func GetName() name.Application {
	return applicationName
}

type Prometheus struct {
	logger logging.Logger
	config *Config
}

func New(logger logging.Logger) *Prometheus {
	return &Prometheus{
		logger: logger,
	}
}
func (p *Prometheus) GetName() name.Application {
	return applicationName
}

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	config := ScrapeMetricsCrdsConfig(toolsetCRDSpec)
	if config == nil {
		return false
	}
	return true
}

func (p *Prometheus) Initial() bool {
	return p.config == nil
}

func (p *Prometheus) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	config := ScrapeMetricsCrdsConfig(toolsetCRDSpec)
	return !reflect.DeepEqual(config, p.config)
}

func (p *Prometheus) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	if toolsetCRDSpec == nil {
		return
	}
	p.config = ScrapeMetricsCrdsConfig(toolsetCRDSpec)
}

func (p *Prometheus) GetNamespace() string {
	return "caos-system"
}

func (p *Prometheus) SpecToHelmValues(toolsetCRDSpec *v1beta1.ToolsetSpec) interface{} {
	config := ScrapeMetricsCrdsConfig(toolsetCRDSpec)

	values := defaultValues(p.GetImageTags())

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

		values.Prometheus.PrometheusSpec.StorageSpec = storageSpec
	}

	if config.MonitorLabels != nil {
		values.Prometheus.PrometheusSpec.ServiceMonitorSelector = &MonitorSelector{
			MatchLabels: config.MonitorLabels,
		}
	}

	if config.ServiceMonitors != nil {
		additionalServiceMonitors := make([]*servicemonitor.Values, 0)
		for _, specServiceMonitor := range config.ServiceMonitors {
			valuesServiceMonitor := servicemonitor.SpecToValues(specServiceMonitor)
			additionalServiceMonitors = append(additionalServiceMonitors, valuesServiceMonitor)
		}

		values.Prometheus.AdditionalServiceMonitors = additionalServiceMonitors
	}

	if config.AdditionalScrapeConfigs != nil {
		values.Prometheus.PrometheusSpec.AdditionalScrapeConfigs = config.AdditionalScrapeConfigs
	}

	ruleLabels := map[string]string{"prometheus": "caos", "instance-name": "operated"}
	rules, _ := getRules(ruleLabels)

	values.Prometheus.PrometheusSpec.RuleSelector = &RuleSelector{MatchLabels: ruleLabels}
	values.DefaultRules.Labels = ruleLabels
	values.KubeTargetVersionOverride = config.KubeVersion
	values.AdditionalPrometheusRules = []*AdditionalPrometheusRules{rules}

	return values
}

func defaultValues(imageTags map[string]string) *Values {
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

	return &Values{
		FullnameOverride: "operated",
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
			KubeletService: &prometheusoperator.KubeletService{
				Enabled: false,
			},
		},
		Prometheus: promValues,
	}
}
