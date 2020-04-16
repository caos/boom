package secret

type Existing struct {
	Name string `json:"name" yaml:"name"`
	Key  string `json:"key" yaml:"key"`
}

type ExistingToFilesystem struct {
	Name         string `json:"name" yaml:"name"`
	Key          string `json:"key" yaml:"key"`
	InternalName string `json:"internalName" yaml:"internalName"`
}

type Secret struct {
	Encryption string `yaml:"encryption"`
	Encoding   string `yaml:"encoding"`
	Value      string `yaml:"value"`
	Masterkey  string `yaml:"-"`
}
