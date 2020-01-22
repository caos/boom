package v1beta1

import (
	"os"
	"path/filepath"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/crd"
	"github.com/caos/boom/internal/crd/v1beta1"
	"github.com/caos/boom/internal/gitcrd/v1beta1/config"

	crdconfig "github.com/caos/boom/internal/crd/config"
	"github.com/caos/boom/internal/git"

	"github.com/caos/boom/internal/helper"
)

type GitCrd struct {
	crd              crd.Crd
	git              *git.Git
	crdDirectoryPath string
	crdPath          string
	status           error
}

func New(conf *config.Config) (*GitCrd, error) {

	gitCrd := &GitCrd{
		crdDirectoryPath: conf.CrdDirectoryPath,
		crdPath:          conf.CrdPath,
		git:              conf.Git,
	}

	// toolsetCRD, err := gitCrd.GetCrdContent()
	// if err != nil {
	// 	return nil, err
	// }

	crdConf := &crdconfig.Config{
		Logger:  conf.Logger,
		Version: v1beta1.GetVersion(),
	}

	crd, err := crd.New(crdConf)
	if err != nil {
		return nil, err
	}
	gitCrd.crd = crd

	return gitCrd, nil
}

func (c *GitCrd) GetStatus() error {
	return c.status
}

func (c *GitCrd) SetBundle(conf *bundleconfig.Config) {
	if c.status != nil {
		return
	}

	bundleConf := &bundleconfig.Config{
		Logger:            conf.Logger,
		CrdName:           conf.CrdName,
		BundleName:        conf.BundleName,
		BaseDirectoryPath: conf.BaseDirectoryPath,
		Templator:         conf.Templator,
	}

	c.crd.SetBundle(bundleConf)
	c.status = c.crd.GetStatus()
}

func (c *GitCrd) CleanUp() {
	if c.status != nil {
		return
	}
	c.status = os.RemoveAll(c.crdDirectoryPath)
}

func (c *GitCrd) Reconcile() {
	if err := c.git.ReloadRepo(); err != nil {
		c.status = err
		return
	}

	toolsetCRD, err := c.GetCrdContent()
	if err != nil {
		c.status = err
		return
	}
	c.crd.Reconcile(toolsetCRD)

	c.status = c.crd.GetStatus()
}

func (c *GitCrd) GetCrdContent() (*toolsetsv1beta1.Toolset, error) {
	crdFilePath := filepath.Join(c.crdDirectoryPath, c.crdPath)

	toolsetCRD := &toolsetsv1beta1.Toolset{}
	if err := helper.YamlToStruct(crdFilePath, toolsetCRD); err != nil {
		return nil, err
	}
	return toolsetCRD, nil
}
