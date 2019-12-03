package apps

/*
import (
	"os"
	"path/filepath"
	"strings"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/kubectl"
	"github.com/caos/toolsop/internal/template"
)

func Reconcile(logger logging.Logger, applicationDirectoryPath, resultsDirectoryName, resultsFileName, applicationName, defaultNamespace, overlay string, helm *template.Helm, spec *toolsetsv1beta1.Ambassador) error {


	logFields := map[string]interface{}{
		"application": applicationName,
	}

	logFields["logID"] = "CRD-rGkpjHLZtVAWumr"
	logger.WithFields(logFields).Info("Reconciling")
	resultsFileDirectory := filepath.Join(applicationDirectoryPath, resultsDirectoryName, overlay)
	_ = os.RemoveAll(resultsFileDirectory)
	_ = os.MkdirAll(resultsFileDirectory, os.ModePerm)
	resultFilePath := filepath.Join(resultsFileDirectory, resultsFileName)

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
		if err := errors.Wrapf(helper.StructToYaml(values, path), "Failed to write values file overlay %s application %s", overlay, applicationName); err != nil {
			logFields["logID"] = "CRD-se7ejQ2L9uj5pv1"
			logger.WithFields(logFields).Error(err)
			return err
		}
		return nil
	}

	if err := helm.Template(applicationName, prefix, namespace, resultFilePath, writeValues); err != nil {
		return err
	}

	kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath).AddParameter("-n", namespace)

	if spec.Deploy {
		if err := errors.Wrapf(kubectlCmd.Run(), "Failed to apply file %s", resultFilePath); err != nil {
			logFields["logID"] = "CRD-auZcsGJbTM8gahX"
			logger.WithFields(logFields).Error(err)
			return err
		}
	}

	return nil
}
*/
