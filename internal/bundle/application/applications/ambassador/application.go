package ambassador

import (
	"github.com/caos/boom/internal/name"
)

const (
	applicationName name.Application = "ambassador"
)

func GetName() name.Application {
	return applicationName
}
