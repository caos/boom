package v1beta1

import (
	"os"
	"path/filepath"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/crd"
	"github.com/caos/boom/internal/crd/v1beta1"
	"github.com/caos/boom/internal/gitcrd/v1beta1/config"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	crdconfig "github.com/caos/boom/internal/crd/config"
	"github.com/caos/boom/internal/current"
	"github.com/caos/boom/internal/git"
)

type GitCrd struct {
	crd              crd.Crd
	git              *git.Client
	crdDirectoryPath string
	crdPath          string
	status           error
	logger           logging.Logger
}

func New(conf *config.Config) (*GitCrd, error) {

	gitCrd := &GitCrd{
		crdDirectoryPath: conf.CrdDirectoryPath,
		crdPath:          conf.CrdPath,
		git:              conf.Git,
		logger:           conf.Logger,
	}

	crdLogger := conf.Logger.WithFields(map[string]interface{}{
		"version": "v1beta1",
	})

	crdConf := &crdconfig.Config{
		Logger:  crdLogger,
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
	toolsetCRD, err := c.GetCrdContent()
	if err != nil {
		c.status = err
		return
	}

	// pre-steps
	if toolsetCRD.Spec.PreApply != nil {
		pre := toolsetCRD.Spec.PreApply
		if err := helper.UseFolder(c.logger, c.git, pre.Deploy, c.crdDirectoryPath, pre.Folder); err != nil {
			c.status = err
			return
		}
	}

	c.crd.Reconcile(toolsetCRD)
	err = c.crd.GetStatus()
	if err != nil {
		c.status = err
		return
	}

	// post-steps
	if toolsetCRD.Spec.PostApply != nil {
		post := toolsetCRD.Spec.PostApply

		if err := helper.UseFolder(c.logger, c.git, post.Deploy, c.crdDirectoryPath, post.Folder); err != nil {
			c.status = err
			return
		}
	}
}

func (c *GitCrd) GetCrdContent() (*toolsetsv1beta1.Toolset, error) {
	if err := c.git.Clone(); err != nil {
		return nil, err
	}

	toolsetCRD := &toolsetsv1beta1.Toolset{}
	err := c.git.ReadYamlIntoStruct(c.crdPath, toolsetCRD)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while unmarshaling yaml %s to struct", c.crdPath)
	}

	return toolsetCRD, nil
}

func (c *GitCrd) WriteBackCurrentState() {
	if c.status != nil {
		return
	}

	toolsetCRD, err := c.GetCrdContent()
	if err != nil {
		c.status = err
		return
	}

	if toolsetCRD.Spec.CurrentState == nil || !toolsetCRD.Spec.CurrentState.WriteBack {
		c.logger.Info("Write-back is deactivated, canceling")
		return
	}

	curr := current.Get(c.logger)

	content, err := yaml.Marshal(curr)
	if err != nil {
		c.status = err
		return
	}

	file := git.File{
		Path:    filepath.Join(toolsetCRD.Spec.CurrentState.Folder, "current.yaml"),
		Content: content,
	}

	if err := c.git.Clone(); err != nil {
		c.status = err
		return
	}

	c.status = c.git.UpdateRemote("current state changed", file)
}
