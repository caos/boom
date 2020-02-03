package v1beta1

type Argocd struct {
	Deploy                bool   `json:"deploy,omitempty"`
	CustomImageWithGopass bool   `json:"customImageWithGopass,omitempty" yaml:"customImageWithGopass,omitempty"`
	ImagePullSecret       string `json:"imagePullSecret,omitempty" yaml:"imagePullSecret,omitempty"`
	GopassGPGKey          string `json:"gopassGPGKey,omitempty" yaml:"gopassGPGKey,omitempty"`
	GopassSSHKey          string `json:"gopassSSHKey,omitempty" yaml:"gopassSSHKey,omitempty"`
	GopassDirectory       string `json:"gopassDirectory,omitempty" yaml:"gopassDirectory,omitempty"`
	GopassStoreName       string `json:"gopassStoreName,omitempty" yaml:"gopassStoreName,omitempty"`
}
