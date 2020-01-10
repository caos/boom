package app

import (
	"path/filepath"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	v1beta1crd "github.com/caos/boom/internal/app/v1beta1/crd"
	v1beta1gitcrd "github.com/caos/boom/internal/app/v1beta1/gitcrd"
	"github.com/caos/boom/internal/git"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/orbiter/logging"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pkg/errors"
)

type Crd interface {
	ReconcileWithFunc(getToolsetCRD func(instance runtime.Object) error) error
	Reconcile(toolsetCRD *toolsetsv1beta1.Toolset) error
	CleanUp() error
}

func NewCrd(logger logging.Logger, version string, getToolsetCRD func(instance runtime.Object) error, toolsDirectoryPath, dashboardsDirectoryPath string) (Crd, error) {

	crdLogger := logger.WithFields(map[string]interface{}{
		"version": version,
	})
	crdLogger.WithFields(map[string]interface{}{
		"logID": "CRD-OieUWt0rdMoRrIh",
	}).Info("New CDR")

	if version != "v1beta1" {
		return nil, errors.Errorf("Unknown CRD version %s", version)
	}

	var toolsetCRD *toolsetsv1beta1.Toolset
	if err := getToolsetCRD(toolsetCRD); err != nil {
		return nil, err
	}

	crd, err := v1beta1crd.New(logger, toolsetCRD.Name, toolsDirectoryPath, dashboardsDirectoryPath)
	if err != nil {
		return nil, err
	}

	return crd, crd.Reconcile(toolsetCRD)
}

type GitCrd interface {
	Reconcile() error
	CleanUp() error
}

func NewGitCrd(logger logging.Logger, crdDirectoryPath, crdUrl string, privateKey []byte, crdPath, toolsDirectoryPath, dashboardDirectoryPath string) (GitCrd, error) {

	logger.WithFields(map[string]interface{}{
		"logID": "CRD-OieUWt0rdMoRrIh",
		"repo":  crdUrl,
		"path":  crdPath,
	}).Info("New GitCRD")

	git, err := git.New(logger, crdDirectoryPath, crdUrl, privateKey)
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

	crd, err := v1beta1gitcrd.New(logger, git, crdDirectoryPath, crdPath, toolsDirectoryPath, dashboardDirectoryPath)
	if err != nil {
		return nil, err
	}

	return crd, crd.Reconcile()
}
