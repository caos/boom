package argocd

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/helm"
	"github.com/caos/boom/internal/templator/helm/chart"
)

func (a *Argocd) SpecToHelmValues(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolsetCRDSpec.Argocd

	imageTags := a.GetImageTags()
	values := helm.DefaultValues(imageTags)
	if spec.CustomImageWithGopass {
		imageRepository := "docker.pkg.github.com/caos/argocd-secrets/argocd"

		values.RepoServer.Image.Repository = imageRepository
		values.RepoServer.Image.Tag = imageTags[imageRepository]
	}

	return values
}

func (a *Argocd) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (a *Argocd) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
