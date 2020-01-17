package grafana

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "grafana"
)

func GetName() name.Application {
	return applicationName
}
