package app

import (
	"path/filepath"

	"github.com/caos/orbiter/logging"
	v1beta1crd "github.com/caos/boom/internal/app/v1beta1/crd"
	v1beta1gitcrd "github.com/caos/boom/internal/app/v1beta1/gitcrd"
	"github.com/caos/boom/internal/git"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/toolset"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pkg/errors"
)

type Crd interface {
	ReconcileWithFunc(func(obj runtime.Object) error, string, *toolset.Toolsets) error
	CleanUp() error
}

func NewCrd(logger logging.Logger, version string, getToolset func(obj runtime.Object) error, toolsDirectoryPath string, toolsets *toolset.Toolsets) (Crd, error) {

	crdLogger := logger.WithFields(map[string]interface{}{
		"version": version,
	})
	crdLogger.WithFields(map[string]interface{}{
		"logID": "CRD-OieUWt0rdMoRrIh",
	}).Info("New CDR")
	if version != "v1beta1" {
		return nil, errors.Errorf("Unknown CRD version %s", version)
	}
	return v1beta1crd.NewWithFunc(logger, getToolset, toolsDirectoryPath, toolsets)
}

type GitCrd interface {
	Reconcile(string, *toolset.Toolsets) error
	CleanUp() error
}

func NewGitCrd(logger logging.Logger, crdDirectoryPath, crdUrl, crdSecretPath, crdPath, toolsDirectoryPath string, toolsets *toolset.Toolsets) (GitCrd, error) {

	logger.WithFields(map[string]interface{}{
		"logID": "CRD-OieUWt0rdMoRrIh",
		"repo":  crdUrl,
		"path":  crdPath,
	}).Info("New GitCRD")

	git, err := git.New(logger, crdDirectoryPath, crdUrl, crdSecretPath)
	if err != nil {
		return nil, err
	}

	crdFilePath := filepath.Join(crdDirectoryPath, crdPath)
	version, err := helper.GetVersionFromYaml(crdFilePath)
	if err != nil {
		return nil, err
	}
	if version != "v1beta1" {
		return nil, errors.Errorf("Unknown CRD version %s", version)
	}
	return v1beta1gitcrd.New(logger, git, crdDirectoryPath, crdPath, toolsDirectoryPath, toolsets)
}
