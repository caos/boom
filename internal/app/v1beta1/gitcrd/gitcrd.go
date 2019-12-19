package gitcrd

import (
	"os"
	"path/filepath"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	appcrd "github.com/caos/boom/internal/app/v1beta1/crd"
	"github.com/caos/boom/internal/git"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/toolset"
)

type GitCrd struct {
	crd              *appcrd.Crd
	git              *git.Git
	crdDirectoryPath string
	crdPath          string
}

func New(logger logging.Logger, git *git.Git, crdDirectoryPath, crdPath, toolsDirectoryPath string, toolsets *toolset.Toolsets) (*GitCrd, error) {

	gitCrd := &GitCrd{
		crdDirectoryPath: crdDirectoryPath,
		crdPath:          crdPath,
		git:              git,
	}

	toolsetCRD, err := gitCrd.GetCrdContent()
	if err != nil {
		return nil, err
	}

	localcrd, err := appcrd.New(logger, toolsetCRD, toolsDirectoryPath, toolsets)
	if err != nil {
		return nil, err
	}
	gitCrd.crd = localcrd
	return gitCrd, nil
}

func (c *GitCrd) CleanUp() error {
	return os.RemoveAll(c.crdDirectoryPath)
}

func (c *GitCrd) Reconcile(toolsDirectoryPath string, toolsets *toolset.Toolsets) error {
	if err := c.git.ReloadRepo(); err != nil {
		return err
	}

	new, err := c.GetCrdContent()
	if err != nil {
		return err
	}

	return c.crd.Reconcile(new, toolsDirectoryPath, toolsets)
}

func (c *GitCrd) GetCrdContent() (*toolsetsv1beta1.Toolset, error) {
	crdFilePath := filepath.Join(c.crdDirectoryPath, c.crdPath)

	toolsetCRD := &toolsetsv1beta1.Toolset{}
	if err := helper.YamlToStruct(crdFilePath, toolsetCRD); err != nil {
		return nil, err
	}
	return toolsetCRD, nil
}
