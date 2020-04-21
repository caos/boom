package yaml

import (
	"path/filepath"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
)

func (y *YAML) Template(appInterface interface{}, spec *v1beta1.ToolsetSpec, resultFunc func(string, string) error) error {
	app, err := checkTemplatorInterface(appInterface)
	if err != nil {
		return err
	}

	yamlInterface := app.GetYaml(y.monitor, spec)
	resultfilepath := y.GetResultsFilePath(app.GetName(), y.overlay, y.templatorDirectoryPath)
	resultfiledirectory := y.getResultsFileDirectory(app.GetName(), y.overlay, y.templatorDirectoryPath)

	resultAbsFilePath, err := filepath.Abs(resultfilepath)
	if err != nil {
		return err
	}
	resultAbsFileDirectory, err := filepath.Abs(resultfiledirectory)
	if err != nil {
		return err
	}

	if err := helper.RecreatePath(resultAbsFileDirectory); err != nil {
		return err
	}

	if yamlStr, isString := yamlInterface.(string); isString {
		err = helper.AddStringObjectToYaml(resultAbsFilePath, yamlStr)
		if err != nil {
			return err
		}
	} else {
		err = helper.AddStructToYaml(resultAbsFilePath, yamlInterface)
		if err != nil {
			return err
		}
	}

	return resultFunc(resultAbsFilePath, "")
}
