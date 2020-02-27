package bundle

import (
	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/orbiter/mntr"
	"github.com/pkg/errors"
)

func deleteWithCurrentState(monitor mntr.Monitor, resourceInfoList []*clientgo.ResourceInfo, app application.Application) func(resultFilePath, namespace string) error {
	logFields := map[string]interface{}{
		"application": app.GetName().String,
		"command":     "delete",
	}
	delMonitor := monitor.WithFields(logFields)

	resultFunc := func(resultFilePath, namespace string) error {
		allApplicationResources, err := clientgo.ListResources(delMonitor, resourceInfoList, labels.GetApplicationLabels(app.GetName()))
		if err != nil {
			err := errors.Wrap(err, "Failed to get current resources")
			delMonitor.Error(err)
			return err
		}

		allForApplicationResources, err := clientgo.ListResources(delMonitor, resourceInfoList, labels.GetForApplicationLabels(app.GetName()))
		if err != nil {
			err := errors.Wrap(err, "Failed to get current for-application resources")
			delMonitor.Error(err)
			return err
		}
		for _, app := range allForApplicationResources {
			allApplicationResources = append(allApplicationResources, app)
		}

		for _, resource := range allApplicationResources {
			if err := clientgo.DeleteResource(resource); err != nil {
				err := errors.Wrap(err, "Failed to delete resource")
				delMonitor.Error(err)
				return err
			}
		}
		return nil
	}

	return resultFunc
}

func delete(monitor mntr.Monitor) func(resultFilePath, namespace string) error {

	logFields := map[string]interface{}{
		"command": "delete",
	}
	delMonitor := monitor.WithFields(logFields)

	resultFunc := func(resultFilePath, namespace string) error {
		kubectlCmd := kubectl.New("delete").AddParameter("-f", resultFilePath).AddFlag("--ignore-not-found")
		if namespace != "" {
			kubectlCmd.AddParameter("-n", namespace)
		}
		err := helper.Run(delMonitor, kubectlCmd.Build())
		return errors.Wrapf(err, "Failed to delete with file %s", resultFilePath)

	}
	return resultFunc
}
