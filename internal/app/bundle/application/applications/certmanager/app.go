package certmanager

import (
	"os"
	"path/filepath"

	"github.com/caos/orbiter/logging"

	"github.com/caos/boom/api/v1beta1"
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/bundle/application/resources/service"
	"github.com/caos/boom/internal/app/name"
	"github.com/caos/boom/internal/helper"
)

const (
	applicationName name.Application = "cert-manager"
)

func GetName() name.Application {
	return applicationName
}

type CertManager struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.CertManager
}

func New(logger logging.Logger) *CertManager {
	c := &CertManager{
		logger: logger,
	}

	return c
}

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.CertManager.Deploy
}

func (a *CertManager) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.CertManager != a.spec
}

func (a *CertManager) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	a.spec = toolsetCRDSpec.CertManager
}

func (c *CertManager) HelmPreApplySteps(resultFilePath string, spec *v1beta1.ToolsetSpec) {

	if err := addService(resultFilePath, "TODO"); err != nil {
		return
	}

	// if err := addCrds(resultFilePath, c.toolsDirectoryPath); err != nil {
	// 	return
	// }
}

func addService(filePath string, namespace string) error {
	name := "cert-manager"
	service := service.New(&service.Config{
		Name:       name,
		Namespace:  namespace,
		Labels:     map[string]string{"app": "cert-manager"},
		Protocol:   "TCP",
		Port:       9402,
		TargetPort: 9402,
		Selector:   map[string]string{"app": "cert-manager"},
	})

	if err := helper.AddStructToYaml(filePath, service); err != nil {
		return err
	}
	return nil
}

func addCrds(filePath string, toolsDirectoryPath string) error {
	chartsPath := filepath.Join(toolsDirectoryPath, "charts", applicationName.String(), "crds")

	var files []string
	err := filepath.Walk(chartsPath, func(path string, info os.FileInfo, err error) error {
		if path != chartsPath {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		err := helper.AddYamlToYaml(filePath, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CertManager) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolset.CertManager
	values := defaultValues(c.GetImageTags())

	if spec.ReplicaCount != 0 {
		values.ReplicaCount = spec.ReplicaCount
	}

	return values
}

func defaultValues(imageTags map[string]string) *Values {
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
			LeaderElection: &LeaderElection{
				Namespace: "TODO",
			},
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
