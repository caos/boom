package helm

func DefaultValues(imageTags map[string]string) *Values {
	return &Values{
		FullNameOverride: "loki",
		Tracing:          &Tracing{},
		Config: &Config{
			AuthEnabled: false,
			Ingester: &Ingester{
				ChunkIdlePeriod:   "3m",
				ChunkBlockSize:    262144,
				ChunkRetainPeriod: "1m",
				Lifecycler: &Lifecycler{
					Ring: &Ring{
						Kvstore: &Kvstore{
							Store: "inmemory",
						},
						ReplicationFactor: 1,
					},
				},
			},
			LimitsConfig: &LimitsConfig{
				EnforceMetricName:      false,
				RejectOldSamples:       true,
				RejectOldSamplesMaxAge: "168h",
			},
			SchemaConfig: &SchemaConfigs{
				Configs: []*SchemaConfig{
					&SchemaConfig{
						From:        "2018-04-15",
						Store:       "boltdb",
						ObjectStore: "filesystem",
						Schema:      "v9",
						Index: &Index{
							Prefix: "index_",
							Period: "168h",
						},
					},
				},
			},
			Server: &Server{
				HTTPListenPort: 3100,
			},
			StorageConfig: &StorageConfig{
				Boltdb: &Boltdb{
					Directory: "/data/loki/index",
				},
				Filesystem: &Filesystem{
					Directory: "/data/loki/chunks",
				},
			},
			ChunkStoreConfig: &ChunkStoreConfig{
				MaxLookBackPeriod: 0,
			},
			TableManager: &TableManager{
				RetentionDeletesEnabled: false,
				RetentionPeriod:         0,
			},
		},
		Image: &Image{
			Repository: "grafana/loki",
			Tag:        imageTags["grafana/loki"],
			PullPolicy: "IfNotPresent",
		},
		LivenessProbe: &LivenessProbe{
			HTTPGet: &HTTPGet{
				Path: "/ready",
				Port: "http-metrics",
			},
			InitialDelaySeconds: 45,
		},
		NetworkPolicy: &NetworkPolicy{
			Enabled: false,
		},
		Persistence: &Persistence{
			Enabled:     false,
			AccessModes: []string{"ReadWriteOnce"},
			Size:        "10Gi",
		},
		PodAnnotations:      map[string]string{},
		PodManagementPolicy: "OrderedReady",
		Rbac: &Rbac{
			Create:     true,
			PspEnabled: true,
		},
		ReadinessProbe: &ReadinessProbe{
			HTTPGet: &HTTPGet{
				Path: "/ready",
				Port: "http-metrics",
			},
			InitialDelaySeconds: 45,
		},
		Replicas: 1,
		SecurityContext: &SecurityContext{
			FsGroup:      10001,
			RunAsGroup:   10001,
			RunAsNonRoot: true,
			RunAsUser:    10001,
		},
		Service: &Service{
			Type: "ClusterIP",
			Port: 3100,
		},
		ServiceAccount: &ServiceAccount{
			Create: true,
		},
		TerminationGracePeriodSeconds: 4800,
		UpdateStrategy: &UpdateStrategy{
			Type: "RollingUpdate",
		},
		ServiceMonitor: &ServiceMonitor{
			Enabled: false,
		},
	}
}
