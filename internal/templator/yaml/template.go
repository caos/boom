package yaml

import (
	"fmt"
	"path/filepath"

	"github.com/caos/boom/api/v1beta1"
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

	yamlInterface := app.GetYaml()
	resultfilepath := y.GetResultsFilePath(app.GetName(), y.overlay, y.templatorDirectoryPath)

	resultAbsFilePath, err := filepath.Abs(resultfilepath)
	if err != nil {
		y.status = err
		return y
	}

	fmt.Println(yamlInterface)
	// if yamlStr, isString := yamlInterface.(string); isString {
	// 	y.status = helper.AddStringObjectToYaml(resultAbsFilePath, yamlStr)
	// } else {
	// 	y.status = helper.AddStringObjectToYaml(resultAbsFilePath, yamlStr)
	// }
	// if y.GetStatus() != nil {
	// 	return y
	// }

	y.status = resultFunc(resultAbsFilePath, app.GetNamespace())
	return y
}
