package prometheus

import (
	"github.com/caos/toolsop/internal/app/v1beta1/crd/prometheus/servicemonitor"
	"github.com/caos/toolsop/internal/app/v1beta1/crd/prometheusoperator"
)

type Service struct {
	Annotations              map[string]string `yaml:"annotations"`
	Labels                   map[string]string `yaml:"labels"`
	ClusterIP                string            `yaml:"clusterIP"`
	Port                     int               `yaml:"port"`
	TargetPort               int               `yaml:"targetPort"`
	ExternalIPs              []interface{}     `yaml:"externalIPs"`
	NodePort                 int               `yaml:"nodePort"`
	LoadBalancerIP           string            `yaml:"loadBalancerIP"`
	LoadBalancerSourceRanges []interface{}     `yaml:"loadBalancerSourceRanges"`
	Type                     string            `yaml:"type"`
	SessionAffinity          string            `yaml:"sessionAffinity"`
}

type ServicePerReplica struct {
	Enabled                  bool              `yaml:"enabled"`
	Annotations              map[string]string `yaml:"annotations"`
	Port                     int               `yaml:"port"`
	TargetPort               int               `yaml:"targetPort"`
	NodePort                 int               `yaml:"nodePort"`
	LoadBalancerSourceRanges []interface{}     `yaml:"loadBalancerSourceRanges"`
	Type                     string            `yaml:"type"`
}
type PodDisruptionBudget struct {
	Enabled        bool   `yaml:"enabled"`
	MinAvailable   int    `yaml:"minAvailable"`
	MaxUnavailable string `yaml:"maxUnavailable"`
}

type Ingress struct {
	Enabled     bool              `yaml:"enabled"`
	Annotations map[string]string `yaml:"annotations"`
	Labels      map[string]string `yaml:"labels"`
	Hosts       []interface{}     `yaml:"hosts"`
	Paths       []interface{}     `yaml:"paths"`
	TLS         []interface{}     `yaml:"tls"`
}
type IngressPerReplica struct {
	Enabled       bool              `yaml:"enabled"`
	Annotations   map[string]string `yaml:"annotations"`
	Labels        map[string]string `yaml:"labels"`
	HostPrefix    string            `yaml:"hostPrefix"`
	HostDomain    string            `yaml:"hostDomain"`
	Paths         []interface{}     `yaml:"paths"`
	TLSSecretName string            `yaml:"tlsSecretName"`
}

type PodSecurityPolicy struct {
	AllowedCapabilities []interface{} `yaml:"allowedCapabilities"`
}

type ServiceMonitor struct {
	Interval          string        `yaml:"interval"`
	SelfMonitor       bool          `yaml:"selfMonitor"`
	BearerTokenFile   interface{}   `yaml:"bearerTokenFile"`
	MetricRelabelings []interface{} `yaml:"metricRelabelings"`
	Relabelings       []interface{} `yaml:"relabelings"`
}

type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

type NamespaceSelector struct {
	Any        bool     `yaml:"any"`
	MatchNames []string `yaml:"matchNames"`
}

type MonitorSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

type Query struct {
	LookbackDelta  string `yaml:"lookbackDelta"`
	MaxConcurrency int32  `yaml:"maxConcurrency"`
	MaxSamples     int32  `yaml:"maxSamples"`
	Timeout        string `yaml:"timeout"`
}
type PodMetadata struct {
	Labels map[string]string `yaml:"labels"`
}
type StorageSpec struct {
	VolumeClaimTemplate *VolumeClaimTemplate `yaml:"volumeClaimTemplate"`
}
type VolumeClaimTemplate struct {
	Spec     *VolumeClaimTemplateSpec `yaml:"spec"`
	selector struct{}                 `yaml:"selector"`
}
type VolumeClaimTemplateSpec struct {
	StorageClassName string     `yaml:"storageClassName"`
	AccessModes      []string   `yaml:"accessModes"`
	Resources        *Resources `yaml:"resources"`
}
type Resources struct {
	Requests *Request `yaml:"request"`
}

