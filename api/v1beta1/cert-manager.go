package v1beta1

type CertManager struct {
	Deploy       bool   `json:"deploy,omitempty"`
	Prefix       string `json:"prefix,omitempty"`
	Namespace    string `json:"namespace,omitempty"`
	ReplicaCount int    `json:"replicaCount.omitempty"`
}
