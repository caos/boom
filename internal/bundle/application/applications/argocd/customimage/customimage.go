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
	for _, store := range spec.CustomImage.GopassStores {

		if store.GPGKey != nil {
			vol := &SecretVolume{
				Name: store.GPGKey.InternalName,
				Secret: &Secret{
					SecretName: store.GPGKey.Name,
					Items: []*Item{&Item{
						Key:  store.GPGKey.Key,
						Path: store.GPGKey.InternalName,
					},
					},
				},
				DefaultMode: 0544,
			}
			vols = append(vols, vol)
			mountPath := filepath.Join(gpgFolderName, store.GPGKey.InternalName)
			volMount := &VolumeMount{
				Name:      store.GPGKey.InternalName,
				MountPath: mountPath,
				SubPath:   store.GPGKey.InternalName,
				ReadOnly:  false,
			}
			volMounts = append(volMounts, volMount)
		}

		if store.SSHKey != nil {
			vol := &SecretVolume{
				Name: store.SSHKey.InternalName,
				Secret: &Secret{
					SecretName: store.SSHKey.Name,
					Items: []*Item{&Item{
						Key:  store.SSHKey.Key,
						Path: store.SSHKey.InternalName,
					},
					},
				},
				DefaultMode: 0544,
			}
			vols = append(vols, vol)
			mountPath := filepath.Join(sshFolderName, store.SSHKey.InternalName)
			volMount := &VolumeMount{
				Name:      store.SSHKey.InternalName,
				MountPath: mountPath,
				SubPath:   store.SSHKey.InternalName,
				ReadOnly:  false,
			}
			volMounts = append(volMounts, volMount)
		}
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

	addCommand := strings.Join([]string{"/home/argocd/initialize_gopass.sh '", jsonStoresStr, "' ", gpgFolderName, " ", sshFolderName}, "")
	addLifecycle := strings.Join([]string{
		tab, tab, tab, tab, "lifecycle:", nl,
		tab, tab, tab, tab, tab, "postStart:", nl,
		tab, tab, tab, tab, tab, tab, "exec:", nl,
		tab, tab, tab, tab, tab, tab, tab, "command: [\"/bin/bash\", \"-c\", \"", addCommand, "\"]", nl,
	}, "")

	return helper.AddStringBeforePointForKindAndName(resultFilePath, "Deployment", "argocd-repo-server", "imagePullPolicy:", addLifecycle)
}
