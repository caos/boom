package v1beta1

type PostApply struct {
	Deploy bool   `json:"deploy,omitempty"`
	Folder string `json:"folder,omitempty"`
}