type Request struct {
	Storage string `yaml:"storage"`
}
type SecurityContext struct {
	RunAsNonRoot bool `yaml:"runAsNonRoot"`
	RunAsUser    int  `yaml:"runAsUser"`
	FsGroup      int  `yaml:"fsGroup"`
}
type PrometheusSpec struct {
	ScrapeInterval                          string             `yaml:"scrapeInterval"`
	EvaluationInterval                      string             `yaml:"evaluationInterval"`
	ListenLocal                             bool               `yaml:"listenLocal"`
	EnableAdminAPI                          bool               `yaml:"enableAdminAPI"`
	Image                                   *Image             `yaml:"image"`
	Tolerations                             []interface{}      `yaml:"tolerations"`
	AlertingEndpoints                       []interface{}      `yaml:"alertingEndpoints"`
	ExternalLabels                          map[string]string  `yaml:"externalLabels"`
	ReplicaExternalLabelName                string             `yaml:"replicaExternalLabelName"`
	ReplicaExternalLabelNameClear           bool               `yaml:"replicaExternalLabelNameClear"`
	PrometheusExternalLabelName             string             `yaml:"prometheusExternalLabelName"`
	PrometheusExternalLabelNameClear        bool               `yaml:"prometheusExternalLabelNameClear"`
	ExternalURL                             string             `yaml:"externalUrl"`
	NodeSelector                            map[string]string  `yaml:"nodeSelector"`
	Secrets                                 []interface{}      `yaml:"secrets"`
	ConfigMaps                              []interface{}      `yaml:"configMaps"`
	Query                                   *Query             `yaml:"query"`
	RuleNamespaceSelector                   *NamespaceSelector `yaml:"ruleNamespaceSelector"`
	RuleSelectorNilUsesHelmValues           bool               `yaml:"ruleSelectorNilUsesHelmValues"`
	RuleSelector                            struct{}           `yaml:"ruleSelector"`
	ServiceMonitorSelectorNilUsesHelmValues bool               `yaml:"serviceMonitorSelectorNilUsesHelmValues"`
	ServiceMonitorSelector                  *MonitorSelector   `yaml:"serviceMonitorSelector"`
	ServiceMonitorNamespaceSelector         *NamespaceSelector `yaml:"serviceMonitorNamespaceSelector"`
	PodMonitorSelectorNilUsesHelmValues     bool               `yaml:"podMonitorSelectorNilUsesHelmValues"`
	PodMonitorSelector                      *MonitorSelector   `yaml:"podMonitorSelector"`
	PodMonitorNamespaceSelector             *NamespaceSelector `yaml:"podMonitorNamespaceSelector"`
	Retention                               string             `yaml:"retention"`
	RetentionSize                           string             `yaml:"retentionSize"`
	WalCompression                          bool               `yaml:"walCompression"`
	Paused                                  bool               `yaml:"paused"`
	Replicas                                int                `yaml:"replicas"`
	LogLevel                                string             `yaml:"logLevel"`
	LogFormat                               string             `yaml:"logFormat"`
	RoutePrefix                             string             `yaml:"routePrefix"`
	PodMetadata                             *PodMetadata       `yaml:"podMetadata"`
	PodAntiAffinity                         string             `yaml:"podAntiAffinity"`
	PodAntiAffinityTopologyKey              string             `yaml:"podAntiAffinityTopologyKey"`
	Affinity                                struct{}           `yaml:"affinity"`
	RemoteRead                              []interface{}      `yaml:"remoteRead"`
	RemoteWrite                             []interface{}      `yaml:"remoteWrite"`
	RemoteWriteDashboards                   bool               `yaml:"remoteWriteDashboards"`
	Resources                               struct{}           `yaml:"resources"`
	StorageSpec                             *StorageSpec       `yaml:"storageSpec"`
	AdditionalScrapeConfigs                 []interface{}      `yaml:"additionalScrapeConfigs"`
	AdditionalAlertManagerConfigs           []interface{}      `yaml:"additionalAlertManagerConfigs"`
	AdditionalAlertRelabelConfigs           []interface{}      `yaml:"additionalAlertRelabelConfigs"`
	SecurityContext                         *SecurityContext   `yaml:"securityContext"`
	PriorityClassName                       string             `yaml:"priorityClassName"`
	Thanos                                  struct{}           `yaml:"thanos"`
	Containers                              []interface{}      `yaml:"containers"`
	AdditionalScrapeConfigsExternal         bool               `yaml:"additionalScrapeConfigsExternal"`
}

type ServiceAccount struct {
	Create bool   `yaml:"create"`
	Name   string `yaml:"name"`
}

