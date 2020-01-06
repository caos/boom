package loggingoperator

import (
	"path/filepath"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/v1beta1/crd/defaults"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/template"
)

var (
	applicationName = "logging-operator"
)

type LoggingOperator struct {
	ApplicationDirectoryPath string
	logger                   logging.Logger
	spec                     *toolsetsv1beta1.LoggingOperator
}

func New(logger logging.Logger, toolsDirectoryPath string) *LoggingOperator {
	lo := &LoggingOperator{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
		logger:                   logger,
	}

	return lo
}

func (l *LoggingOperator) Reconcile(overlay, specNamespace string, helm *template.Helm, spec *toolsetsv1beta1.LoggingOperator) error {

	logFields := map[string]interface{}{
		"application": applicationName,
		"logID":       "CRD-8Z2ueDkBmkBONgc",
	}

	l.logger.WithFields(logFields).Info("Reconciling")

	resultFilePath := defaults.GetResultFilePath(overlay, l.ApplicationDirectoryPath, applicationName)
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
		if err := defaults.PrepareForResultOutput(defaults.GetResultFileDirectory(overlay, l.ApplicationDirectoryPath, applicationName)); err != nil {
			return err
		}

		if err := helm.Template(applicationName, resultFilePath); err != nil {
			return err
		}

		if err := helper.DeleteKindFromYaml(resultFilePath, "Namespace"); err != nil {
			return err
		}

		kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)

		if err := errors.Wrapf(helper.Run(l.logger, kubectlCmd.Build()), "Failed to apply file %s", resultFilePath); err != nil {
			return err
		}

		l.spec = spec
	} else if !spec.Deploy && l.spec != nil {
		kubectlCmd := kubectl.New("delete").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)

		if err := errors.Wrapf(helper.Run(l.logger, kubectlCmd.Build()), "Failed to apply file %s", resultFilePath); err != nil {
			return err
		}

		l.spec = nil
	}

	return nil
}

func specToValues(imageTags map[string]string, spec *toolsetsv1beta1.LoggingOperator) *Values {
	values := &Values{
		FullnameOverride: "logging-operator",
		ReplicaCount:     1,
		Image: Image{
			Repository: "banzaicloud/logging-operator",
			Tag:        imageTags["banzaicloud/logging-operator"],
			PullPolicy: "IfNotPresent",
		},
		HTTP: HTTP{
			Port: 8080,
			Service: Service{
				Type: "ClusterIP",
			},
		},
		RBAC: RBAC{
			Enabled: true,
			PSP: PSP{
				Enabled: true,
			},
		},
	}
	return values
}
