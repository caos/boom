package templator

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/name"
)

type Templator interface {
	Template(interface{}, *v1beta1.ToolsetSpec, func(string) error) Templator
	GetResultsFilePath(appName name.Application, overlay, basePath string) string
	CleanUp() Templator
	GetStatus() error
}
