package v1beta1

type Prometheus struct {
	Deploy  bool         `json:"deploy,omitempty"`
	Metrics *Metrics     `json:"metrics,omitempty"`
	Storage *StorageSpec `json:"storage,omitempty"`
}

type StorageSpec struct {
	StorageClass string   `json:"storageClass,omitempty"`
	AccessModes  []string `json:"accessModes,omitempty"`
	Size         string   `json:"size,omitempty"`
}

type Metrics struct {
	Ambassador             bool `json:"ambassador"`
	Argocd                 bool `json:"argocd"`
	KubeStateMetrics       bool `json:"kube-state-metrics"`
	PrometheusNodeExporter bool `json:"prometheus-node-exporter"`
	APIServer              bool `json:"api-server"`
	PrometheusOperator     bool `json:"prometheus-operator"`
	LoggingOperator        bool `json:"logging-operator"`
	Loki                   bool `json:"loki"`
}
