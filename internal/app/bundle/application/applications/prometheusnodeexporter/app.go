package prometheusnodeexporter

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/name"
)

const (
	applicationName name.Application = "prometheus-node-exporter"
)

func GetName() name.Application {
	return applicationName
}

type PrometheusNodeExporter struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.PrometheusNodeExporter
}

func New(logger logging.Logger) *PrometheusNodeExporter {
	pne := &PrometheusNodeExporter{
		logger: logger,
	}

	return pne
}
func (pne *PrometheusNodeExporter) GetName() name.Application {
	return applicationName
}

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.PrometheusNodeExporter.Deploy
}

func (pne *PrometheusNodeExporter) Initial() bool {
	return pne.spec == nil
}

func (pne *PrometheusNodeExporter) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.PrometheusNodeExporter, pne.spec)
}

func (pne *PrometheusNodeExporter) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	pne.spec = toolsetCRDSpec.PrometheusNodeExporter
}

func (pne *PrometheusNodeExporter) GetNamespace() string {
	return "caos-system"
}

func (p *PrometheusNodeExporter) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.PrometheusNodeExporter
	values := defaultValues(p.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	return values
}

func defaultValues(imageTags map[string]string) *Values {
	return &Values{
		FullnameOverride: "node-exporter",
		Image: &Image{
			Repository: "quay.io/prometheus/node-exporter",
			Tag:        imageTags["quay.io/prometheus/node-exporter"],
			PullPolicy: "IfNotPresent",
		},
		Service: &Service{
			Type:        "ClusterIP",
			Port:        9100,
			TargetPort:  9100,
			NodePort:    "",
			Annotations: map[string]string{"prometheus.io/scrape": "true"},
		},
		Prometheus: &Prometheus{
			Monitor: &Monitor{
				Enabled:          false,
				AdditionalLabels: map[string]string{},
				Namespace:        "",
				ScrapeTimeout:    "10s",
			},
		},
		ServiceAccount: &ServiceAccount{
			Create: true,
		},
		SecurityContext: &SecurityContext{
			RunAsNonRoot: true,
			RunAsUser:    65534,
		},
		Rbac: &Rbac{
			Create:     true,
			PspEnabled: true,
		},
		HostNetwork: false,
		Tolerations: []*Toleration{&Toleration{
			Effect:   "NoSchedule",
			Operator: "Exists",
		}},
	}
}
