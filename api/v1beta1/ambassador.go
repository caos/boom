package v1beta1

type Ambassador struct {
	Deploy         bool   `json:"deploy,omitempty"`
	ReplicaCount   int    `json:"replicaCount,omitempty"`
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`
}
