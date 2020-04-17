package customimage

import (
	"encoding/json"
	"github.com/caos/boom/internal/bundle/application/applications/argocd/info"
	"github.com/caos/boom/internal/bundle/application/resources"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/boom/internal/secret"
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

func getSecretName(store string, ty string) string {
	return strings.Join([]string{"argocd", getInternalName(store, ty)}, "-")
}
func getSecretKey(store string, ty string) string {
	return strings.Join([]string{store, ty}, "-")
}
func getInternalName(store string, ty string) string {
	return strings.Join([]string{"store", store, ty}, "-")
}

func GetSecrets(spec *toolsetsv1beta1.Argocd) []interface{} {
	namespace := "caos-system"

	secrets := make([]interface{}, 0)
	for _, store := range spec.CustomImage.GopassStores {
		if helper.IsCrdSecret(store.GPGKey, store.ExistingGPGKeySecret) {
			ty := "gpg"
			data := map[string]string{
				getSecretKey(store.StoreName, ty): store.GPGKey.Value,
			}

			conf := &resources.SecretConfig{
				Name:      getSecretName(store.StoreName, ty),
				Namespace: namespace,
				Labels:    labels.GetAllApplicationLabels(info.GetName()),
				Data:      data,
			}
			secretRes := resources.NewSecret(conf)
			secrets = append(secrets, secretRes)
		}
	}

	return secrets
}

func FromSpec(spec *toolsetsv1beta1.Argocd, imageTags map[string]string) *CustomImage {
	imageRepository := "docker.pkg.github.com/caos/argocd-secrets/argocd"

	vols := make([]*SecretVolume, 0)
	volMounts := make([]*VolumeMount, 0)
	for _, store := range spec.CustomImage.GopassStores {

		volGPG, volMountGPG := getVolAndVolMount(store.StoreName, "gpg", store.GPGKey, store.ExistingGPGKeySecret)
		vols = append(vols, volGPG)
		volMounts = append(volMounts, volMountGPG)

		volSSH, volMountSSH := getVolAndVolMount(store.StoreName, "ssh", store.SSHKey, store.ExistingSSHKeySecret)
		vols = append(vols, volSSH)
		volMounts = append(volMounts, volMountSSH)
	}

	return &CustomImage{
		ImageRepository:  imageRepository,
		ImageTag:         imageTags[imageRepository],
		AddSecretVolumes: vols,
		AddVolumeMounts:  volMounts,
	}
}

func getVolAndVolMount(storeName string, ty string, secret *secret.Secret, existent *secret.Existing) (*SecretVolume, *VolumeMount) {
	internalName := ""
	name := ""
	key := ""

	if helper.IsCrdSecret(secret, existent) {
		internalName = getInternalName(storeName, ty)
		name = getSecretName(storeName, ty)
		key = getSecretKey(storeName, ty)
	} else if helper.IsExistentSecret(secret, existent) {
		internalName = existent.InternalName
		name = existent.Name
		key = existent.Key
	} else {
		//TODO
	}

	return getVol(internalName, name, key), getVolMount(internalName)
}

func getVol(internal string, name string, key string) *SecretVolume {
	return &SecretVolume{
		Name: internal,
		Secret: &Secret{
			SecretName: name,
			Items: []*Item{&Item{
				Key:  key,
				Path: internal,
			},
			},
		},
		DefaultMode: 0544,
	}
}

func getVolMount(internal string) *VolumeMount {
	mountPath := filepath.Join(gpgFolderName, internal)
	return &VolumeMount{
		Name:      internal,
		MountPath: mountPath,
		SubPath:   internal,
		ReadOnly:  false,
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
