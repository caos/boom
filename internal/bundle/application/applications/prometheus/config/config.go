package config

import (
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/helm"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
)

type Config struct {
	Prefix                  string
	Namespace               string
	MonitorLabels           map[string]string
	ServiceMonitors         []*servicemonitor.Config
	ReplicaCount            int
	StorageSpec             *StorageSpec
	AdditionalScrapeConfigs []*helm.AdditionalScrapeConfig
	KubeVersion             string
}

type StorageSpec struct {
	StorageClass string
	AccessModes  []string
	Storage      string
}
