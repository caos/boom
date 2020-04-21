package gitcrd

import (
	"context"
	"github.com/caos/boom/internal/metrics"
	"strings"

	toolsetv1beta1 "github.com/caos/boom/api/v1beta1"
	bundleconfig "github.com/caos/boom/internal/bundle/config"
	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/git"
	"github.com/caos/boom/internal/gitcrd/config"
	"github.com/caos/boom/internal/gitcrd/v1beta1"
	v1beta1config "github.com/caos/boom/internal/gitcrd/v1beta1/config"
	"github.com/caos/boom/internal/helper"
	"github.com/pkg/errors"
)

type GitCrd interface {
	SetBundle(*bundleconfig.Config)
	Reconcile([]*clientgo.Resource)
	WriteBackCurrentState([]*clientgo.Resource)
	CleanUp()
	GetStatus() error
	SetBackStatus()
	GetRepoURL() string
	GetRepoCRDPath() string
}

func New(conf *config.Config) (GitCrd, error) {

	conf.Monitor.Info("New GitCRD")

	gitInternal := git.New(context.Background(), conf.Monitor, conf.User, conf.Email, conf.CrdUrl)
	err := gitInternal.Init(conf.PrivateKey)
	if err != nil {
		conf.Monitor.Error(err)
		return nil, err
	}

	err = gitInternal.Clone()
	if err != nil {
		metrics.FailedGitClone(conf.CrdUrl)
		conf.Monitor.Error(err)
		return nil, err
	}
	metrics.SuccessfulGitClone(conf.CrdUrl)

	crdFileStruct := &helper.Resource{}
	if err := gitInternal.ReadYamlIntoStruct(conf.CrdPath, crdFileStruct); err != nil {
		metrics.WrongCRDFormat(conf.CrdUrl, conf.CrdPath)
		conf.Monitor.Error(err)
		return nil, err
	}

	groupVersion := toolsetv1beta1.GroupVersion
	parts := strings.Split(crdFileStruct.ApiVersion, "/")
	if parts[0] != "boom.caos.ch" {
		err := errors.Errorf("Unknown CRD apiGroup %s", parts[0])
		conf.Monitor.Error(err)
		metrics.UnsupportedAPIGroup(conf.CrdUrl, conf.CrdPath)
		return nil, err
	}

	if parts[1] != groupVersion.Version {
		err := errors.Errorf("Unknown CRD version %s", parts[1])
		conf.Monitor.Error(err)
		metrics.UnsupportedVersion(conf.CrdUrl, conf.CrdPath)
		return nil, err
	}

	metrics.SuccessfulUnmarshalCRD(conf.CrdUrl, conf.CrdPath)
	monitor := conf.Monitor.WithFields(map[string]interface{}{
		"type": "gitcrd",
	})

	v1beta1conf := &v1beta1config.Config{
		Monitor:          monitor,
		Git:              gitInternal,
		CrdDirectoryPath: conf.CrdDirectoryPath,
		CrdPath:          conf.CrdPath,
	}

	return v1beta1.New(v1beta1conf)
}
