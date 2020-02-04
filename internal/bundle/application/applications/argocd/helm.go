package argocd

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/auth"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/helm"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/templator/helm/chart"
	"gopkg.in/yaml.v3"
)

func (a *Argocd) HelmMutate(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec, resultFilePath string) error {
	if toolsetCRDSpec.Argocd.CustomImageWithGopass && toolsetCRDSpec.Argocd.ImagePullSecret != "" {
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

type connector struct {
	Type   string
	Name   string
	ID     string
	Config interface{}
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

	if spec.Auth.OIDC != nil {
		oidc, err := auth.GetOIDC(spec.Auth.OIDC)
		if err == nil {
			values.Server.Config.OIDC = oidc
		}
	}

	dex := make([]*connector, 0)
	if spec.Auth.GithubConnector != nil {
		github, err := auth.GetGithub(spec.Auth.GithubConnector)
		if err == nil {
			dex = append(dex, &connector{
				Name:   spec.Auth.GithubConnector.Name,
				ID:     spec.Auth.GithubConnector.ID,
				Type:   "github",
				Config: github,
			})
		}
	}
	if spec.Auth.GitlabConnector != nil {
		gitlab, err := auth.GetGitlab(spec.Auth.GitlabConnector)
		if err == nil {
			dex = append(dex, &connector{
				Name:   spec.Auth.GitlabConnector.Name,
				ID:     spec.Auth.GitlabConnector.ID,
				Type:   "gitlab",
				Config: gitlab,
			})
		}
	}
	if spec.Auth.GoogleConnector != nil {
		google, err := auth.GetGoogle(spec.Auth.GoogleConnector)
		if err == nil {
			dex = append(dex, &connector{
				Name:   spec.Auth.GoogleConnector.Name,
				ID:     spec.Auth.GoogleConnector.ID,
				Type:   "google",
				Config: google,
			})
		}
	}

	if len(dex) > 0 {
		data, err := yaml.Marshal(dex)
		if err == nil {
			values.Server.Config.Dex = string(data)
		}
		values.Dex = helm.DefaultDexValues(imageTags)
	}

	return values
}

func (a *Argocd) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (a *Argocd) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
