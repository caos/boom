package grafana

import (
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/caos/orbiter/logging"

	"github.com/caos/boom/api/v1beta1"
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/bundle/application/applications/grafanastandalone"
	"github.com/caos/boom/internal/app/bundle/application/applications/prometheusoperator"
	"github.com/caos/boom/internal/app/name"
)

const (
	applicationName name.Application = "grafana"
)

func GetName() name.Application {
	return applicationName
}

type Grafana struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.Grafana
}

func New(logger logging.Logger) *Grafana {
	return &Grafana{
		logger: logger,
	}
}
func (g *Grafana) GetName() name.Application {
	return applicationName
}

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.Grafana.Deploy
}

func (g *Grafana) Initial() bool {
	return g.spec == nil
}

func (g *Grafana) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.Grafana, g.spec)
}

func (g *Grafana) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	g.spec = toolsetCRDSpec.Grafana
}

func (g *Grafana) GetNamespace() string {
	return "caos-system"
}

func (g *Grafana) HelmPreApplySteps(spec *v1beta1.ToolsetSpec) ([]interface{}, error) {
	config := newConfig(spec.KubeVersion, spec)

	folders := make([]string, 0)
	for _, provider := range config.DashboardProviders {
		folders = append(folders, provider.Folder)
	}

	outs, err := getKustomizeOutput(folders)
	if err != nil {
		return nil, err
	}

	ret := make([]interface{}, len(outs))
	for k, v := range outs {
		ret[k] = v
	}
	return ret, nil
}

func (g *Grafana) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {

	config := newConfig(toolset.KubeVersion, toolset)
	values := defaultValues(g.GetImageTags())

	values.KubeTargetVersionOverride = config.KubeVersion

	if config.Datasources != nil {
		datasources := make([]*grafanastandalone.Datasource, 0)
		for _, datasource := range config.Datasources {
			valuesDatasource := &grafanastandalone.Datasource{
				Name:      datasource.Name,
				Type:      datasource.Type,
				URL:       datasource.Url,
				Access:    datasource.Access,
				IsDefault: datasource.IsDefault,
			}
			datasources = append(datasources, valuesDatasource)
		}
		values.Grafana.AdditionalDataSources = datasources
	}

	if config.DashboardProviders != nil {
		providers := make([]*Provider, 0)
		dashboards := make(map[string]string, 0)
		for _, provider := range config.DashboardProviders {
			for _, configmap := range provider.ConfigMaps {
				providers = append(providers, getProvider(configmap))
				dashboards[configmap] = configmap
			}
		}
		values.Grafana.DashboardProviders = &DashboardProviders{
			Providers: &Providersyaml{
				APIVersion: 1,
				Providers:  providers,
			},
		}
		values.Grafana.DashboardsConfigMaps = dashboards
	}

	return values
}

func defaultValues(imageTags map[string]string) *Values {
	grafana := &GrafanaValues{
		FullnameOverride:         "grafana",
		Enabled:                  true,
		DefaultDashboardsEnabled: true,
		AdminPassword:            "admin",
		Ingress: &Ingress{
			Enabled: false,
		},
		Sidecar: &Sidecar{
			Dashboards: &Dashboards{
				Enabled: true,
				Label:   "grafana_dashboard",
			},
			Datasources: &Datasources{
				Enabled: true,
				Label:   "grafana_datasource",
			},
		},
		ServiceMonitor: &ServiceMonitor{
			SelfMonitor: false,
		},
	}

	return &Values{
		DefaultRules: &DefaultRules{
			Create: false,
			Rules: &Rules{
				Alertmanager:                false,
				Etcd:                        false,
				General:                     false,
				K8S:                         false,
				KubeApiserver:               false,
				KubePrometheusNodeAlerting:  false,
				KubePrometheusNodeRecording: false,
				KubernetesAbsent:            false,
				KubernetesApps:              false,
				KubernetesResources:         false,
				KubernetesStorage:           false,
				KubernetesSystem:            false,
				KubeScheduler:               false,
				Network:                     false,
				Node:                        false,
				Prometheus:                  false,
				PrometheusOperator:          false,
				Time:                        false,
			},
		},
		Global: &Global{
			Rbac: &Rbac{
				Create:     false,
				PspEnabled: false,
			},
		},
		FullnameOverride: "grafana",
		Alertmanager: &DisabledTool{
			Enabled: false,
		},
		Grafana: grafana,
		KubeAPIServer: &DisabledTool{
			Enabled: false,
		},
		Kubelet: &DisabledTool{
			Enabled: false,
		},
		KubeControllerManager: &DisabledTool{
			Enabled: false,
		},
		CoreDNS: &DisabledTool{
			Enabled: false,
		},
		KubeDNS: &DisabledTool{
			Enabled: false,
		},
		KubeEtcd: &DisabledTool{
			Enabled: false,
		},
		KubeScheduler: &DisabledTool{
			Enabled: false,
		},
		KubeProxy: &DisabledTool{
			Enabled: false,
		},
		KubeStateMetricsScrap: &DisabledTool{
			Enabled: false,
		},
		KubeStateMetrics: &DisabledTool{
			Enabled: false,
		},
		NodeExporter: &DisabledTool{
			Enabled: false,
		},
		PrometheusNodeExporter: &DisabledTool{
			Enabled: false,
		},
		PrometheusOperator: &prometheusoperator.PrometheusOperatorValues{
			Enabled: false,
			TLSProxy: &prometheusoperator.TLSProxy{
				Enabled: false,
				Image: &prometheusoperator.Image{
					Repository: "squareup/ghostunnel",
					Tag:        imageTags["squareup/ghostunnel"],
					PullPolicy: "IfNotPresent",
				},
			},
			AdmissionWebhooks: &prometheusoperator.AdmissionWebhooks{
				FailurePolicy: "Fail",
				Enabled:       false,
				Patch: &prometheusoperator.Patch{
					Enabled: false,
					Image: &prometheusoperator.Image{
						Repository: "jettech/kube-webhook-certgen",
						Tag:        imageTags["jettech/kube-webhook-certgen"],
						PullPolicy: "IfNotPresent",
					},
					PriorityClassName: "",
				},
			},
			ServiceAccount: &prometheusoperator.ServiceAccount{
				Create: false,
			},
			ServiceMonitor: &prometheusoperator.ServiceMonitor{
				Interval:    "",
				SelfMonitor: false,
			},
			CreateCustomResource: true,
			KubeletService: &prometheusoperator.KubeletService{
				Enabled: false,
			},
		},
		Prometheus: &DisabledTool{
			Enabled: false,
		},
	}
}

func getKustomizeOutput(folders []string) ([]string, error) {
	ret := make([]string, len(folders))
	for n, folder := range folders {
		command := strings.Join([]string{"kustomize build", folder}, " ")

		cmd := exec.Command("/bin/sh", "-c", command)
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		ret[n] = string(out)
	}
	return ret, nil
}

func getProvider(appName string) *Provider {
	return &Provider{
		Name:            appName,
		Type:            "file",
		DisableDeletion: false,
		Editable:        true,
		Options: map[string]string{
			"path": filepath.Join("/var/lib/grafana/dashboards", appName),
		},
	}
}
