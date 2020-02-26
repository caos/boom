package yaml

import (
	"path/filepath"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/templator"
)

func (y *YAML) Template(appInterface interface{}, spec *v1beta1.ToolsetSpec, resultFunc func(string, string) error) templator.Templator {
	if y.GetStatus() != nil {
		return y
	}

	app, err := checkTemplatorInterface(appInterface)
	if err != nil {
		y.status = err
		return y
	}

	yamlInterface := app.GetYaml(y.monitor, spec)
	resultfilepath := y.GetResultsFilePath(app.GetName(), y.overlay, y.templatorDirectoryPath)
	resultfiledirectory := y.getResultsFileDirectory(app.GetName(), y.overlay, y.templatorDirectoryPath)

	resultAbsFilePath, err := filepath.Abs(resultfilepath)
	if err != nil {
		y.status = err
		return y
	}
	resultAbsFileDirectory, err := filepath.Abs(resultfiledirectory)
	if err != nil {
		y.status = err
		return y
	}

	if err := helper.RecreatePath(resultAbsFileDirectory); err != nil {
		y.status = err
		return y
	}

	if yamlStr, isString := yamlInterface.(string); isString {
		y.status = helper.AddStringObjectToYaml(resultAbsFilePath, yamlStr)
	} else {
		y.status = helper.AddStringObjectToYaml(resultAbsFilePath, yamlStr)
	}
	if y.GetStatus() != nil {
		return y
	}

	y.status = resultFunc(resultAbsFilePath, "")
	return y
}
