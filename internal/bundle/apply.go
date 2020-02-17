package bundle

import (
	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

func apply(logger logging.Logger, app application.Application) func(resultFilePath, namespace string) error {

	logFields := map[string]interface{}{
		"application": app.GetName().String,
		"logID":       "CRD-7ifHfIFSKZ2jKgZ",
		"command":     "apply",
	}

	resultFunc := func(resultFilePath, namespace string) error {
		// apply resources
		kubectlCmd := kubectl.New("apply").AddParameter("-f", resultFilePath)
		if namespace != "" {
			kubectlCmd.AddParameter("-n", namespace)
		}
		err := helper.Run(logger.WithFields(logFields), kubectlCmd.Build())
		if err != nil {
			return errors.Wrapf(err, "Failed to apply with file %s", resultFilePath)
		}
		// label resources
		kubectlLabel := kubectl.NewLabel(resultFilePath)
		if namespace != "" {
			kubectlLabel.AddParameter("-n", namespace)
		}
		err = kubectlLabel.Apply(logger.WithFields(logFields), labels.GetApplicationLabels(app.GetName()))
		return errors.Wrapf(err, "Failed to label with file %s", resultFilePath)

		//TODO cleanup unnecessary resources

	}

	return resultFunc
}
