package v1beta1

type KubeStateMetrics struct {
	Deploy       bool `json:"deploy,omitempty"`
	ReplicaCount int  `json:"replicaCount,omitempty"`
}
