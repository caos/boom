package kubestatemetrics

import (
	"path/filepath"

	"github.com/caos/boom/internal/app/v1beta1/crd/defaults"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/template"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

var (
	applicationName = "kube-state-metrics"
)

type KubeStateMetrics struct {
	ApplicationDirectoryPath string
	logger                   logging.Logger
	spec                     *toolsetsv1beta1.KubeStateMetrics
}

func New(logger logging.Logger, toolsDirectoryPath string) *KubeStateMetrics {
	lo := &KubeStateMetrics{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
		logger:                   logger,
	}

	return lo
}

func (k *KubeStateMetrics) Reconcile(overlay, specNamespace string, helm *template.Helm, spec *toolsetsv1beta1.KubeStateMetrics) error {

	logFields := map[string]interface{}{
		"application": applicationName,
		"logID":       "CRD-8Z2ueDkBmkBONgc",
	}

	k.logger.WithFields(logFields).Info("Reconciling")

	resultFilePath := defaults.GetResultFilePath(overlay, k.ApplicationDirectoryPath, applicationName)
	prefix := defaults.GetPrefix(overlay, applicationName, spec.Prefix)
	namespace := defaults.GetNamespace(overlay, applicationName, specNamespace, spec.Namespace)

	values := specToValues(helm.GetImageTags(applicationName), spec)
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
		if err := defaults.PrepareForResultOutput(defaults.GetResultFileDirectory(overlay, k.ApplicationDirectoryPath, applicationName)); err != nil {
			return err
		}

		if err := helm.Template(applicationName, resultFilePath); err != nil {
			return err
		}

		if err := helper.DeleteKindFromYaml(resultFilePath, "Namespace"); err != nil {
			return err
		}

		kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)

		if err := errors.Wrapf(helper.Run(k.logger, kubectlCmd.Build()), "Failed to apply file %s", resultFilePath); err != nil {
			return err
		}

		k.spec = spec
	} else if !spec.Deploy && k.spec != nil {
		kubectlCmd := kubectl.New("delete").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)

		if err := errors.Wrapf(helper.Run(k.logger, kubectlCmd.Build()), "Failed to apply file %s", resultFilePath); err != nil {
			return err
		}

		k.spec = nil
	}

	return nil
}

func specToValues(imageTags map[string]string, spec *toolsetsv1beta1.KubeStateMetrics) *Values {
	values := &Values{
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

	return values
}
