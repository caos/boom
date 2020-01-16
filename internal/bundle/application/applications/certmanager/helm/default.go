package helm

func DefaultValues(imageTags map[string]string) *Values {
	return &Values{
		FullnameOverride: "cert-manager",
		Global: &Global{
			IsOpenshift: false,
			Rbac: &Rbac{
				Create: true,
			},
			PodSecurityPolicy: &PodSecurityPolicy{
				Enabled: false,
			},
			LogLevel: 2,
		},
		ReplicaCount: 1,
		Image: &Image{
			Repository: "quay.io/jetstack/cert-manager-controller",
			Tag:        imageTags["quay.io/jetstack/cert-manager-controller"],
			PullPolicy: "IfNotPresent",
		},
		ServiceAccount: &ServiceAccount{
			Create: true,
		},
		SecurityContext: &SecurityContext{
			Enabled:   false,
			FsGroup:   1001,
			RunAsUser: 1001,
		},
		Prometheus: &Prometheus{
			Enabled: false,
		},
		Webhook: &Webhook{
			Enabled:      false,
			ReplicaCount: 1,
			Image: &Image{
				Repository: "quay.io/jetstack/cert-manager-webhook",
				Tag:        imageTags["quay.io/jetstack/cert-manager-webhook"],
				PullPolicy: "IfNotPresent",
			},
			InjectAPIServerCA: true,
			SecurePort:        10250,
		},
		Cainjector: &Cainjector{
			ReplicaCount: 1,
			Image: &Image{
				Repository: "quay.io/jetstack/cert-manager-cainjector",
				Tag:        imageTags["quay.io/jetstack/cert-manager-cainjector"],
				PullPolicy: "IfNotPresent",
			},
		},
	}
}
