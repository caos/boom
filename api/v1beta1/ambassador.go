package v1beta1

type Ambassador struct {
	Deploy       bool               `json:"deploy,omitempty"`
	ReplicaCount int                `json:"replicaCount,omitempty"`
	Service      *AmbassadorService `json:"service,omitempty"`
	Hosts        *Hosts             `json:"hosts,omitempty" yaml:"hosts,omitempty"`
}

type AmbassadorService struct {
	Type           string  `json:"type,omitempty" yaml:"type,omitempty"`
	LoadBalancerIP string  `json:"loadBalancerIP,omitempty" yaml:"loadBalancerIP,omitempty"`
	Ports          []*Port `json:"ports,omitempty" yaml:"ports,omitempty"`
}

type Port struct {
	Name       string `json:"name" yaml:"name"`
	Port       uint16 `json:"port,omitempty" yaml:"port,omitempty"`
	TargetPort uint16 `json:"targetPort,omitempty" yaml:"targetPort,omitempty"`
	NodePort   uint16 `json:"nodePort,omitempty" yaml:"nodePort,omitempty"`
}

type Hosts struct {
	Argocd  *Host `json:"argocd"`
	Grafana *Host `json:"grafana"`
}

type Host struct {
	Domain        string `json:"domain" yaml:"domain"`
	Email         string `json:"email" yaml:"email"`
	AcmeAuthority string `json:"acmeAuthority,omitempty" yaml:"acmeAuthority,omitempty"`
}
