package v1beta1

type PreApply struct {
	Deploy bool   `json:"deploy,omitempty"`
	Folder string `json:"folder,omitempty"`
}
