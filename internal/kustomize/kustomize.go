package kustomize

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type Kustomize struct {
	path  string
	apply bool
}

func New(path string, apply bool) (*Kustomize, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return &Kustomize{
		path:  abspath,
		apply: apply,
	}, nil
}

func (k *Kustomize) Build() exec.Cmd {
	all := strings.Join([]string{"kustomize", "build", k.path}, " ")
	if k.apply {
		all = strings.Join([]string{all, "| kubectl apply -f -"}, " ")
	}
	cmd := exec.Command("/bin/sh", "-c", all)
	return *cmd
}
