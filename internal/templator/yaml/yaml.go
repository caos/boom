package yaml

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator"
)

type Templator interface {
	GetName() name.Application
	SpecToYAML(spec *v1beta1.ToolsetSpec) interface{}
}

type YAML struct {
}

func (y *YAML) Template(appInterface interface{}, spec *v1beta1.ToolsetSpec, resultFunc func(string) error) templator.Templator {

	return nil
}
func (y *YAML) GetResultsFilePath(appName name.Application, overlay, basePath string) string {
	return ""
}
func (y *YAML) CleanUp() templator.Templator {
	return nil
}
func (y *YAML) GetStatus() error {
	return nil
}
