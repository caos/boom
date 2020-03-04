package customimage

import (
	"encoding/json"
	"path/filepath"
	"strings"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	"github.com/pkg/errors"
)

const (
	tab           string = "  "
	nl            string = "\n"
	sshFolderName string = "/home/argocd/ssh-keys"
	gpgFolderName string = "/home/argocd/gpg-import"
)

type SecretVolume struct {
	Name        string  `yaml:"name"`
	Secret      *Secret `yaml:"secret,omitempty"`
	DefaultMode int     `yaml:"defaultMode"`
}

type Secret struct {
	SecretName string  `yaml:"secretName,omitempty"`
	Items      []*Item `yaml:"items,omitempty"`
}

type Item struct {
	Key  string `yaml:"key"`
	Path string `yaml:"path"`
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
	if spec.CustomImage.GopassGPGKey != nil {
		vol := &SecretVolume{
			Name: spec.CustomImage.GopassGPGKey.InternalName,
			Secret: &Secret{
				SecretName: spec.CustomImage.GopassGPGKey.Name,
				Items: []*Item{&Item{
					Key:  spec.CustomImage.GopassGPGKey.Key,
					Path: spec.CustomImage.GopassGPGKey.InternalName,
				},
				},
			},
			DefaultMode: 0544,
		}
		vols = append(vols, vol)
		volMount := &VolumeMount{
			Name:      spec.CustomImage.GopassGPGKey.InternalName,
			MountPath: gpgFolderName,
			ReadOnly:  false,
		}
		volMounts = append(volMounts, volMount)
	}

	if spec.CustomImage.GopassSSHKey != nil {
		vol := &SecretVolume{
			Name: spec.CustomImage.GopassSSHKey.InternalName,
			Secret: &Secret{
				SecretName: spec.CustomImage.GopassSSHKey.Name,
				Items: []*Item{&Item{
					Key:  spec.CustomImage.GopassSSHKey.Key,
					Path: spec.CustomImage.GopassSSHKey.InternalName,
				},
				},
			},
			DefaultMode: 0544,
		}
		vols = append(vols, vol)
		volMount := &VolumeMount{
			Name:      spec.CustomImage.GopassSSHKey.InternalName,
			MountPath: sshFolderName,
			ReadOnly:  false,
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

type stores struct {
	Stores []*store `json:"stores"`
}

type store struct {
	Directory string `json:"directory"`
	StoreName string `json:"storename"`
}

func AddPostStartFromSpec(spec *toolsetsv1beta1.Argocd, resultFilePath string) error {
	stores := &stores{}
	for _, v := range spec.CustomImage.GopassStores {
		stores.Stores = append(stores.Stores, &store{Directory: v.Directory, StoreName: v.StoreName})
	}
	jsonStores, err := json.Marshal(stores)
	if err != nil {
		return errors.Wrap(err, "Error while marshaling gopass stores in json")
	}
	jsonStoresStr := strings.ReplaceAll(string(jsonStores), "\"", "\\\"")

	gpgFileName := filepath.Join(gpgFolderName, spec.CustomImage.GopassGPGKey.InternalName)
	sshFileName := filepath.Join(sshFolderName, spec.CustomImage.GopassSSHKey.InternalName)

	addCommand := strings.Join([]string{"/home/argocd/initialize_gopass.sh '", jsonStoresStr, "' '", gpgFileName, "' '", sshFileName, "'"}, "")
	addLifecycle := strings.Join([]string{
		tab, tab, tab, tab, "lifecycle:", nl,
		tab, tab, tab, tab, tab, "postStart:", nl,
		tab, tab, tab, tab, tab, tab, "exec:", nl,
		tab, tab, tab, tab, tab, tab, tab, "command: [\"/bin/bash\", \"-c\", \"", addCommand, "\"]", nl,
	}, "")

	return helper.AddStringBeforePointForKindAndName(resultFilePath, "Deployment", "argocd-repo-server", "imagePullPolicy:", addLifecycle)
}
