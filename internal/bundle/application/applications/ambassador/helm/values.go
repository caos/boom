package helm

type AdminService struct {
	Annotations map[string]string `yaml:"annotations"`
	Create      bool              `yaml:"create"`
	Port        int               `yaml:"port"`
	Type        string            `yaml:"type"`
}
type AuthService struct {
	Create                 bool        `yaml:"create"`
	OptionalConfigurations interface{} `yaml:"optional_configurations"`
}
type Resource struct {
	Name                     string `yaml:"name"`
	TargetAverageUtilization int    `yaml:"targetAverageUtilization"`
}
type Metrics []struct {
	Resource *Resource `yaml:"resource"`
	Type     string    `yaml:"type"`
}
type Autoscaling struct {
	Enabled     bool     `yaml:"enabled"`
	MaxReplicas int      `yaml:"maxReplicas"`
	Metrics     *Metrics `yaml:"metrics"`
	MinReplicas int      `yaml:"minReplicas"`
}
type Crds struct {
	Create  bool `yaml:"create"`
	Enabled bool `yaml:"enabled"`
	Keep    bool `yaml:"keep"`
}
type DeploymentStrategy struct {
	Type string `yaml:"type"`
}
type Image struct {
	PullPolicy string `yaml:"pullPolicy"`
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}
type LicenseKey struct {
	CreateSecret bool        `yaml:"createSecret"`
	Value        interface{} `yaml:"value"`
}
type LivenessProbe struct {
	FailureThreshold    int `yaml:"failureThreshold"`
	InitialDelaySeconds int `yaml:"initialDelaySeconds"`
	PeriodSeconds       int `yaml:"periodSeconds"`
}
type PrometheusExporter struct {
	Enabled    bool     `yaml:"enabled"`
	PullPolicy string   `yaml:"pullPolicy"`
	Repository string   `yaml:"repository"`
	Resources  struct{} `yaml:"resources"`
	Tag        string   `yaml:"tag"`
}
type RateLimit struct {
	Create bool `yaml:"create"`
}
type Rbac struct {
	Create              bool     `yaml:"create"`
	PodSecurityPolicies struct{} `yaml:"podSecurityPolicies"`
}
type ReadinessProbe struct {
	FailureThreshold    int `yaml:"failureThreshold"`
	InitialDelaySeconds int `yaml:"initialDelaySeconds"`
	PeriodSeconds       int `yaml:"periodSeconds"`
}
type RedisAnnotations struct {
	Deployment map[string]string `yaml:"deployment"`
	Service    map[string]string `yaml:"service"`
}
type Redis struct {
	Annotations *RedisAnnotations `yaml:"annotations"`
	Create      bool              `yaml:"create"`
	Resources   struct{}          `yaml:"resources"`
}
type Scope struct {
	SingleNamespace bool `yaml:"singleNamespace"`
}
type SecurityContext struct {
	RunAsUser int `yaml:"runAsUser"`
}
type Port struct {
	Name       string `yaml:"name"`
	Port       uint16 `yaml:"port,omitempty"`
	TargetPort uint16 `yaml:"targetPort,omitempty"`
	NodePort   uint16 `yaml:"nodePort,omitempty"`
}
type Service struct {
	Annotations    interface{} `yaml:"annotations,omitempty"`
	Ports          []*Port     `yaml:"ports"`
	Type           string      `yaml:"type"`
	LoadBalancerIP string      `yaml:"loadBalancerIP,omitempty"`
}
type ServiceAccount struct {
	Create bool        `yaml:"create"`
	Name   interface{} `yaml:"name"`
}

type Values struct {
	AdminService          *AdminService       `yaml:"adminService"`
	Affinity              struct{}            `yaml:"affinity"`
	AmbassadorConfig      string              `yaml:"ambassadorConfig"`
	AuthService           *AuthService        `yaml:"authService"`
	Autoscaling           *Autoscaling        `yaml:"autoscaling"`
	Crds                  *Crds               `yaml:"crds"`
	DaemonSet             bool                `yaml:"daemonSet"`
	DeploymentAnnotations map[string]string   `yaml:"deploymentAnnotations"`
	DeploymentStrategy    *DeploymentStrategy `yaml:"deploymentStrategy"`
	DNSPolicy             string              `yaml:"dnsPolicy"`
	Env                   map[string]string   `yaml:"env"`
	FullnameOverride      string              `yaml:"fullnameOverride"`
	HostNetwork           bool                `yaml:"hostNetwork"`
	Image                 *Image              `yaml:"image"`
	ImagePullSecrets      []interface{}       `yaml:"imagePullSecrets"`
	InitContainers        []interface{}       `yaml:"initContainers"`
	LicenseKey            *LicenseKey         `yaml:"licenseKey"`
	LivenessProbe         *LivenessProbe      `yaml:"livenessProbe"`
	NameOverride          string              `yaml:"nameOverride"`
	NodeSelector          struct{}            `yaml:"nodeSelector"`
	PodAnnotations        map[string]string   `yaml:"podAnnotations"`
	PodDisruptionBudget   struct{}            `yaml:"podDisruptionBudget"`
	PodLabels             map[string]string   `yaml:"podLabels"`
	PriorityClassName     string              `yaml:"priorityClassName"`
	PrometheusExporter    *PrometheusExporter `yaml:"prometheusExporter"`
	RateLimit             *RateLimit          `yaml:"rateLimit"`
	Rbac                  *Rbac               `yaml:"rbac"`
	ReadinessProbe        *ReadinessProbe     `yaml:"readinessProbe"`
	Redis                 *Redis              `yaml:"redis"`
	RedisURL              interface{}         `yaml:"redisURL"`
	ReplicaCount          int                 `yaml:"replicaCount"`
	Resources             struct{}            `yaml:"resources"`
	Scope                 *Scope              `yaml:"scope"`
	SecurityContext       *SecurityContext    `yaml:"securityContext"`
	Service               *Service            `yaml:"service"`
	ServiceAccount        *ServiceAccount     `yaml:"serviceAccount"`
	SidecarContainers     []interface{}       `yaml:"sidecarContainers"`
	Tolerations           []interface{}       `yaml:"tolerations"`
	VolumeMounts          []interface{}       `yaml:"volumeMounts"`
	Volumes               []interface{}       `yaml:"volumes"`
}
