package prometheus

import (
	"errors"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/config"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/helm"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/info"
	"github.com/caos/boom/internal/bundle/application/applications/prometheus/servicemonitor"
	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/mntr"
)

func (p *Prometheus) SpecToHelmValues(monitor mntr.Monitor, toolsetCRDSpec *v1beta1.ToolsetSpec) interface{} {
	version, err := kubectl.NewVersion().GetKubeVersion(monitor)
	if err != nil {
		// TODO: Better error handling?
		return nil
	}

	_, err = clientgo.GetSecret("grafana-cloud", "caos-system")
	if err != nil && !errors.Is(err, clientgo.ErrNotFound{}) {
		// TODO: Better error handling?
		panic(err)
	}
	ingestionSecretPresent := err == nil

	config := config.ScrapeMetricsCrdsConfig(info.GetInstanceName(), toolsetCRDSpec)

	values := helm.DefaultValues(p.GetImageTags())

	if config.StorageSpec != nil {
		storageSpec := &helm.StorageSpec{
			VolumeClaimTemplate: &helm.VolumeClaimTemplate{
				Spec: &helm.VolumeClaimTemplateSpec{
					StorageClassName: config.StorageSpec.StorageClass,
					AccessModes:      config.StorageSpec.AccessModes,
					Resources: &helm.Resources{
						Requests: &helm.Request{
							Storage: config.StorageSpec.Storage,
						},
					},
				},
			},
		}

		values.Prometheus.PrometheusSpec.StorageSpec = storageSpec
	}

	if config.MonitorLabels != nil {
		values.Prometheus.PrometheusSpec.ServiceMonitorSelector = &helm.MonitorSelector{
			MatchLabels: config.MonitorLabels,
		}
	}

	if config.ServiceMonitors != nil {
		additionalServiceMonitors := make([]*servicemonitor.Values, 0)
		for _, specServiceMonitor := range config.ServiceMonitors {
			valuesServiceMonitor := servicemonitor.SpecToValues(specServiceMonitor)
			additionalServiceMonitors = append(additionalServiceMonitors, valuesServiceMonitor)
		}

		values.Prometheus.AdditionalServiceMonitors = additionalServiceMonitors
	}

	if config.AdditionalScrapeConfigs != nil {
		values.Prometheus.PrometheusSpec.AdditionalScrapeConfigs = config.AdditionalScrapeConfigs
	}

	if ingestionSecretPresent {
		if values.Prometheus.PrometheusSpec.ExternalLabels == nil {
			values.Prometheus.PrometheusSpec.ExternalLabels = make(map[string]string)
		}
		values.Prometheus.PrometheusSpec.ExternalLabels["orb"] = p.orb
		values.Prometheus.PrometheusSpec.RemoteWrite = append(values.Prometheus.PrometheusSpec.RemoteWrite, &helm.RemoteWrite{
			URL: "https://prometheus-us-central1.grafana.net/api/prom/push",
			BasicAuth: &helm.BasicAuth{
				Username: &helm.SecretKeySelector{
					Name: "grafana-cloud",
					Key:  "username",
				},
				Password: &helm.SecretKeySelector{
					Name: "grafana-cloud",
					Key:  "password",
				},
			},
			WriteRelabelConfigs: []*helm.RelabelConfig{{
				Action: "keep",
				SourceLabels: []string{
					"__name__",
					"job",
				},
				Regex: "caos_.+;.*|up;caos_remote_.+",
			}},
		})
	}

	if toolsetCRDSpec.Prometheus.RemoteWrite != nil {
		values.Prometheus.PrometheusSpec.RemoteWrite = append(values.Prometheus.PrometheusSpec.RemoteWrite, &helm.RemoteWrite{
			URL: toolsetCRDSpec.Prometheus.RemoteWrite.URL,
			BasicAuth: &helm.BasicAuth{
				Username: &helm.SecretKeySelector{
					Name: toolsetCRDSpec.Prometheus.RemoteWrite.BasicAuth.Username.Name,
					Key:  toolsetCRDSpec.Prometheus.RemoteWrite.BasicAuth.Username.Key,
				},
				Password: &helm.SecretKeySelector{
					Name: toolsetCRDSpec.Prometheus.RemoteWrite.BasicAuth.Password.Name,
					Key:  toolsetCRDSpec.Prometheus.RemoteWrite.BasicAuth.Password.Key,
				},
			},
		})
	}

	ruleLabels := labels.GetRuleLabels(info.GetInstanceName())
	rules, _ := helm.GetDefaultRules(ruleLabels)

	values.Prometheus.PrometheusSpec.RuleSelector = &helm.RuleSelector{MatchLabels: ruleLabels}
	values.DefaultRules.Labels = ruleLabels
	values.KubeTargetVersionOverride = version
	values.AdditionalPrometheusRules = []*helm.AdditionalPrometheusRules{rules}

	values.FullnameOverride = info.GetInstanceName()

	return values
}

func (p *Prometheus) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (p *Prometheus) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
