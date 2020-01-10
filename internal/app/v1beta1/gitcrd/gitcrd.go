package gitcrd

import (
	"os"
	"path/filepath"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/v1beta1/crd"
	"github.com/caos/boom/internal/git"
	"github.com/caos/boom/internal/helper"
)

type GitCrd struct {
	crd              *crd.Crd
	git              *git.Git
	crdDirectoryPath string
	crdPath          string
}

func New(logger logging.Logger, git *git.Git, crdDirectoryPath, crdPath, toolsDirectoryPath, dashboardDirectoryPath string) (*GitCrd, error) {

	gitCrd := &GitCrd{
		crdDirectoryPath: crdDirectoryPath,
		crdPath:          crdPath,
		git:              git,
	}

	toolsetCRD, err := gitCrd.GetCrdContent()
	if err != nil {
		return nil, err
	}

	localcrd, err := crd.New(logger, toolsetCRD.Name, toolsDirectoryPath, dashboardDirectoryPath)
	if err != nil {
		return nil, err
	}

	gitCrd.crd = localcrd
	return gitCrd, nil
}

func (c *GitCrd) CleanUp() error {
	return os.RemoveAll(c.crdDirectoryPath)
}

func (c *GitCrd) Reconcile() error {
	if err := c.git.ReloadRepo(); err != nil {
		return err
	}

	toolsetCRD, err := c.GetCrdContent()
	if err != nil {
		return err
	}

	return c.crd.Reconcile(toolsetCRD)
}

func (c *GitCrd) GetCrdContent() (*toolsetsv1beta1.Toolset, error) {
	crdFilePath := filepath.Join(c.crdDirectoryPath, c.crdPath)

	toolsetCRD := &toolsetsv1beta1.Toolset{}
	if err := helper.YamlToStruct(crdFilePath, toolsetCRD); err != nil {
		return nil, err
	}
	return toolsetCRD, nil
}
