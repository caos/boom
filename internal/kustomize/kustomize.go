package kustomize

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type Kustomize struct {
	path string
}

func New(path string) (*Kustomize, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return &Kustomize{
		path: abspath,
	}, nil
}

func (k *Kustomize) Build() exec.Cmd {
	all := strings.Join([]string{"kustomize", "build", k.path, "| kubectl apply -f -"}, " ")
	cmd := exec.Command("/bin/sh", "-c", all)
	return *cmd
}
