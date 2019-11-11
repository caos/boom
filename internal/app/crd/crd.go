package app

import (
	"os"
	"strings"

	"github.com/caos/toolsop/internal/git"
	"github.com/caos/toolsop/internal/helper"
	"github.com/caos/toolsop/internal/template"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
)

type Crd struct {
	git           *git.Git
	crdPath       string
	directoryPath string
	generateFunc  func(string, string, string) error
	reconcileFunc func(string, *toolsetsv1beta1.ToolsetSpec) error
	helm          *template.Helm
}

func New(directoryPath, url, secretPath, crdPath string, generateFunc func(string, string, string) error, reconcileFunc func(string, *toolsetsv1beta1.ToolsetSpec) error) (*Crd, error) {
	gitCrd := &Crd{
		directoryPath: directoryPath,
		crdPath:       crdPath,
		generateFunc:  generateFunc,
		reconcileFunc: reconcileFunc,
	}

	g, err := git.New(directoryPath, url, secretPath)
	if err != nil {
		return nil, err
	}
	gitCrd.git = g

	return gitCrd, nil
}

func (c *Crd) CleanUp() error {
	return os.RemoveAll(c.directoryPath)
}

func (c *Crd) Maintain() error {

	changed, err := c.git.IsFileChanged(c.crdPath)
	if err != nil || !changed {
		return err
	}

	return c.Apply()
}

func (c *Crd) Apply() error {
	crdFilePath := strings.Join([]string{c.directoryPath, c.crdPath}, "/")

	toolsetCRD := &toolsetsv1beta1.Toolset{}
	if err := helper.YamlToStruct(crdFilePath, toolsetCRD); err != nil {
		return err
	}

	if err := c.generateFunc("caos", toolsetCRD.Spec.Name, toolsetCRD.Spec.Version); err != nil {
		return err
	}

	if err := c.reconcileFunc("caos", toolsetCRD.Spec); err != nil {
		return err
	}

	return nil
}
