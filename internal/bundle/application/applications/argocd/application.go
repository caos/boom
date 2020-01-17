package argocd

import "github.com/caos/boom/internal/name"

const (
	applicationName name.Application = "argocd"
)

func GetName() name.Application {
	return applicationName
}
