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
	"github.com/caos/boom/internal/kubectl"
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

	crdLogger := conf.Logger.WithFields(map[string]interface{}{
		"version": "v1beta1",
	})

	gitConf := *conf.Git
	gitCrd := &GitCrd{
		crdDirectoryPath: conf.CrdDirectoryPath,
		crdPath:          conf.CrdPath,
		git:              &gitConf,
		logger:           crdLogger,
	}

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

	toolsetCRD, err := c.getCrdContent()
	if err != nil {
		c.status = err
		return
	}

	logger := c.logger.WithFields(map[string]interface{}{
		"CRD": toolsetCRD.GetName(),
	})

	bundleConf := &bundleconfig.Config{
		Logger:            logger,
		CrdName:           toolsetCRD.GetName(),
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
	if c.status != nil {
		return
	}

	logger := c.logger.WithFields(map[string]interface{}{
		"action": "reconiling",
	})

	toolsetCRD, err := c.getCrdContent()
	if err != nil {
		c.status = err
		return
	}

	// pre-steps
	if toolsetCRD.Spec.PreApply != nil {
		pre := toolsetCRD.Spec.PreApply
		if err := helper.CopyFolderToLocal(c.git, c.crdDirectoryPath, pre.Folder); err != nil {
			c.status = err
			return
		}

		if err := useFolder(logger, pre.Deploy, c.crdDirectoryPath, pre.Folder); err != nil {
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
		if err := helper.CopyFolderToLocal(c.git, c.crdDirectoryPath, post.Folder); err != nil {
			c.status = err
			return
		}

		if err := useFolder(logger, post.Deploy, c.crdDirectoryPath, post.Folder); err != nil {
			c.status = err
			return
		}
	}
}

func (c *GitCrd) getCrdContent() (*toolsetsv1beta1.Toolset, error) {
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

	logger := c.logger.WithFields(map[string]interface{}{
		"action": "current",
	})

	curr := current.Get(logger)

	content, err := yaml.Marshal(curr)
	if err != nil {
		c.status = err
		return
	}

	currentFolder := filepath.Join("internal", "boom")

	file := git.File{
		Path:    filepath.Join(currentFolder, "current.yaml"),
		Content: content,
	}

	c.status = c.git.UpdateRemote("current state changed", file)
}

func useFolder(logger logging.Logger, deploy bool, tempDirectory, folderRelativePath string) error {
	folderPath := filepath.Join(tempDirectory, folderRelativePath)

	command := kubectl.NewApply(folderPath).Build()
	if !deploy {
		command = kubectl.NewDelete(folderPath).Build()
	}

	return helper.Run(logger, command)
}
