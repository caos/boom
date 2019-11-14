package servicemonitor

import toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"

func SpecToValues(labels map[string]string, spec *toolsetsv1beta1.ServiceMonitor) *Values {

	endpoints := make([]*Endpoint, 0)
	for _, endpoint := range spec.Endpoints {
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
		Name:             spec.Name,
		AdditionalLabels: labels,
		Selector: &Selector{
			MatchLabels: spec.ServiceMatchingLabels,
		},
		NamespaceSelector: &NamespaceSelector{
			Any: true,
		},
		Endpoints: endpoints,
	}

	return values
}
