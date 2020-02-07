package network

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/ambassador/crds"
)

func GetHostConfig(spec *toolsetsv1beta1.Network) *crds.HostConfig {
	return &crds.HostConfig{
		Name:             spec.Domain,
		Namespace:        "caos-system",
		InsecureAction:   "route",
		Hostname:         spec.Domain,
		AcmeProvider:     spec.AcmeAuthority,
		PrivateKeySecret: spec.Domain,
		Email:            spec.Email,
		TLSSecret:        spec.Domain,
	}
}

func GetMappingConfig(spec *toolsetsv1beta1.Network) *crds.MappingConfig {
	return &crds.MappingConfig{
		Name:      "argocd",
		Namespace: "caos-system",
		Prefix:    "/",
		Service:   "https://argocd-server.caos-system:443",
		Host:      spec.Domain,
	}
}
