package customimage

import (
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
)

const (
	tab string = "  "
	nl  string = "\n"
)

type SecretVolume struct {
	SecretName  string `yaml:"secretName,omitempty"`
	Name        string `yaml:"name"`
	DefaultMode int    `yaml:"defaultMode"`
}

type VolumeMount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath,omitempty"`
	SubPath   string `yaml:"subPath,omitempty"`
	ReadOnly  bool   `yaml:"readOnly,omitempty"`
}

type CustomImage struct {
	ImageRepository  string
	ImageTag         string
	AddSecretVolumes []*SecretVolume
	AddVolumeMounts  []*VolumeMount
}

func FromSpec(spec *toolsetsv1beta1.Argocd, imageTags map[string]string) *CustomImage {
	imageRepository := "docker.pkg.github.com/caos/argocd-secrets/argocd"

	vols := make([]*SecretVolume, 0)
	volMounts := make([]*VolumeMount, 0)
	if spec.CustomImage.GopassGPGKey != "" {
		vol := &SecretVolume{
			Name:        spec.CustomImage.GopassGPGKey,
			SecretName:  spec.CustomImage.GopassGPGKey,
			DefaultMode: 0444,
		}
		vols = append(vols, vol)
		volMount := &VolumeMount{
			Name:      spec.CustomImage.GopassGPGKey,
			MountPath: "/home/argocd/gpg-import",
			ReadOnly:  true,
		}
		volMounts = append(volMounts, volMount)
	}

	if spec.CustomImage.GopassSSHKey != "" {
		vol := &SecretVolume{
			Name:        spec.CustomImage.GopassSSHKey,
			SecretName:  spec.CustomImage.GopassSSHKey,
			DefaultMode: 0444,
		}
		vols = append(vols, vol)
		volMount := &VolumeMount{
			Name:      spec.CustomImage.GopassSSHKey,
			MountPath: "/home/argocd/ssh-key",
			ReadOnly:  true,
		}
		volMounts = append(volMounts, volMount)
	}

	return &CustomImage{
		ImageRepository:  imageRepository,
		ImageTag:         imageTags[imageRepository],
		AddSecretVolumes: vols,
		AddVolumeMounts:  volMounts,
	}
}

func AddImagePullSecretFromSpec(spec *toolsetsv1beta1.Argocd, resultFilePath string) error {
	addContent := strings.Join([]string{
		tab, tab, tab, "imagePullSecrets:", nl,
		tab, tab, tab, "- name: ", spec.CustomImage.ImagePullSecret, nl,
	}, "")

	return helper.AddStringBeforePointForKindAndName(resultFilePath, "Deployment", "argocd-repo-server", "volumes:", addContent)
}

func AddPostStartFromSpec(spec *toolsetsv1beta1.Argocd, resultFilePath string) error {
	addCommand := strings.Join([]string{"/home/argocd/initialize_gopass.sh", spec.CustomImage.GopassDirectory, spec.CustomImage.GopassStoreName}, " ")
	addLifecycle := strings.Join([]string{
		tab, tab, tab, tab, "lifecycle:", nl,
		tab, tab, tab, tab, tab, "postStart:", nl,
		tab, tab, tab, tab, tab, tab, "exec:", nl,
		tab, tab, tab, tab, tab, tab, tab, "command: [\"/bin/bash\", \"-c\", \"", addCommand, "\"]", nl,
	}, "")

	return helper.AddStringBeforePointForKindAndName(resultFilePath, "Deployment", "argocd-repo-server", "imagePullPolicy:", addLifecycle)
}
