package v1beta1

type Logs struct {
	Ambassador             bool `json:"ambassador"`
	Argocd                 bool `json:"argocd"`
	KubeStateMetrics       bool `json:"kube-state-metrics"`
	PrometheusNodeExporter bool `json:"prometheus-node-exporter"`
	APIServer              bool `json:"api-server"`
	PrometheusOperator     bool `json:"prometheus-operator"`
	LoggingOperator        bool `json:"logging-operator"`
}
