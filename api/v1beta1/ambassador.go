package v1beta1

type Ambassador struct {
	Deploy         bool   `json:"deploy,omitempty"`
	Prefix         string `json:"prefix,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`
	ScrapeMetrics  bool   `json:"scrapeMetrics,omitempty" yaml:"scrapeMetrics,omitempty"`
}
