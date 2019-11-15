package loggingoperator

import (
	"os"
	"path/filepath"
	"strings"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/kubectl"
	"github.com/caos/toolsop/internal/template"
)

var (
	applicationName      = "logging-operator"
	resultsDirectoryName = "results"
	resultsFileName      = "results.yaml"
	defaultNamespace     = "logging"
)

type LoggingOperator struct {
	ApplicationDirectoryPath string
}

func New(toolsDirectoryPath string) *LoggingOperator {
	lo := &LoggingOperator{
		ApplicationDirectoryPath: filepath.Join(toolsDirectoryPath, applicationName),
	}

	return lo
}

func (l *LoggingOperator) Reconcile(overlay string, helm *template.Helm, spec *toolsetsv1beta1.LoggingOperator) error {
	resultsFileDirectory := filepath.Join(l.ApplicationDirectoryPath, resultsDirectoryName, overlay)
	_ = os.RemoveAll(resultsFileDirectory)
	_ = os.MkdirAll(resultsFileDirectory, os.ModePerm)
	resultFilePath := filepath.Join(resultsFileDirectory, resultsFileName)

	values := specToValues(helm.GetImageTags(applicationName), spec)
	writeValues := func(path string) error {
		if err := helper.StructToYaml(values, path); err != nil {
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
			return err
		}
	}

	return nil
}

func specToValues(imageTags map[string]string, spec *toolsetsv1beta1.LoggingOperator) *Values {
	values := &Values{
		ReplicaCount: 1,
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
