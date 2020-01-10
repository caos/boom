package prometheus

import "github.com/caos/boom/internal/app/bundle/application/applications/prometheus/servicemonitor"

type Config struct {
	Prefix                  string
	Namespace               string
	MonitorLabels           map[string]string
	ServiceMonitors         []*servicemonitor.Config
	ReplicaCount            int
	StorageSpec             *ConfigStorageSpec
	AdditionalScrapeConfigs []*AdditionalScrapeConfig
	KubeVersion             string
}

type ConfigStorageSpec struct {
	StorageClass string
	AccessModes  []string
	Storage      string
}
