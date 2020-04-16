package v1beta1

import "github.com/caos/boom/api/v1beta1/storage"

type Loki struct {
	Deploy        bool          `json:"deploy,omitempty"`
	Logs          *Logs         `json:"logs,omitempty"`
	Storage       *storage.Spec `json:"storage,omitempty"`
	ClusterOutput bool          `json:"clusterOutput,omitempty" yaml:"clusterOutput,omitempty"`
}

type Logs struct {
	Ambassador             bool `json:"ambassador"`
	Grafana                bool `json:"grafana"`
	Argocd                 bool `json:"argocd"`
	KubeStateMetrics       bool `json:"kube-state-metrics" yaml:"kube-state-metrics"`
	PrometheusNodeExporter bool `json:"prometheus-node-exporter"  yaml:"prometheus-node-exporter"`
	PrometheusOperator     bool `json:"prometheus-operator" yaml:"prometheus-operator"`
	LoggingOperator        bool `json:"logging-operator" yaml:"logging-operator"`
	Loki                   bool `json:"loki"`
	Prometheus             bool `json:"prometheus"`
}
