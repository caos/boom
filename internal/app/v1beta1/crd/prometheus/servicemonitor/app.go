package servicemonitor

type ConfigEndpoint struct {
	Port       string
	TargetPort string
	Interval   string
	Scheme     string
	Path       string
}

type Config struct {
	Name                  string
	Endpoints             []*ConfigEndpoint
	MonitorMatchingLabels map[string]string
	ServiceMatchingLabels map[string]string
}

func SpecToValues(config *Config) *Values {

	endpoints := make([]*Endpoint, 0)
	for _, endpoint := range config.Endpoints {
		valueEndpoint := &Endpoint{
			Port:       endpoint.Port,
			TargetPort: endpoint.TargetPort,
			Interval:   endpoint.Interval,
			Scheme:     endpoint.Scheme,
			Path:       endpoint.Path,
		}
		endpoints = append(endpoints, valueEndpoint)
	}

	values := &Values{
		Name:             config.Name,
		AdditionalLabels: config.MonitorMatchingLabels,
		Selector: &Selector{
			MatchLabels: config.ServiceMatchingLabels,
		},
		NamespaceSelector: &NamespaceSelector{
			Any: true,
		},
		Endpoints: endpoints,
	}

	return values
}
