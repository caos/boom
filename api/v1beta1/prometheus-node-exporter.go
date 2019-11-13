package v1beta1

type PrometheusNodeExporter struct {
	Deploy    bool     `json:"deploy,omitempty"`
	Prefix    string   `json:"prefix,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
	Monitor   *Monitor `json:"monitor,omitempty"`
}

type Monitor struct {
	Enabled   bool   `json:"enabled"`
	Namespace string `json:"namespace"`
	// AdditionalLabels map[string]string `json:"additionalLabels"`
	// ScrapeTimeout    string            `json:"scrapeTimeout"`
}
