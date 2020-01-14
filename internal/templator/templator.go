package templator

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/name"
)

type Templator interface {
	Template(interface{}, *v1beta1.ToolsetSpec, func(string) error) Templator
	GetResultsFilePath(name.Application, string, string) string
	CleanUp() Templator
	GetStatus() error
}
