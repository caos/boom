package certmanager

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/kubectl"
	"github.com/caos/toolsop/internal/template"
)

var (
	applicationName  = "cert-manager"
	resultsFilename  = "results.yaml"
	defaultNamespace = "kube-system"
)

type CertManager struct {
	ApplicationDirectoryPath string
}

func New(toolsDirectoryPath string) *CertManager {
	c := &CertManager{
		ApplicationDirectoryPath: strings.Join([]string{toolsDirectoryPath, applicationName}, "/"),
	}

	return c
}

func (c *CertManager) Reconcile(overlay string, helm *template.Helm, spec *toolsetsv1beta1.CertManager) error {

	resultFilePath := strings.Join([]string{c.ApplicationDirectoryPath, resultsFilename}, "/")

	prefix := spec.Prefix
	if prefix == "" {
		prefix = overlay
	}
	namespace := spec.Namespace
	if namespace == "" {
		namespace = strings.Join([]string{overlay, defaultNamespace}, "-")
	}

	values := specToValues(helm.GetImageTags(applicationName), spec, namespace)
	writeValues := func(path string) error {
		if err := helper.StructToYaml(values, path); err != nil {
			return err
		}
		return nil
	}

	if err := helm.Template(applicationName, prefix, namespace, resultFilePath, writeValues); err != nil {
		return err
	}

	kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath)

	if spec.Deploy {
		if err := kubectlCmd.Run(); err != nil {
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
			Enabled:      true,
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
	return values
}
