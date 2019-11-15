package grafana

import (
	"os"
	"path/filepath"
	"strings"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/kubectl"
	"github.com/caos/toolsop/internal/template"
	"github.com/caos/utils/logging"
)

var (
	applicationName      = "grafana"
	resultsDirectoryName = "results"
	resultsFileName      = "results.yaml"
	defaultNamespace     = "monitoring"
)

type Grafana struct {
	ApplicationDirectoryPath string
}

func New(toolsDirectoryPath string) *Grafana {
	lo := &Grafana{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
	}

	return lo
}

func (g *Grafana) Reconcile(overlay string, helm *template.Helm, spec *toolsetsv1beta1.Grafana) error {

	logging.Log("CRD-tS3NCOfewXYGvDE").Infof("Reconciling application %s", applicationName)
	resultsFileDirectory := filepath.Join(g.ApplicationDirectoryPath, resultsDirectoryName, overlay)
	_ = os.RemoveAll(resultsFileDirectory)
	_ = os.MkdirAll(resultsFileDirectory, os.ModePerm)
	resultFilePath := filepath.Join(resultsFileDirectory, resultsFileName)

	values := specToValues(helm.GetImageTags(applicationName), spec)
	writeValues := func(path string) error {
		if err := helper.StructToYaml(values, path); err != nil {
			logging.Log("CRD-ZXIlvuoOW1WWBpU").Debugf("Failed to write values file overlay %s application %s", overlay, applicationName)
			return err
		}
		return nil
	}

	prefix := spec.Prefix
	if prefix == "" {
		prefix = overlay
	}
	namespace := spec.Namespace
	if namespace == "" {
		namespace = strings.Join([]string{overlay, defaultNamespace}, "-")
	}

	if err := helm.Template(applicationName, prefix, namespace, resultFilePath, writeValues); err != nil {
		return err
	}

	kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)

	if spec.Deploy {
		if err := kubectlCmd.Run(); err != nil {
			logging.Log("CRD-HcT1sFDBfJMCQHG").OnError(err).Debugf("Failed to apply file %s", resultFilePath)
			return err
		}
	}

	return nil
}

func specToValues(imageTags map[string]string, spec *toolsetsv1beta1.Grafana) *Values {
	values := &Values{
		Rbac: &Rbac{
			Create:         true,
			PspEnabled:     true,
			PspUseAppArmor: true,
			Namespaced:     false,
		},
		ServiceAccount: &ServiceAccount{
			Create: true,
		},
		Replicas: 1,
		DeploymentStrategy: &DeploymentStrategy{
			Type: "RollingUpdate",
		},
		ReadinessProbe: &ReadinessProbe{
			HTTPGet: &HTTPGet{
				Port: 3000,
				Path: "/api/health",
			},
		},
		LivenessProbe: &LivenessProbe{
			HTTPGet: &HTTPGet{
				Port: 3000,
				Path: "/api/health",
			},
			InitialDelaySeconds: 60,
			TimeoutSeconds:      30,
			FailureThreshold:    10,
		},
		Image: &Image{
			Repository: "grafana/grafana",
			Tag:        imageTags["grafana/grafana"],
			PullPolicy: "IfNotPresent",
		},
		TestFramework: &TestFramework{
			Enabled: true,
			Image:   "dduportal/bats",
			Tag:     imageTags["dduportal/bats"],
		},
		SecurityContext: &SecurityContext{
			RunAsUser: 472,
			FsGroup:   472,
		},
		DownloadDashboardsImage: &DownloadDashboardsImage{
			Repository: "appropriate/curl",
			Tag:        imageTags["appropriate/curl"],
			PullPolicy: "IfNotPresent",
		},
		DownloadDashboards: &DownloadDashboards{},
		PodPortName:        "grafana",
		Service: &Service{
			Type:       "ClusterIP",
			Port:       80,
			TargetPort: 3000,
			PortName:   "service",
		},
		Ingress: &Ingress{
			Enabled: false,
		},
		Persistence: &Persistence{
			Type:        "pvc",
			Enabled:     false,
			AccessModes: []string{"ReadWriteOnce"},
			Size:        "10Gi",
			Finalizers:  []string{"kubernetes.io/pvc-protection"},
		},
		InitChownData: &InitChownData{
			Enabled: true,
			Image: &Image{
				Repository: "busybox",
				Tag:        imageTags["busybox"],
				PullPolicy: "IfNotPresent",
			},
		},
		AdminUser:     "admin",
		AdminPassword: "admin",
		Admin: &Admin{
			ExistingSecret: "",
			UserKey:        "admin-user",
			PasswordKey:    "admin-password",
		},
		// Datasources             *Datasources             `yaml:"datasources"`
		// Dashboards              *Dashboards              `yaml:"dashboards"`
		// DashboardsConfigMaps    map[string]string        `yaml:"dashboardsConfigMaps"`
		GrafanaIni: &GrafanaIni{
			Paths: &Paths{
				Data:         "/var/lib/grafana/data",
				Logs:         "/var/log/grafana",
				Plugins:      "/var/lib/grafana/plugins",
				Provisioning: "/etc/grafana/provisioning",
			},
			Analytics: &Analytics{
				CheckForUpdates: true,
			},
			Log: &Log{
				Mode: "console",
			},
			GrafanaNet: &GrafanaNet{
				URL: "https://grafana.net",
			},
		},
		Ldap: &Ldap{
			Enabled: false,
		},
		SMTP: &SMTP{
			ExistingSecret: "",
			UserKey:        "user",
			PasswordKey:    "password",
		},
		Sidecar: &Sidecar{
			Image:           "kiwigrid/k8s-sidecar:0.1.20",
			ImagePullPolicy: "IfNotPresent",
			Dashboards: &DashboardsSidecar{
				Enabled: false,
			},
			Datasources: &DatasourcesSidecar{
				Enabled: false,
			},
		},
	}

	if spec.Admin != nil {
		values.Admin.ExistingSecret = spec.Admin.ExistingSecret
		values.Admin.UserKey = spec.Admin.UserKey
		values.Admin.PasswordKey = spec.Admin.PasswordKey
	}

	if spec.Datasources != nil {
		datasources := make([]*Datasource, 0)
		for _, datasource := range spec.Datasources {
			valuesDatasource := &Datasource{
				Name:      datasource.Name,
				Type:      datasource.Type,
				URL:       datasource.Url,
				Access:    datasource.Access,
				IsDefault: datasource.IsDefault,
			}
			datasources = append(datasources, valuesDatasource)
		}
		values.Datasources = &Datasources{
			Datasources: &Datasourcesyaml{
				APIVersion:  1,
				Datasources: datasources,
			},
		}
	}

	if spec.Dashboards != nil {
		for _, dConfigMap := range spec.Dashboards {
			values.DashboardsConfigMaps[dConfigMap.ConfigMap] = dConfigMap.ConfigMap

			values.Dashboards.Dashboards = make(map[string]map[string]*DashboardFile, 0)
			for _, dashboard := range dConfigMap.FileNames {
				filePath := filepath.Join("dashboards", dashboard.FileName)
				values.Dashboards.Dashboards[dConfigMap.ConfigMap][dashboard.Name] = &DashboardFile{
					File: filePath,
				}
			}
		}
	}

	return values
}
