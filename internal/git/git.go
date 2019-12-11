package git

import (
	"context"
	"os"
	"time"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	git "gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Git struct {
	Repo     *git.Repository
	prevTree *object.Tree
	logger   logging.Logger
	auth     *gitssh.PublicKeys
}

func New(logger logging.Logger, localPath, url string, privateKey []byte) (*Git, error) {

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	auth := &gitssh.PublicKeys{User: "git", Signer: signer}
	auth.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	g := &Git{logger: logger, auth: auth}

	repoLogger := g.logger.WithFields(map[string]interface{}{
		"repo": url,
	})

	repo, err := g.cloneRepo(localPath, url)
	if err != nil {
		return nil, errors.Wrapf(err, "Cloning repo %s failed", url)
	}
	g.Repo = repo

	g.logger.WithFields(map[string]interface{}{
		"logID": "GIT-4Sia0VjJ79gb7cw",
	}).Info("Cloned...")
	ref, err := g.Repo.Head()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get head from repo %s", url)
	}

	g.logger.WithFields(map[string]interface{}{
		"logID": "GIT-4Sia0VjJ79gb7cw",
	}).Info("Get last commit...")
	commit, err := g.Repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get last commit from repo %s", url)
	}
	prevTree, err := commit.Tree()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get tree of last commit from repo %s", url)
	}
	g.prevTree = prevTree
	repoLogger.WithFields(map[string]interface{}{
		"logID": "GIT-pQnw5FfIqAk0eWc",
		"path":  localPath,
	}).Info("Cloned new GitCRD")

	return g, nil
}

func (g *Git) cloneRepo(localPath, url string) (*git.Repository, error) {

	g.logger.WithFields(map[string]interface{}{
		"logID": "GIT-vNU9maj2Rfo5rRU",
		"repo":  url,
		"to":    localPath,
	}).Info("Cloning plain")

	ctx := context.TODO()
	toCtx, _ := context.WithTimeout(ctx, 10*time.Second)
	return git.PlainCloneContext(toCtx, localPath, false, &git.CloneOptions{
		URL:          url,
		SingleBranch: true,
		Depth:        1,
		Progress:     os.Stdout,
		Auth:         g.auth,
	})
}

func (g *Git) IsFileChanged(path string) (changed bool, err error) {

	var action string
	defer func() {
		if err != nil {
			err = errors.Wrapf(err, "Failed to %s of repo", action)
		}
	}()

	w, err := g.Repo.Worktree()
	if err != nil {
		action = "get workingtree"
		return false, err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       g.auth,
	})
	if err == git.NoErrAlreadyUpToDate {
		return false, nil
	}
	if err != nil {
		action = "pull"
		return false, err
	}

	ref, err := g.Repo.Head()
	if err != nil {
		action = "get the head"
		return false, err
	}

	commit, err := g.Repo.CommitObject(ref.Hash())
	if err != nil {
		action = "get last commit"
		return false, err
	}

	currentTree, err := commit.Tree()
	if err != nil {
		action = "get tree of last commit"
		return false, err
	}

	changes, err := currentTree.Diff(g.prevTree)
	if err != nil {
		action = "diff changes"
		return false, err
	}
	g.prevTree = currentTree

	for _, c := range changes {
		if c.To.Name == path {
			return true, nil
		}
	}
	return false, nil

}
