package app

import (
	"strings"

	"github.com/caos/toolsop/internal/git"
	"github.com/caos/toolsop/internal/helper"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
)

func (a *App) MaintainGitCrd(directoryPath, url, secretPath, crdPath string) error {
	if url != "" {
		if _, err := git.CloneRepo(directoryPath, url, secretPath); err != nil {
			return err
		}
	}

	crdFilePath := strings.Join([]string{directoryPath, crdPath}, "/")

	toolsetCRD := &toolsetsv1beta1.Toolset{}
	if err := helper.YamlToStruct(crdFilePath, toolsetCRD); err != nil {
		return err
	}

	if err := a.GenerateTemplateComponents("caos", toolsetCRD.Spec.Name, toolsetCRD.Spec.Version); err != nil {
		return err
	}

	if err := a.Reconcile("caos", toolsetCRD.Spec); err != nil {
		return err
	}

	return nil
}
