package helm

type Rbac struct {
	Create bool `yaml:"create"`
}

type PodSecurityPolicy struct {
	Enabled bool `yaml:"enabled"`
}

type LeaderElection struct {
	Namespace string `yaml:"namespace"`
}

type Global struct {
	ImagePullSecrets  []interface{}      `yaml:"imagePullSecrets"`
	IsOpenshift       bool               `yaml:"isOpenshift"`
	PriorityClassName string             `yaml:"priorityClassName"`
	Rbac              *Rbac              `yaml:"rbac"`
	PodSecurityPolicy *PodSecurityPolicy `yaml:"podSecurityPolicy"`
	LogLevel          int                `yaml:"logLevel"`
	LeaderElection    *LeaderElection    `yaml:"leaderElection"`
}

type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}

type ServiceAccount struct {
	Create      bool              `yaml:"create"`
	Name        interface{}       `yaml:"name"`
	Annotations map[string]string `yaml:"annotations"`
}

type SecurityContext struct {
	Enabled   bool `yaml:"enabled"`
	FsGroup   int  `yaml:"fsGroup"`
	RunAsUser int  `yaml:"runAsUser"`
}

type Webhook struct {
	Enabled           bool              `yaml:"enabled"`
	ReplicaCount      int               `yaml:"replicaCount"`
	Strategy          struct{}          `yaml:"strategy"`
	PodAnnotations    map[string]string `yaml:"podAnnotations"`
	ExtraArgs         []interface{}     `yaml:"extraArgs"`
	Resources         struct{}          `yaml:"resources"`
	NodeSelector      struct{}          `yaml:"nodeSelector"`
	Affinity          struct{}          `yaml:"affinity"`
	Tolerations       []interface{}     `yaml:"tolerations"`
	Image             *Image            `yaml:"image"`
	InjectAPIServerCA bool              `yaml:"injectAPIServerCA"`
	SecurePort        int               `yaml:"securePort"`
}

type Cainjector struct {
	ReplicaCount   int               `yaml:"replicaCount"`
	Strategy       struct{}          `yaml:"strategy"`
	PodAnnotations map[string]string `yaml:"podAnnotations"`
	ExtraArgs      []interface{}     `yaml:"extraArgs"`
	Resources      struct{}          `yaml:"resources"`
	NodeSelector   struct{}          `yaml:"nodeSelector"`
	Affinity       struct{}          `yaml:"affinity"`
	Tolerations    []interface{}     `yaml:"tolerations"`
	Image          *Image            `yaml:"image"`
}

type Prometheus struct {
	Enabled        bool `yaml:"enabled"`
	Servicemonitor struct {
		Enabled            bool   `yaml:"enabled"`
		PrometheusInstance string `yaml:"prometheusInstance"`
		TargetPort         int    `yaml:"targetPort"`
		Path               string `yaml:"path"`
		Interval           string `yaml:"interval"`
		ScrapeTimeout      string `yaml:"scrapeTimeout"`
		Labels             struct {
		} `yaml:"labels"`
	} `yaml:"servicemonitor"`
}

type Values struct {
	FullnameOverride         string            `yaml:"fullnameOverride,omitempty"`
	Global                   *Global           `yaml:"global"`
	ReplicaCount             int               `yaml:"replicaCount"`
	Strategy                 struct{}          `yaml:"strategy"`
	Image                    *Image            `yaml:"image"`
	ClusterResourceNamespace string            `yaml:"clusterResourceNamespace"`
	ServiceAccount           *ServiceAccount   `yaml:"serviceAccount"`
	ExtraArgs                []interface{}     `yaml:"extraArgs"`
	ExtraEnv                 []interface{}     `yaml:"extraEnv"`
	Resources                struct{}          `yaml:"resources"`
	SecurityContext          *SecurityContext  `yaml:"securityContext"`
	PodAnnotations           map[string]string `yaml:"podAnnotations"`
	PodLabels                map[string]string `yaml:"podLabels"`
	NodeSelector             struct{}          `yaml:"nodeSelector"`
	IngressShim              struct{}          `yaml:"ingressShim"`
	Prometheus               *Prometheus       `yaml:"prometheus"`
	Affinity                 struct{}          `yaml:"affinity"`
	Tolerations              []interface{}     `yaml:"tolerations"`
	Webhook                  *Webhook          `yaml:"webhook"`
	Cainjector               *Cainjector       `yaml:"cainjector"`
}
