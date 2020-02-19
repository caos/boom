package v1beta1

import (
	"io/ioutil"
	"os"
	"path/filepath"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/crd"
	"github.com/caos/boom/internal/crd/v1beta1"
	"github.com/caos/boom/internal/gitcrd/v1beta1/config"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	crdconfig "github.com/caos/boom/internal/crd/config"
	"github.com/caos/boom/internal/current"
	"github.com/caos/boom/internal/git"

	"github.com/caos/boom/internal/helper"
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

	if toolsetCRD.Spec.PreApply != nil || toolsetCRD.Spec.PreApply.Deploy {
		command := "delete"
		if toolsetCRD.Spec.PostApply.Deploy {
			command = "apply"
		}

		err := c.ApplyFolders(c.git, command, toolsetCRD.Spec.PreApply.Folder)
		if err != nil {
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

	if toolsetCRD.Spec.PostApply != nil {
		command := "delete"
		if toolsetCRD.Spec.PostApply.Deploy {
			command = "apply"
		}

		err := c.ApplyFolders(c.git, command, toolsetCRD.Spec.PostApply.Folder)
		if err != nil {
			c.status = err
			return
		}
	}

}

func (c *GitCrd) ApplyFolders(git *git.Client, command, folderRelativePath string) error {
	err := git.Clone()
	if err != nil {
		return err
	}

	folderPath := filepath.Join(c.crdDirectoryPath, folderRelativePath)
	err = helper.RecreatePath(folderPath)
	if err != nil {
		return err
	}

	files, err := git.ReadFolder(folderRelativePath)
	for filename, file := range files {
		filePath := filepath.Join(folderPath, filename)
		err := ioutil.WriteFile(filePath, file, os.ModePerm)
		if err != nil {
			return err
		}
	}

	applyCmd := kubectl.New(command).AddParameter("-f", folderPath).Build()
	err = helper.Run(c.logger, applyCmd)
	return err
}

func (c *GitCrd) GetCrdContent() (*toolsetsv1beta1.Toolset, error) {
	if err := c.git.Clone(); err != nil {
		return nil, err
	}

	data, err := c.git.Read(c.crdPath)

	toolsetCRD := &toolsetsv1beta1.Toolset{}
	err = yaml.Unmarshal(data, toolsetCRD)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while unmarshaling yaml %s to struct", c.crdPath)
	}

	return toolsetCRD, nil
}

func (c *GitCrd) WriteBackCurrentState() {
	if c.status != nil {
		return
	}

	curr := current.Get(c.logger)

	content, err := yaml.Marshal(curr)
	if err != nil {
		c.status = err
		return
	}

	file := git.File{
		Path:    filepath.Join("internal", "boom", "current.yaml"),
		Content: content,
	}

	if err := c.git.Clone(); err != nil {
		c.status = err
		return
	}

	c.status = c.git.UpdateRemote("current state changed", file)
}
