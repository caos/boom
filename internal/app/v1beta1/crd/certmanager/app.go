package certmanager

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/v1beta1/crd/defaults"
	"github.com/caos/boom/internal/app/v1beta1/crd/service"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/template"
)

var (
	applicationName = "cert-manager"
)

type CertManager struct {
	ApplicationDirectoryPath string
	toolsDirectoryPath       string
	logger                   logging.Logger
	spec                     *toolsetsv1beta1.CertManager
}

func New(logger logging.Logger, toolsDirectoryPath string) *CertManager {
	c := &CertManager{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
		toolsDirectoryPath:       toolsDirectoryPath,
		logger:                   logger,
	}

	return c
}

func (c *CertManager) Reconcile(overlay string, specNamespace string, helm *template.Helm, spec *toolsetsv1beta1.CertManager) error {

	logFields := map[string]interface{}{
		"application": applicationName,
		"logID":       "CRD-S9vTqVD7gqyLPrQ",
	}

	c.logger.WithFields(logFields).Info("Reconciling")

	resultFilePath := defaults.GetResultFilePath(overlay, c.ApplicationDirectoryPath, applicationName)
	prefix := defaults.GetPrefix(overlay, applicationName, spec.Prefix)
	namespace := defaults.GetNamespace(overlay, applicationName, specNamespace, spec.Namespace)

	values := specToValues(helm.GetImageTags(applicationName), spec, namespace)
	writeValues := func(path string) error {
		if err := errors.Wrapf(helper.StructToYaml(values, path), "Failed to write values file overlay %s application %s", overlay, applicationName); err != nil {
			return err
		}
		return nil
	}

	if err := helm.PrepareTemplate(applicationName, prefix, namespace, writeValues); err != nil {
		return err
	}

	if spec.Deploy {
		if err := defaults.PrepareForResultOutput(defaults.GetResultFileDirectory(overlay, c.ApplicationDirectoryPath, applicationName)); err != nil {
			return err
		}

		if err := helm.Template(applicationName, resultFilePath); err != nil {
			return err
		}

		if err := addService(resultFilePath, prefix, namespace); err != nil {
			return err
		}

		if err := addCrds(resultFilePath, c.toolsDirectoryPath); err != nil {
			return err
		}

		if err := helper.DeleteKindFromYaml(resultFilePath, "Namespace"); err != nil {
			return err
		}

		kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath).AddFlag("--validate=false")
		if err := errors.Wrapf(helper.Run(c.logger, kubectlCmd.Build()), "Failed to apply with file %s", resultFilePath); err != nil {
			return err
		}

		c.spec = spec
	} else if !spec.Deploy && c.spec != nil {
		kubectlCmd := kubectl.New("delete").AddParameter("-f", resultFilePath)
		if err := errors.Wrapf(helper.Run(c.logger, kubectlCmd.Build()), "Failed to apply with file %s", resultFilePath); err != nil {
			return err
		}

		c.spec = nil
	}

	return nil
}

func addService(filePath string, prefix string, namespace string) error {
	name := "cert-manager"
	if prefix != "" {
		name = strings.Join([]string{prefix, name}, "-")
	}

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
	chartsPath := filepath.Join(toolsDirectoryPath, "charts", applicationName, "crds")

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

func specToValues(imageTags map[string]string, spec *toolsetsv1beta1.CertManager, namespace string) *Values {
	values := &Values{
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
				Namespace: namespace,
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

	if spec.ReplicaCount != 0 {
		values.ReplicaCount = spec.ReplicaCount
	}

	return values
}
