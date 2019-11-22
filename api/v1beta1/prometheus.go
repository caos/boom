package v1beta1

type Prometheus struct {
	Prefix          string            `json:"prefix" yaml:"prefix"`
	Namespace       string            `json:"namespace" yaml:"namespace"`
	Deploy          bool              `json:"deploy" yaml:"deploy"`
	MonitorLabels   map[string]string `json:"monitorLabel" yaml:"monitorLabel"`
	ServiceMonitors []*ServiceMonitor `json:"serviceMonitors" yaml:"serviceMonitors"`
	// Annotations           map[string]string `json:"annotations,omitempty"`
	// RemoteWrite           []*RemoteWrite    `json:"remotewrite,omitempty"`
}

type Relabeling struct {
	Action       string   `json:"action,omitempty"`
	Regex        string   `json:"regex,omitempty"`
	Replacement  string   `json:"replacement,omitempty"`
	SourceLabels []string `json:"sourcelabels,omitempty"`
	TargetLabel  string   `json:"targetlabel,omitempty"`
}

type ServiceMonitor struct {
	Name                  string            `json:"name,omitempty"`
	Interval              string            `json:"interval,omitempty"`
	Relabelings           []*Relabeling     `json:"relabelings,omitempty"`
	ServiceMatchingLabels map[string]string `json:"serviceMatchingLabels,omitempty"`
	Endpoints             []*Endpoint       `json:"endpoints,omitempty"`
	// Annotations map[string]string `json:"annotations,omitempty"`
}

type Endpoint struct {
	Port       string `json:"port,omitempty"`
	TargetPort string `json:"targetPort,omitempty"`
	Interval   string `json:"interval,omitempty"`
	Path       string `json:"path,omitempty"`
	Scheme     string `json:"scheme,omitempty"`
}
