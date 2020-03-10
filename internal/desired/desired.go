package desired

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kustomize"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/boom/internal/name"
	"github.com/caos/orbiter/mntr"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Apply(monitor mntr.Monitor, resultFilePath, namespace string, appName name.Application, force bool) error {
	resultFileDirPath := filepath.Dir(resultFilePath)

	if err := prepareAdditionalFiles(resultFilePath, namespace, appName); err != nil {
		return err
	}

	// apply resources
	cmd, err := kustomize.New(resultFileDirPath, true, force)
	if err != nil {
		return err
	}

	return errors.Wrapf(helper.Run(monitor, cmd.Build()), "Failed to apply with file %s", resultFilePath)
}

func Get(monitor mntr.Monitor, resultFilePath, namespace string, appName name.Application) ([]*helper.Resource, error) {
	resultFileDirPath := filepath.Dir(resultFilePath)

	if err := prepareAdditionalFiles(resultFilePath, namespace, appName); err != nil {
		return nil, err
	}

	// apply resources
	cmd, err := kustomize.New(resultFileDirPath, false, false)
	if err != nil {
		return nil, err
	}

	out, err := helper.RunWithOutput(monitor, cmd.Build())
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to apply with file %s", resultFilePath)
	}

	resources := make([]*helper.Resource, 0)

	parts := strings.Split(string(out), "\n---\n")
	for _, part := range parts {
		if part == "" {
			continue
		}
		var resource helper.Resource

		if err := yaml.Unmarshal([]byte(part), &resource); err != nil {
			return nil, err
		}

		resources = append(resources, &resource)
	}

	return resources, nil
}

func prepareAdditionalFiles(resultFilePath, namespace string, appName name.Application) error {
	resultFileDirPath := filepath.Dir(resultFilePath)

	resultFileKustomizePath := filepath.Join(resultFileDirPath, "kustomization.yaml")
	resultFileTransformerPath := filepath.Join(resultFileDirPath, "transformer.yaml")

	if helper.FileExists(resultFileKustomizePath) {
		if err := os.Remove(resultFileKustomizePath); err != nil {
			return err
		}
	}

	if helper.FileExists(resultFileTransformerPath) {
		if err := os.Remove(resultFileTransformerPath); err != nil {
			return err
		}
	}

	transformer := &kustomize.LabelTransformer{
		ApiVersion: "builtin",
		Kind:       "LabelTransformer",
		Metadata: &kustomize.Metadata{
			Name: "LabelTransformer",
		},
		Labels:     labels.GetAllApplicationLabels(appName),
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
	return nil
}