type PrometheusValues struct {
	Enabled                   bool                     `yaml:"enabled"`
	Annotations               map[string]string        `yaml:"annotations"`
	ServiceAccount            *ServiceAccount          `yaml:"serviceAccount"`
	Service                   *Service                 `yaml:"service"`
	ServicePerReplica         *ServicePerReplica       `yaml:"servicePerReplica"`
	PodDisruptionBudget       *PodDisruptionBudget     `yaml:"podDisruptionBudget"`
	Ingress                   *Ingress                 `yaml:"ingress"`
	IngressPerReplica         *IngressPerReplica       `yaml:"ingressPerReplica"`
	PodSecurityPolicy         *PodSecurityPolicy       `yaml:"podSecurityPolicy"`
	ServiceMonitor            *ServiceMonitor          `yaml:"serviceMonitor"`
	PrometheusSpec            *PrometheusSpec          `yaml:"prometheusSpec"`
	AdditionalServiceMonitors []*servicemonitor.Values `yaml:"additionalServiceMonitors"`
	AdditionalPodMonitors     []interface{}            `yaml:"additionalPodMonitors"`
}

type Rules struct {
	Alertmanager                bool `yaml:"alertmanager"`
	Etcd                        bool `yaml:"etcd"`
	General                     bool `yaml:"general"`
	K8S                         bool `yaml:"k8s"`
	KubeApiserver               bool `yaml:"kubeApiserver"`
	KubePrometheusNodeAlerting  bool `yaml:"kubePrometheusNodeAlerting"`
	KubePrometheusNodeRecording bool `yaml:"kubePrometheusNodeRecording"`
	KubernetesAbsent            bool `yaml:"kubernetesAbsent"`
	KubernetesApps              bool `yaml:"kubernetesApps"`
	KubernetesResources         bool `yaml:"kubernetesResources"`
	KubernetesStorage           bool `yaml:"kubernetesStorage"`
	KubernetesSystem            bool `yaml:"kubernetesSystem"`
	KubeScheduler               bool `yaml:"kubeScheduler"`
	Network                     bool `yaml:"network"`
	Node                        bool `yaml:"node"`
	Prometheus                  bool `yaml:"prometheus"`
	PrometheusOperator          bool `yaml:"prometheusOperator"`
	Time                        bool `yaml:"time"`
}

type DefaultRules struct {
	Create      bool              `yaml:"create"`
	Rules       *Rules            `yaml:"rules"`
	Labels      map[string]string `yaml:"labels"`
	Annotations map[string]string `yaml:"annotations"`
}

type Global struct {
	Rbac             *Rbac         `yaml:"rbac"`
	ImagePullSecrets []interface{} `yaml:"imagePullSecrets"`
}

type Rbac struct {
	Create     bool `yaml:"create"`
	PspEnabled bool `yaml:"pspEnabled"`
}

type DisabledTool struct {
	Enabled bool `yaml:"enabled"`
}

type Values struct {
	NameOverride              string                                       `yaml:"nameOverride"`
	FullnameOverride          string                                       `yaml:"fullnameOverride"`
	CommonLabels              map[string]string                            `yaml:"commonLabels"`
	DefaultRules              *DefaultRules                                `yaml:"defaultRules"`
	AdditionalPrometheusRules []interface{}                                `yaml:"additionalPrometheusRules"`
	Global                    *Global                                      `yaml:"global"`
	Alertmanager              *DisabledTool                                `yaml:"alertmanager"`
	Grafana                   *DisabledTool                                `yaml:"grafana"`
	KubeAPIServer             *DisabledTool                                `yaml:"kubeApiServer"`
	Kubelet                   *DisabledTool                                `yaml:"kubelet"`
	KubeControllerManager     *DisabledTool                                `yaml:"kubeControllerManager"`
	CoreDNS                   *DisabledTool                                `yaml:"coreDns"`
	KubeDNS                   *DisabledTool                                `yaml:"kubeDns"`
	KubeEtcd                  *DisabledTool                                `yaml:"kubeEtcd"`
	KubeScheduler             *DisabledTool                                `yaml:"kubeScheduler"`
	KubeProxy                 *DisabledTool                                `yaml:"kubeProxy"`
	KubeStateMetricsScrap     *DisabledTool                                `yaml:"kubeStateMetrics"`
	KubeStateMetrics          *DisabledTool                                `yaml:"kube-state-metrics"`
	NodeExporter              *DisabledTool                                `yaml:"nodeExporter"`
	PrometheusNodeExporter    *DisabledTool                                `yaml:"prometheus-node-exporter"`
	PrometheusOperator        *prometheusoperator.PrometheusOperatorValues `yaml:"prometheusOperator"`
	Prometheus                *PrometheusValues                            `yaml:"prometheus"`
}
