package v1beta1

type CurrentState struct {
	WriteBack bool   `json:"writeBack" yaml:"writeBack"`
	Folder    string `json:"folder" yaml:"folder"`
}
