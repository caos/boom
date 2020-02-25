package bundle

import (
	"path/filepath"

	"github.com/caos/boom/internal/bundle/application"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kustomize"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

func apply(logger logging.Logger, app application.Application) func(resultFilePath, namespace string) error {

	logFields := map[string]interface{}{
		"command": "apply",
	}

	resultFunc := func(resultFilePath, namespace string) error {
		resultFileDirPath := filepath.Dir(resultFilePath)
		resultFileKustomizePath := filepath.Join(resultFileDirPath, "kustomization.yaml")
		resultFileTransformerPath := filepath.Join(resultFileDirPath, "transformer.yaml")

		transformer := &kustomize.LabelTransformer{
			ApiVersion: "builtin",
			Kind:       "LabelTransformer",
			Metadata: &kustomize.Metadata{
				Name: "LabelTransformer",
			},
			Labels:     labels.GetApplicationLabels(app.GetName()),
			FieldSpecs: []*kustomize.FieldSpec{&kustomize.FieldSpec{Path: "metadata/labels", Create: true}},
		}
		if err := helper.AddStructToYaml(resultFileTransformerPath, transformer); err != nil {
			return err
		}

		kustomizeFile := kustomize.File{
			Namespace:    "caos-system",
			Resources:    []string{filepath.Base(resultFilePath)},
			Transformers: []string{filepath.Base(resultFileTransformerPath)},
		}
		if err := helper.AddStructToYaml(resultFileKustomizePath, kustomizeFile); err != nil {
			return err
		}

		// apply resources
		cmd, err := kustomize.New(resultFileDirPath, true)
		if err != nil {
			return err
		}
		err = helper.Run(logger.WithFields(logFields), cmd.Build())
		if err != nil {
			return errors.Wrapf(err, "Failed to apply with file %s", resultFilePath)
		}

		//TODO cleanup unnecessary resources
		return nil
	}

	return resultFunc
}
