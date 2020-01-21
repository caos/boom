package v1beta1

type KubeStateMetrics struct {
	Deploy       bool        `json:"deploy,omitempty"`
	ReplicaCount int         `json:"replicaCount,omitempty"`
	Collectors   *Collectors `json:"collectors,omitempty"`
}

type Collectors struct {
	Certificatesigningrequests bool `json:"certificatesigningrequests,omitempty" yaml:"certificatesigningrequests"`
	Configmaps                 bool `json:"configmaps,omitempty" yaml:"configmaps"`
	Cronjobs                   bool `json:"cronjobs,omitempty" yaml:"cronjobs"`
	Daemonsets                 bool `json:"daemonsets,omitempty" yaml:"daemonsets"`
	Deployments                bool `json:"deployments,omitempty" yaml:"deployments"`
	Endpoints                  bool `json:"endpoints,omitempty" yaml:"endpoints"`
	Horizontalpodautoscalers   bool `json:"horizontalpodautoscalers,omitempty" yaml:"horizontalpodautoscalers"`
	Ingresses                  bool `json:"ingresses,omitempty" yaml:"ingresses"`
	Jobs                       bool `json:"jobs,omitempty" yaml:"jobs"`
	Limitranges                bool `json:"limitranges,omitempty" yaml:"limitranges"`
	Namespaces                 bool `json:"namespaces,omitempty" yaml:"namespaces"`
	Nodes                      bool `json:"nodes,omitempty" yaml:"nodes"`
	Persistentvolumeclaims     bool `json:"persistentvolumeclaims,omitempty" yaml:"persistentvolumeclaims"`
	Persistentvolumes          bool `json:"persistentvolumes,omitempty" yaml:"persistentvolumes"`
	Poddisruptionbudgets       bool `json:"poddisruptionbudgets,omitempty" yaml:"poddisruptionbudgets"`
	Pods                       bool `json:"pods,omitempty" yaml:"pods"`
	Replicasets                bool `json:"replicasets,omitempty" yaml:"replicasets"`
	Replicationcontrollers     bool `json:"replicationcontrollers,omitempty" yaml:"replicationcontrollers"`
	Resourcequotas             bool `json:"resourcequotas,omitempty" yaml:"resourcequotas"`
	Secrets                    bool `json:"secrets,omitempty" yaml:"secrets"`
	Services                   bool `json:"services,omitempty" yaml:"services"`
	Statefulsets               bool `json:"statefulsets,omitempty" yaml:"statefulsets"`
	Storageclasses             bool `json:"storageclasses,omitempty" yaml:"storageclasses"`
	Verticalpodautoscalers     bool `json:"verticalpodautoscalers,omitempty" yaml:"verticalpodautoscalers"`
}
