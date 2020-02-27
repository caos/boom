package bundle

import (
	"path/filepath"

	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/desired"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/orbiter/mntr"
	"github.com/pkg/errors"
)

func applyWithCurrentState(monitor mntr.Monitor, resourceInfoList []*clientgo.ResourceInfo, app application.Application) func(resultFilePath, namespace string) error {

	logFields := map[string]interface{}{
		"command": "apply",
	}
	applyMonitor := monitor.WithFields(logFields)

	resultFunc := func(resultFilePath, namespace string) error {
		applyFunc := apply(monitor, app)

		desiredResources, err := desired.Get(monitor, resultFilePath, namespace, app.GetName())
		if err != nil {
			return err
		}

		currentApplicationResources, err := clientgo.ListResources(applyMonitor, resourceInfoList, labels.GetApplicationLabels(app.GetName()))
		if err != nil {
			err := errors.Wrap(err, "Failed to get current resources")
			applyMonitor.Error(err)
			return err
		}

		currentForApplicationResources, err := clientgo.ListResources(applyMonitor, resourceInfoList, labels.GetForApplicationLabels(app.GetName()))
		if err != nil {
			err := errors.Wrap(err, "Failed to get current for-application resources")
			applyMonitor.Error(err)
			return err
		}
		for _, app := range currentForApplicationResources {
			currentApplicationResources = append(currentApplicationResources, app)
		}

		if err := applyFunc(resultFilePath, namespace); err != nil {
			return err
		}
		deleteResources := make([]*clientgo.Resource, 0)
		for _, currentResource := range currentApplicationResources {
			found := false
			for _, desiredResource := range desiredResources {
				apiVersion := filepath.Join(currentResource.Group, currentResource.Version)
				if desiredResource.ApiVersion == apiVersion &&
					desiredResource.Kind == currentResource.Kind &&
					desiredResource.Metadata.Name == currentResource.Name &&
					(currentResource.Namespace == "" || desiredResource.Metadata.Namespace == currentResource.Namespace) {
					found = true
					break
				}
			}
			if found == false {
				otherAPI := false
				for _, desiredResource := range desiredResources {
					apiVersion := filepath.Join(currentResource.Group, currentResource.Version)
					if apiVersion != desiredResource.ApiVersion &&
						currentResource.Kind == desiredResource.Kind &&
						currentResource.Name == desiredResource.Metadata.Name &&
						(currentResource.Namespace == "" || desiredResource.Metadata.Namespace == currentResource.Namespace) {
						otherAPI = true
						break
					}
				}

				if !otherAPI {
					deleteResources = append(deleteResources, currentResource)
				}
			}
		}

		if deleteResources != nil && len(deleteResources) > 0 {
			for _, deleteResource := range deleteResources {
				if err := clientgo.DeleteResource(deleteResource); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return resultFunc
}

func apply(monitor mntr.Monitor, app application.Application) func(resultFilePath, namespace string) error {

	logFields := map[string]interface{}{
		"command": "apply",
	}
	applyMonitor := monitor.WithFields(logFields)

	resultFunc := func(resultFilePath, namespace string) error {
		return desired.Apply(applyMonitor, resultFilePath, namespace, app.GetName())
	}

	return resultFunc
}
