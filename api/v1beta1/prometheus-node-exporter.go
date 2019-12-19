package v1beta1

type PrometheusNodeExporter struct {
	Deploy    bool   `json:"deploy,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}
