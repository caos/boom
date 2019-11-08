package loggingoperator

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/kubectl"
	"github.com/caos/toolsop/internal/template"
)

var (
	applicationName = "logging-operator"
	resultsFilename = "results.yaml"
)

type LoggingOperator struct {
	ApplicationDirectoryPath string
}

func New(toolsDirectoryPath string) *LoggingOperator {
	lo := &LoggingOperator{
		ApplicationDirectoryPath: strings.Join([]string{toolsDirectoryPath, applicationName}, "/"),
	}

	return lo
}

func (l *LoggingOperator) Reconcile(overlay string, helm *template.Helm, logopspec *toolsetsv1beta1.LoggingOperator) error {

	resultFilePath := strings.Join([]string{l.ApplicationDirectoryPath, resultsFilename}, "/")

	values := specToValues(helm.GetImageTags(applicationName), logopspec)
	writeValues := func(path string) error {
		if err := helper.StructToYaml(values, path); err != nil {
			return err
		}
		return nil
	}

	if err := helm.Template(applicationName, logopspec.Prefix, logopspec.Namespace, resultFilePath, writeValues); err != nil {
		return err
	}

	kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath).AddParameter("-n", logopspec.Namespace)

	if err := kubectlCmd.Run(); err != nil {
		return err
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
