package plugin

type Plugin struct {
	Name     string   `yaml:"name"`
	Init     *Command `yaml:"init,omitempty"`
	Generate *Command `yaml:"generate,omitempty"`
}

type Command struct {
	Command []string `yaml:"command,omitempty"`
	Args    []string `yaml:"args,omitempty"`
}

func New(name string, init *Command, generate *Command) *Plugin {
	return &Plugin{
		Name:     name,
		Init:     init,
		Generate: generate,
	}
}
