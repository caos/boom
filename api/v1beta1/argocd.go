package v1beta1

type Argocd struct {
	Deploy                bool   `json:"deploy,omitempty"`
	CustomImageWithGopass bool   `json:"customImageWithGopass,omitempty" yaml:"customImageWithGopass,omitempty"`
	ImagePullSecret       string `json:"imagePullSecret,omitempty" yaml:"imagePullSecret,omitempty"`
}
