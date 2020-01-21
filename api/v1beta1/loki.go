package v1beta1

type Loki struct {
	Deploy  bool         `json:"deploy,omitempty"`
	Logs    *Logs        `json:"logs,omitempty"`
	Storage *StorageSpec `json:"storage,omitempty"`
}

type Logs struct {
	Ambassador             bool `json:"ambassador"`
	Grafana                bool `json:"grafana"`
	Argocd                 bool `json:"argocd"`
	KubeStateMetrics       bool `json:"kube-state-metrics"`
	PrometheusNodeExporter bool `json:"prometheus-node-exporter"`
	PrometheusOperator     bool `json:"prometheus-operator"`
	LoggingOperator        bool `json:"logging-operator"`
	Loki                   bool `json:"loki"`
}
