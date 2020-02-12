package gitcrd

import (
	"path/filepath"

	"github.com/caos/boom/internal/git"
	"github.com/caos/boom/internal/helper"
	"github.com/pkg/errors"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/gitcrd/config"
	"github.com/caos/boom/internal/gitcrd/v1beta1"
	v1beta1config "github.com/caos/boom/internal/gitcrd/v1beta1/config"
)

type GitCrd interface {
	SetBundle(*bundleconfig.Config)
	Reconcile()
	CleanUp()
	GetStatus() error
	GetCrdContent() (*toolsetsv1beta1.Toolset, error)
}

func New(conf *config.Config) (GitCrd, error) {

	logFields := map[string]interface{}{
		"logID": "CRD-OieUWt0rdMoRrIh",
	}

	conf.Logger.WithFields(logFields).Info("New GitCRD")

	git, err := git.New(conf.Logger, conf.CrdDirectoryPath, conf.CrdUrl, conf.PrivateKey)
	if err != nil {
		return nil, err
	}

	crdFilePath := filepath.Join(conf.CrdDirectoryPath, conf.CrdPath)
	group, err := helper.GetApiGroupFromYaml(crdFilePath)
	if err != nil {
		conf.Logger.WithFields(logFields).Error(err)
		return nil, err
	}

	if group != "boom.caos.ch" {
		return nil, errors.Errorf("Unknown CRD apiGroup %s", group)
	}

	version, err := helper.GetVersionFromYaml(crdFilePath)
	if err != nil {
		conf.Logger.WithFields(logFields).Error(err)
		return nil, err
	}

	if version != "v1beta1" {
		return nil, errors.Errorf("Unknown CRD version %s", version)
	}

	crdLogger := conf.Logger.WithFields(map[string]interface{}{
		"type": "gitcrd",
	})

	v1beta1conf := &v1beta1config.Config{
		Logger:           crdLogger,
		Git:              git,
		CrdDirectoryPath: conf.CrdDirectoryPath,
		CrdPath:          conf.CrdPath,
	}

	return v1beta1.New(v1beta1conf)
}
