package app

import (
	"path/filepath"

	v1beta1crd "github.com/caos/toolsop/internal/app/v1beta1/crd"
	v1beta1gitcrd "github.com/caos/toolsop/internal/app/v1beta1/gitcrd"
	"github.com/caos/toolsop/internal/git"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/toolset"
	"k8s.io/apimachinery/pkg/runtime"
)

type Crd interface {
	ReconcileWithFunc(func(obj runtime.Object) error, string, *toolset.Toolsets) error
}

func NewCrd(version string, getToolset func(obj runtime.Object) error, toolsDirectoryPath string, toolsets *toolset.Toolsets) (Crd, error) {

	if version == "v1beta1" {
		return v1beta1crd.NewWithFunc(getToolset, toolsDirectoryPath, toolsets)
	}

	return nil, nil
}

type GitCrd interface {
	Reconcile(string, *toolset.Toolsets) error
	CleanUp() error
}

func NewGitCrd(crdDirectoryPath, crdUrl, crdSecretPath, crdPath, toolsDirectoryPath string, toolsets *toolset.Toolsets) (GitCrd, error) {

	git, err := git.New(crdDirectoryPath, crdUrl, crdSecretPath)
	if err != nil {
		return nil, err
	}

	crdFilePath := filepath.Join(crdDirectoryPath, crdPath)
	version, err := helper.GetVersionFromYaml(crdFilePath)
	if err != nil {
		return nil, err
	}
	if version == "v1beta1" {
		return v1beta1gitcrd.New(git, crdDirectoryPath, crdPath, toolsDirectoryPath, toolsets)
	}

	return nil, nil
}
