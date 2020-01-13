package kubestatemetrics

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/name"
)

const (
	applicationName name.Application = "kube-state-metrics"
)

func GetName() name.Application {
	return applicationName
}

type KubeStateMetrics struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.KubeStateMetrics
}

func New(logger logging.Logger) *KubeStateMetrics {
	lo := &KubeStateMetrics{
		logger: logger,
	}

	return lo
}
func (k *KubeStateMetrics) GetName() name.Application {
	return applicationName
}

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.KubeStateMetrics.Deploy
}

func (k *KubeStateMetrics) Initial() bool {
	return k.spec == nil
}

func (k *KubeStateMetrics) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.KubeStateMetrics, k.spec)
}

func (k *KubeStateMetrics) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	k.spec = toolsetCRDSpec.KubeStateMetrics
}

func (k *KubeStateMetrics) GetNamespace() string {
	return "caos-system"
}

func (k *KubeStateMetrics) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.CertManager
	values := defaultValues(k.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	return values
}

func defaultValues(imageTags map[string]string) *Values {
	return &Values{
		FullnameOverride: "kube-state-metrics",
		PrometheusScrape: true,
		Image: &Image{
			Repository: "quay.io/coreos/kube-state-metrics",
			Tag:        imageTags["quay.io/coreos/kube-state-metrics"],
			PullPolicy: "IfNotPresent",
		},
		Replicas: 1,
		Service: &Service{
			Port:           8080,
			Type:           "ClusterIP",
			NodePort:       0,
			LoadBalancerIP: "",
			Annotations:    map[string]string{},
		},
		CustomLabels: map[string]string{},
		HostNetwork:  false,
		Rbac: &Rbac{
			Create: true,
		},
		ServiceAccount: &ServiceAccount{
			Create: true,
			Name:   "",
		},
		Prometheus: &Prometheus{
			Monitor: &Monitor{
				Enabled: false,
			},
		},
		PodSecurityPolicy: &PodSecurityPolicy{
			Enabled: false,
		},
		SecurityContext: &SecurityContext{
			Enabled:   true,
			RunAsUser: 65534,
			FsGroup:   65534,
		},
		NodeSelector:   map[string]string{},
		Affinity:       nil,
		Tolerations:    nil,
		PodAnnotations: map[string]string{},
		Collectors: &Collectors{
			Certificatesigningrequests: true,
			Configmaps:                 true,
			Cronjobs:                   true,
			Daemonsets:                 true,
			Deployments:                true,
			Endpoints:                  true,
			Horizontalpodautoscalers:   true,
			Ingresses:                  true,
			Jobs:                       true,
			Limitranges:                true,
			Namespaces:                 true,
			Nodes:                      true,
			Persistentvolumeclaims:     true,
			Persistentvolumes:          true,
			Poddisruptionbudgets:       true,
			Pods:                       true,
			Replicasets:                true,
			Replicationcontrollers:     true,
			Resourcequotas:             true,
			Secrets:                    true,
			Services:                   true,
			Statefulsets:               true,
			Storageclasses:             true,
			Verticalpodautoscalers:     false,
		},
	}
}
