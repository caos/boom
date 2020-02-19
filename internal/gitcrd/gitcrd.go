package gitcrd

import (
	"context"
	"strings"

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
	WriteBackCurrentState()
	CleanUp()
	GetStatus() error
	GetCrdContent() (*toolsetsv1beta1.Toolset, error)
}

func New(conf *config.Config) (GitCrd, error) {

	logFields := map[string]interface{}{
		"logID": "CRD-OieUWt0rdMoRrIh",
	}

	conf.Logger.WithFields(logFields).Info("New GitCRD")

	git := git.New(context.Background(), conf.Logger, conf.User, conf.Email, conf.CrdUrl)
	err := git.Init(conf.PrivateKey)
	if err != nil {
		return nil, err
	}

	err = git.Clone()
	if err != nil {
		return nil, err
	}

	crdFileStruct := &helper.Resource{}
	if err := git.ReadYamlIntoStruct(conf.CrdPath, crdFileStruct); err != nil {
		conf.Logger.WithFields(logFields).Error(err)
		return nil, err
	}

	parts := strings.Split(crdFileStruct.ApiVersion, "/")
	if parts[0] != "boom.caos.ch" {
		return nil, errors.Errorf("Unknown CRD apiGroup %s", parts[0])
	}

	if parts[1] != "v1beta1" {
		return nil, errors.Errorf("Unknown CRD version %s", parts[1])
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
