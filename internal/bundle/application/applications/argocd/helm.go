package argocd

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/auth"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/customimage"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/helm"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/orbiter/logging"
	"gopkg.in/yaml.v3"
)

func (a *Argocd) HelmMutate(logger logging.Logger, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec, resultFilePath string) error {
	spec := toolsetCRDSpec.Argocd

	if spec.CustomImage.Enabled && spec.CustomImage.ImagePullSecret != "" {
		if err := customimage.AddImagePullSecretFromSpec(spec, resultFilePath); err != nil {
			return err
		}

		if spec.CustomImage.GopassDirectory != "" && spec.CustomImage.GopassStoreName != "" {
			if err := customimage.AddPostStartFromSpec(spec, resultFilePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Argocd) SpecToHelmValues(logger logging.Logger, toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolsetCRDSpec.Argocd

	imageTags := a.GetImageTags()
	values := helm.DefaultValues(imageTags)
	if spec.CustomImage.Enabled {
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

	if spec.Auth.OIDC != nil {
		oidc, err := auth.GetOIDC(spec.Auth.OIDC)
		if err == nil {
			values.Server.Config.OIDC = oidc
		}
	}

	dexConfig := auth.GetDexConfigFromSpec(logger, spec)
	if dexConfig != nil {
		data, err := yaml.Marshal(dexConfig)
		if err == nil {
			values.Server.Config.Dex = string(data)
		}
		values.Dex = helm.DefaultDexValues(imageTags)
		values.Server.Config.URL = spec.Auth.RootURL
	}

	return values
}

func (a *Argocd) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (a *Argocd) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
