package bundle

import (
	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

func deleteWithLabels(logger logging.Logger, app application.Application) error {
	logFields := map[string]interface{}{
		"application": app.GetName().String,
		"logID":       "CRD-4SwuPEofj23d6cP",
		"command":     "delete",
	}

	allApplicationResources, err := clientgo.ListResources(logger, labels.GetApplicationLabels(app.GetName()))
	if err != nil {
		err := errors.Wrap(err, "Failed to get current resources")
		logger.WithFields(logFields).Error(err)
		return err
	}

	for _, resource := range allApplicationResources {
		if err := clientgo.DeleteResource(resource); err != nil {
			err := errors.Wrap(err, "Failed to delete resource")
			logger.WithFields(logFields).Error(err)
			return err
		}
	}

	return nil
}

func delete(logger logging.Logger, app application.Application) func(resultFilePath, namespace string) error {

	logFields := map[string]interface{}{
		"application": app.GetName().String,
		"logID":       "CRD-E9Hu4zKkaQkzuzo",
		"command":     "delete",
	}

	resultFunc := func(resultFilePath, namespace string) error {
		kubectlCmd := kubectl.New("delete").AddParameter("-f", resultFilePath).AddFlag("--ignore-not-found")
		if namespace != "" {
			kubectlCmd.AddParameter("-n", namespace)
		}
		err := helper.Run(logger.WithFields(logFields), kubectlCmd.Build())
		return errors.Wrapf(err, "Failed to delete with file %s", resultFilePath)

	}
	return resultFunc
}
