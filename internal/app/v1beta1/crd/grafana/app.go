package grafana

import (
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	"github.com/caos/boom/internal/app/v1beta1/crd/defaults"
	"github.com/caos/boom/internal/app/v1beta1/crd/grafanastandalone"
	"github.com/caos/boom/internal/app/v1beta1/crd/prometheusoperator"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/template"
)

var (
	applicationName = "grafana"
)

type Grafana struct {
	ApplicationDirectoryPath string
	logger                   logging.Logger
	config                   *Config
}

func New(logger logging.Logger, toolsDirectoryPath string) *Grafana {
	lo := &Grafana{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
		logger:                   logger,
	}

	return lo
}

func (g *Grafana) Reconcile(overlay, specNamespace string, helm *template.Helm, config *Config) error {

	logFields := map[string]interface{}{
		"application": applicationName,
		"logID":       "CRD-ehpIBYPdrZPAynI",
	}
	g.logger.WithFields(logFields).Info("Reconciling")

	resultFilePath := defaults.GetResultFilePath(overlay, g.ApplicationDirectoryPath, applicationName)
	prefix := defaults.GetPrefix(overlay, applicationName, config.Prefix)
	namespace := defaults.GetNamespace(overlay, applicationName, specNamespace, config.Namespace)

	values := specToValues(helm.GetImageTags(applicationName), config)
	writeValues := func(path string) error {
		if err := errors.Wrapf(helper.StructToYaml(values, path), "Failed to write values file overlay %s application %s", overlay, applicationName); err != nil {
			return err
		}
		return nil
	}

	if err := helm.PrepareTemplate(applicationName, prefix, namespace, writeValues); err != nil {
		return err
	}

	if config.Deploy && !reflect.DeepEqual(g.config, config) {
		if err := defaults.PrepareForResultOutput(defaults.GetResultFileDirectory(overlay, g.ApplicationDirectoryPath, applicationName)); err != nil {
			return err
		}

		if err := helm.Template(applicationName, resultFilePath); err != nil {
			return err
		}

		if err := helper.DeleteKindFromYaml(resultFilePath, "Namespace"); err != nil {
			return err
		}

		folders := make([]string, 0)
		for _, provider := range config.DashboardProviders {
			folders = append(folders, provider.Folder)
		}

		if err := applyKustomize(folders); err != nil {
			return err
		}

		kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)
		if err := errors.Wrapf(helper.Run(g.logger, kubectlCmd.Build()), "Failed to apply file %s", resultFilePath); err != nil {
			return err
		}

		g.config = config
	} else if !config.Deploy && g.config != nil {
		kubectlCmd := kubectl.New("delete").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)
		if err := errors.Wrapf(helper.Run(g.logger, kubectlCmd.Build()), "Failed to apply file %s", resultFilePath); err != nil {
			return err
		}

		g.config = nil
	}
	return nil
}

func specToValues(imageTags map[string]string, config *Config) *Values {
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

	values := &Values{
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
		FullnameOverride:          "grafana",
		KubeTargetVersionOverride: config.KubeVersion,
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

func applyKustomize(folders []string) error {
	for _, folder := range folders {
		command := strings.Join([]string{"kustomize build", folder, "| kubectl apply -f -"}, " ")

		cmd := exec.Command("/bin/sh", "-c", command)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
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
