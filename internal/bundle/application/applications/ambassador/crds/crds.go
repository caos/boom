package crds

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

func GetCrdsFromSpec(spec *toolsetsv1beta1.Ambassador) []interface{} {
	ret := make([]interface{}, 0)

	if spec.Hosts.Argocd != nil {
		// privateKeySecret := strings.ReplaceAll(url.QueryEscape(spec.Hosts.Argocd.AcmeAuthority), "%", "-")
		host := GetHostFromConfig(&HostConfig{
			Name:             spec.Hosts.Argocd.Domain,
			Namespace:        "caos-system",
			InsecureAction:   "route",
			Hostname:         spec.Hosts.Argocd.Domain,
			AcmeProvider:     spec.Hosts.Argocd.AcmeAuthority,
			PrivateKeySecret: spec.Hosts.Argocd.Domain,
			Email:            spec.Hosts.Argocd.Email,
			TLSSecret:        spec.Hosts.Argocd.Domain,
		})
		ret = append(ret, host)
		mapping := GetMappingFromConfig(&MappingConfig{
			Name:      "argocd",
			Namespace: "caos-system",
			Prefix:    "/",
			Service:   "https://argocd-server.caos-system:443",
			Host:      spec.Hosts.Argocd.Domain,
		})
		ret = append(ret, mapping)
	}

	if spec.Hosts.Grafana != nil {
		// privateKeySecret := strings.ReplaceAll(url.QueryEscape(spec.Hosts.Argocd.AcmeAuthority), "%", "-")
		host := GetHostFromConfig(&HostConfig{
			Name:             spec.Hosts.Grafana.Domain,
			Namespace:        "caos-system",
			InsecureAction:   "route",
			Hostname:         spec.Hosts.Grafana.Domain,
			AcmeProvider:     spec.Hosts.Grafana.AcmeAuthority,
			PrivateKeySecret: spec.Hosts.Grafana.Domain,
			Email:            spec.Hosts.Grafana.Email,
			TLSSecret:        spec.Hosts.Grafana.Domain,
		})
		ret = append(ret, host)
		mapping := GetMappingFromConfig(&MappingConfig{
			Name:      "grafana",
			Namespace: "caos-system",
			Prefix:    "/",
			Service:   "grafana.caos-system",
			Host:      spec.Hosts.Grafana.Domain,
		})
		ret = append(ret, mapping)
	}

	return ret
}
