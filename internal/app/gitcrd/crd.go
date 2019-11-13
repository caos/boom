package gitcrd

import (
	"os"
	"strings"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
	appcrd "github.com/caos/toolsop/internal/app/crd"
	"github.com/caos/toolsop/internal/git"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/toolset"
)

type GitCrd struct {
	crd              *appcrd.Crd
	git              *git.Git
	crdDirectoryPath string
	crdPath          string
}

func New(crdDirectoryPath, crdUrl, crdSecretPath, crdPath, toolsDirectoryPath string, toolsets *toolset.Toolsets) (*GitCrd, error) {

	gitCrd := &GitCrd{
		crdDirectoryPath: crdDirectoryPath,
		crdPath:          crdPath,
	}

	g, err := git.New(crdDirectoryPath, crdUrl, crdSecretPath)
	if err != nil {
		return nil, err
	}
	gitCrd.git = g

	toolsetCRD, err := gitCrd.GetCrdContent()
	if err != nil {
		return nil, err
	}

	localcrd, err := appcrd.New(toolsetCRD, toolsDirectoryPath, toolsets)
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
	changed, err := c.git.IsFileChanged(c.crdPath)
	if err != nil || !changed {
		return err
	}

	new, err := c.GetCrdContent()
	if err != nil {
		return err
	}

	return c.crd.Reconcile(new, toolsDirectoryPath, toolsets)
}

func (c *GitCrd) GetCrdContent() (*toolsetsv1beta1.Toolset, error) {
	crdFilePath := strings.Join([]string{c.crdDirectoryPath, c.crdPath}, "/")

	toolsetCRD := &toolsetsv1beta1.Toolset{}
	if err := helper.YamlToStruct(crdFilePath, toolsetCRD); err != nil {
		return nil, err
	}
	return toolsetCRD, nil
}
