package argocd

import (
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/helm"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/templator/helm/chart"
	"strings"
)

func (a *Argocd) HelmMutate(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec, resultFilePath string) error {
	if toolsetCRDSpec.Argocd.ImagePullSecret != "" {
		tab := "  "
		nl := "\n"
		addContent := strings.Join([]string{
			tab, tab, tab, "imagePullSecrets:", nl,
			tab, tab, tab, "- name: ", toolsetCRDSpec.Argocd.ImagePullSecret, nl,
		}, "")

		if err := helper.AddStringBeforePointForKindAndName(resultFilePath, "Deployment", "argocd-repo-server", "volumes:", addContent); err != nil {
			return err
		}

		if toolsetCRDSpec.Argocd.GopassDirectory != "" && toolsetCRDSpec.Argocd.GopassStoreName != "" {
			addCommand := strings.Join([]string{"/home/argocd/initialize_gopass.sh", toolsetCRDSpec.Argocd.GopassDirectory, toolsetCRDSpec.Argocd.GopassStoreName}, " ")
			addLifecycle := strings.Join([]string{
				tab, tab, tab, tab, "lifecycle:", nl,
				tab, tab, tab, tab, tab, "postStart:", nl,
				tab, tab, tab, tab, tab, tab, "exec:", nl,
				tab, tab, tab, tab, tab, tab, tab, "command: [\"/bin/bash\", \"-c\", \"", addCommand, "\"]", nl,
			}, "")

			if err := helper.AddStringBeforePointForKindAndName(resultFilePath, "Deployment", "argocd-repo-server", "imagePullPolicy:", addLifecycle); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Argocd) SpecToHelmValues(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolsetCRDSpec.Argocd

	imageTags := a.GetImageTags()
	values := helm.DefaultValues(imageTags)
	if spec.CustomImageWithGopass {
		imageRepository := "docker.pkg.github.com/caos/argocd-secrets/argocd"

		values.RepoServer.Image.Repository = imageRepository
		values.RepoServer.Image.Tag = imageTags[imageRepository]
		if spec.GopassGPGKey != "" {
			vol := &helm.Volume{
				Name: toolsetCRDSpec.Argocd.GopassGPGKey,
				Secret: &helm.VolumeSecret{
					SecretName:  toolsetCRDSpec.Argocd.GopassGPGKey,
					DefaultMode: 0444,
				},
			}
			values.RepoServer.Volumes = append(values.RepoServer.Volumes, vol)
			volMount := &helm.VolumeMount{
				Name:      toolsetCRDSpec.Argocd.GopassGPGKey,
				MountPath: "/home/argocd/gpg-import",
				ReadOnly:  true,
			}
			values.RepoServer.VolumeMounts = append(values.RepoServer.VolumeMounts, volMount)
		}

		if spec.GopassSSHKey != "" {
			vol := &helm.Volume{
				Name: toolsetCRDSpec.Argocd.GopassSSHKey,
				Secret: &helm.VolumeSecret{
					SecretName:  toolsetCRDSpec.Argocd.GopassSSHKey,
					DefaultMode: 0444,
				},
			}
			values.RepoServer.Volumes = append(values.RepoServer.Volumes, vol)
			volMount := &helm.VolumeMount{
				Name:      toolsetCRDSpec.Argocd.GopassSSHKey,
				MountPath: "/home/argocd/ssh-key",
				ReadOnly:  true,
			}
			values.RepoServer.VolumeMounts = append(values.RepoServer.VolumeMounts, volMount)
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
