package ambassador

type Resource struct {
	Name                     string `yaml:"name"`
	TargetAverageUtilization int    `yaml:"targetAverageUtilization"`
}

type Metric struct {
	Type     string    `yaml:"type"`
	Resource *Resource `yaml:"resource"`
}

type Autoscaling struct {
	Enabled     bool      `yaml:"enabled"`
	MinReplicas int       `yaml:"minReplicas"`
	MaxReplicas int       `yaml:"maxReplicas"`
	Metrics     []*Metric `yaml:"metrics"`
}

type SecurityContext struct {
	RunAsUser int `yaml:"runAsUser"`
}

type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}

type Port struct {
	Name       string `yaml:"name"`
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
	NodePort   int    `yaml:"nodePort"`
}

type Service struct {
	Type        string            `yaml:"type"`
	Ports       []*Port           `yaml:"ports"`
	Annotations map[string]string `yaml:"annotations"`
}

type AdminService struct {
	Create      bool              `yaml:"create"`
	Type        string            `yaml:"type"`
	Port        int               `yaml:"port"`
	Annotations map[string]string `yaml:"annotations"`
}

type Rbac struct {
	Create              bool `yaml:"create"`
	PodSecurityPolicies struct {
	} `yaml:"podSecurityPolicies"`
}

type Scope struct {
	SingleNamespace bool `yaml:"singleNamespace"`
}

type ServiceAccount struct {
	Create bool        `yaml:"create"`
	Name   interface{} `yaml:"name"`
}
type Crds struct {
	Enabled bool `yaml:"enabled"`
	Create  bool `yaml:"create"`
	Keep    bool `yaml:"keep"`
}

type Pro struct {
	Enabled bool   `yaml:"enabled"`
	Image   *Image `yaml:"image"`
	Ports   struct {
		Auth      int `yaml:"auth"`
		Ratelimit int `yaml:"ratelimit"`
	} `yaml:"ports"`
	LogLevel   string `yaml:"logLevel"`
	LicenseKey struct {
		Value  string `yaml:"value"`
		Secret struct {
			Enabled bool `yaml:"enabled"`
			Create  bool `yaml:"create"`
		} `yaml:"secret"`
	} `yaml:"licenseKey"`
	Env struct {
	} `yaml:"env"`
	Resources struct {
	} `yaml:"resources"`
	AuthService struct {
		Enabled                bool        `yaml:"enabled"`
		OptionalConfigurations interface{} `yaml:"optional_configurations"`
	} `yaml:"authService"`
	RateLimit struct {
		Enabled bool `yaml:"enabled"`
		Redis   struct {
			Annotations struct {
				Deployment struct {
				} `yaml:"deployment"`
				Service struct {
				} `yaml:"service"`
			} `yaml:"annotations"`
			Resources struct {
			} `yaml:"resources"`
		} `yaml:"redis"`
	} `yaml:"rateLimit"`
	DevPortal struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"devPortal"`
}

type PrometheusExporter struct {
	Enabled    bool   `yaml:"enabled"`
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
	Resources  struct {
	} `yaml:"resources"`
	Configuration string `yaml:"configuration"`
}

type Namespace struct {
	Name string `yaml:"name,omitempty"`
}

type Values struct {
	ReplicaCount          int                 `yaml:"replicaCount"`
	DaemonSet             bool                `yaml:"daemonSet"`
	Autoscaling           *Autoscaling        `yaml:"autoscaling"`
	PodDisruptionBudget   struct{}            `yaml:"podDisruptionBudget"`
	Namespace             *Namespace          `yaml:"namespace"`
	Env                   map[string]string   `yaml:"env"`
	ImagePullSecrets      []interface{}       `yaml:"imagePullSecrets"`
	SecurityContext       *SecurityContext    `yaml:"securityContext"`
	Image                 *Image              `yaml:"image"`
	NameOverride          string              `yaml:"nameOverride"`
	FullnameOverride      string              `yaml:"fullnameOverride,omitempty"`
	DNSPolicy             string              `yaml:"dnsPolicy"`
	HostNetwork           bool                `yaml:"hostNetwork"`
	Service               *Service            `yaml:"service"`
	AdminService          *AdminService       `yaml:"adminService"`
	Rbac                  *Rbac               `yaml:"rbac"`
	Scope                 *Scope              `yaml:"scope"`
	ServiceAccount        *ServiceAccount     `yaml:"serviceAccount"`
	InitContainers        []interface{}       `yaml:"initContainers"`
	SidecarContainers     []interface{}       `yaml:"sidecarContainers"`
	Volumes               []interface{}       `yaml:"volumes"`
	VolumeMounts          []interface{}       `yaml:"volumeMounts"`
	PodLabels             map[string]string   `yaml:"podLabels"`
	PodAnnotations        map[string]string   `yaml:"podAnnotations"`
	DeploymentAnnotations map[string]string   `yaml:"deploymentAnnotations"`
	Resources             struct{}            `yaml:"resources"`
	PriorityClassName     string              `yaml:"priorityClassName"`
	NodeSelector          struct{}            `yaml:"nodeSelector"`
	Tolerations           []interface{}       `yaml:"tolerations"`
	Affinity              struct{}            `yaml:"affinity"`
	AmbassadorConfig      string              `yaml:"ambassadorConfig"`
	Crds                  *Crds               `yaml:"crds"`
	Pro                   *Pro                `yaml:"pro"`
	PrometheusExporter    *PrometheusExporter `yaml:"prometheusExporter"`
}
