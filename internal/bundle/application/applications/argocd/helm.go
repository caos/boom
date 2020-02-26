package argocd

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/config"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/customimage"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/helm"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/mntr"
)

func (a *Argocd) HelmMutate(monitor mntr.Monitor, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec, resultFilePath string) error {
	spec := toolsetCRDSpec.Argocd

	if spec.CustomImage != nil && spec.CustomImage.Enabled && spec.CustomImage.ImagePullSecret != "" {
		if err := customimage.AddImagePullSecretFromSpec(spec, resultFilePath); err != nil {
			return err
		}

		if spec.CustomImage.GopassStores != nil && len(spec.CustomImage.GopassStores) > 0 {
			if err := customimage.AddPostStartFromSpec(spec, resultFilePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Argocd) SpecToHelmValues(monitor mntr.Monitor, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolsetCRDSpec.Argocd

	imageTags := a.GetImageTags()
	values := helm.DefaultValues(imageTags)
	if spec.CustomImage != nil && spec.CustomImage.Enabled {
		conf := customimage.FromSpec(spec, imageTags)
		values.RepoServer.Image = &helm.Image{
			Repository:      conf.ImageRepository,
			Tag:             conf.ImageTag,
			ImagePullPolicy: "IfNotPresent",
		}
		if conf.AddSecretVolumes != nil {
			for _, v := range conf.AddSecretVolumes {
				values.RepoServer.Volumes = append(values.RepoServer.Volumes, &helm.Volume{
					Secret: &helm.VolumeSecret{
						SecretName:  v.SecretName,
						DefaultMode: v.DefaultMode,
					},
					Name: v.Name,
				})
			}
		}
		if conf.AddVolumeMounts != nil {
			for _, v := range conf.AddVolumeMounts {
				values.RepoServer.VolumeMounts = append(values.RepoServer.VolumeMounts, &helm.VolumeMount{
					Name:      v.Name,
					MountPath: v.MountPath,
					SubPath:   v.SubPath,
					ReadOnly:  v.ReadOnly,
				})
			}
		}
	}

	conf := config.GetFromSpec(monitor, spec)
	if conf.Repositories != "" && conf.Repositories != "[]\n" {
		values.Server.Config.Repositories = conf.Repositories
	}

	if conf.ConfigManagementPlugins != "" {
		values.Server.Config.ConfigManagementPlugins = conf.ConfigManagementPlugins
	}

	if spec.Network != nil && spec.Network.Domain != "" {

		if conf.OIDC != "" {
			values.Server.Config.OIDC = conf.OIDC
		}

		if conf.Connectors != "" && conf.Connectors != "{}\n" {
			values.Server.Config.Dex = conf.Connectors

			values.Dex = helm.DefaultDexValues(imageTags)
			values.Server.Config.URL = strings.Join([]string{"https://", spec.Network.Domain}, "")
		}
	}

	return values
}

func (a *Argocd) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (a *Argocd) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
