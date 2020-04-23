package admin

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/helm"
	"github.com/caos/boom/internal/bundle/application/applications/grafana/info"
	"github.com/caos/boom/internal/bundle/application/resources"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/labels"
	"strings"
)

func getSecretName() string {
	return strings.Join([]string{"grafana", "admin"}, "-")
}

func getUserKey() string {
	return "username"
}

func getPasswordKey() string {
	return "password"
}

func GetSecrets(adminSpec *toolsetsv1beta1.Admin) []interface{} {
	namespace := "caos-system"

	secrets := make([]interface{}, 0)

	if !helper.IsExistentClientSecret(adminSpec.ExistingSecret) {
		data := map[string]string{
			getUserKey():     adminSpec.Username.Value,
			getPasswordKey(): adminSpec.Password.Value,
		}

		conf := &resources.SecretConfig{
			Name:      getSecretName(),
			Namespace: namespace,
			Labels:    labels.GetAllApplicationLabels(info.GetName()),
			Data:      data,
		}
		secretRes := resources.NewSecret(conf)
		secrets = append(secrets, secretRes)
	}
	return secrets
}

func GetConfig(adminSpec *toolsetsv1beta1.Admin) *helm.Admin {
	if helper.IsExistentClientSecret(adminSpec.ExistingSecret) {

		return &helm.Admin{
			ExistingSecret: adminSpec.ExistingSecret.Name,
			UserKey:        adminSpec.ExistingSecret.IDKey,
			PasswordKey:    adminSpec.ExistingSecret.SecretKey,
		}
	}

	return &helm.Admin{
		ExistingSecret: getSecretName(),
		UserKey:        getUserKey(),
		PasswordKey:    getPasswordKey(),
	}

}
